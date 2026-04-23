package main

import (
	"context"
	"fmt"
	"os"
	"os/signal"

	"gateway-svc/adapter/catalog_grpc_adapter"
	"gateway-svc/adapter/finance_rest_adapter"
	"gateway-svc/adapter/iam_grpc_adapter"
	"gateway-svc/adapter/iam_rest_adapter"
	"gateway-svc/api/rest_oapi"
	"gateway-svc/service"
	"gateway-svc/util/config"
	"gateway-svc/util/logging"
	"gateway-svc/util/tracing"

	"go.opentelemetry.io/contrib/instrumentation/google.golang.org/grpc/otelgrpc"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

func start() {
	const op = "main.start"

	// --- Init logger (pick options from logging package; no config) ---
	rootLogger := logging.NewLogger(logging.Options{
		Level:    logging.LevelDebug,
		TimeZone: logging.TZWIB,
	})
	logger := &rootLogger

	// --- Load config ---
	config, err := config.LoadConfig(".")
	if err != nil {
		logger.Error().
			Str("op", op).
			Str("scope", "Load config").
			Err(err).
			Msg("")
		os.Exit(1)
	}

	// --- Init otel tracer ---
	cleanup, err := tracing.InitTracer(config.OtelTracer.Name, config.OtelTracer.Endpoint)
	if err != nil {
		logger.Error().
			Str("op", op).
			Str("scope", "Init otel tracer").
			Err(err).
			Msg("")
	}
	defer func() {
		if err := cleanup(context.Background()); err != nil {
			logger.Error().
				Str("op", op).
				Str("scope", "Cleanup otel tracer").
				Err(err).
				Msg("")
		}
	}()

	tracer := tracing.GetTracer(config.OtelTracer.Name)

	logger.Info().
		Str("op", op).
		Str("config", fmt.Sprintf("%+v", config)).
		Msg(fmt.Sprintf("Starting '%s' service ...", config.App.Name))

	// --- Dial iam-svc gRPC for edge bearer validation (BL-GTW-001 / F1-W7) ---
	//
	// Per ADR 0009, gateway validates every authenticated request once at the
	// edge via iam-svc.ValidateToken (gRPC). Traffic stays inside the
	// docker-compose network — insecure is fine today; the trust-contract
	// card (BL-GTW-100) adds mTLS. Unary stats handler propagates the current
	// trace context so iam-svc spans continue the gateway trace.
	iamConn, err := grpc.NewClient(
		config.External.IamSvc.GrpcTarget,
		grpc.WithTransportCredentials(insecure.NewCredentials()),
		grpc.WithStatsHandler(otelgrpc.NewClientHandler()),
	)
	if err != nil {
		logger.Error().
			Str("op", op).
			Str("scope", "Dial iam-svc gRPC").
			Str("target", config.External.IamSvc.GrpcTarget).
			Err(err).
			Msg("")
		os.Exit(1)
	}
	defer func() {
		if err := iamConn.Close(); err != nil {
			logger.Error().Err(err).Msg("close iam gRPC conn")
		}
	}()
	iamGrpcAdapter := iam_grpc_adapter.NewAdapter(logger, tracer, iamConn)

	// --- Dial catalog-svc gRPC for public catalog read (BL-GTW-002) ---
	//
	// Per ADR 0009, gateway's /v1/packages* + /v1/package-departures/{id}
	// proxy to catalog-svc.CatalogService over gRPC. Same dial + OTel
	// handler pattern as the iam conn above.
	catalogConn, err := grpc.NewClient(
		config.External.CatalogSvc.GrpcTarget,
		grpc.WithTransportCredentials(insecure.NewCredentials()),
		grpc.WithStatsHandler(otelgrpc.NewClientHandler()),
	)
	if err != nil {
		logger.Error().
			Str("op", op).
			Str("scope", "Dial catalog-svc gRPC").
			Str("target", config.External.CatalogSvc.GrpcTarget).
			Err(err).
			Msg("")
		os.Exit(1)
	}
	defer func() {
		if err := catalogConn.Close(); err != nil {
			logger.Error().Err(err).Msg("close catalog gRPC conn")
		}
	}()
	catalogGrpcAdapter := catalog_grpc_adapter.NewAdapter(logger, tracer, catalogConn)

	// --- Init REST adapters (interim; retire as each backend graduates to gRPC) ---
	// gateway-svc has no DB and no internal store; the service layer dispatches
	// to these per-backend adapters.
	//
	// - iam_rest_adapter: still serves the interim /v1/iam/system/live and
	//   /v1/iam/system/diagnostics/db-tx probes. Retires with BL-IAM-018 / S1-E-12.
	// - finance_rest_adapter: still serves /v1/finance/system/live and
	//   /v1/finance/ping (BL-IAM-002 permission-gate demo). Retires with
	//   BL-IAM-019 / S1-E-14.
	//
	// The seven pure-scaffold backends (booking/crm/jamaah/logistics/ops/
	// payment/visa) retired their REST surfaces in BL-REFACTOR-002..008 /
	// S1-E-13 — gateway has no adapters for them.
	iamAdapter := iam_rest_adapter.NewAdapter(logger, tracer, config.External.IamSvc.Address)
	financeAdapter := finance_rest_adapter.NewAdapter(logger, tracer, config.External.FinanceSvc.Address)

	// --- Init service layer ---
	svc := service.NewService(service.NewServiceParams{
		Logger:      logger,
		Tracer:      tracer,
		AppName:     config.App.Name,
		IamRest:     iamAdapter,
		CatalogGrpc: catalogGrpcAdapter,
		FinanceRest: financeAdapter,
	})

	// --- Init API layer (REST only — gateway is the edge proxy, no gRPC server) ---
	restServer := rest_oapi.NewServer(logger, tracer, svc, iamGrpcAdapter)

	// --- Run server ---
	runRestServer(config.Api.Rest.Port, restServer, config.Api.Metrics.Enabled, config.OtelTracer.Name)

	// --- Wait for signal ---
	ch := make(chan os.Signal, 1)
	signal.Notify(ch, os.Interrupt)

	<-ch

	logger.Info().Msg("end of program...")
}
