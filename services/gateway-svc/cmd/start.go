package main

import (
	"context"
	"fmt"
	"os"
	"os/signal"

	"gateway-svc/adapter/booking_grpc_adapter"
	"gateway-svc/adapter/catalog_grpc_adapter"
	"gateway-svc/adapter/crm_grpc_adapter"
	"gateway-svc/adapter/finance_grpc_adapter"
	"gateway-svc/adapter/finance_rest_adapter"
	"gateway-svc/adapter/iam_grpc_adapter"
	"gateway-svc/adapter/iam_rest_adapter"
	"gateway-svc/adapter/jamaah_grpc_adapter"
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

	// --- Dial booking-svc gRPC (BL-GTW-003 / S1-E-03) ---
	//
	// Per ADR-0009, gateway's /v1/bookings proxy to booking-svc.BookingService
	// over gRPC via booking_grpc_adapter. Same dial + OTel handler pattern as
	// the iam and catalog conns above.
	bookingConn, err := grpc.NewClient(
		config.External.BookingSvc.GrpcTarget,
		grpc.WithTransportCredentials(insecure.NewCredentials()),
		grpc.WithStatsHandler(otelgrpc.NewClientHandler()),
	)
	if err != nil {
		logger.Error().
			Str("op", op).
			Str("scope", "Dial booking-svc gRPC").
			Str("target", config.External.BookingSvc.GrpcTarget).
			Err(err).
			Msg("")
		os.Exit(1)
	}
	defer func() {
		if err := bookingConn.Close(); err != nil {
			logger.Error().Err(err).Msg("close booking gRPC conn")
		}
	}()
	bookingGrpcAdapter := booking_grpc_adapter.NewAdapter(logger, tracer, bookingConn)

	// --- Dial crm-svc gRPC (S4-E-02) ---
	//
	// Per ADR-0009, gateway's /v1/leads* routes proxy to crm-svc.CrmService over gRPC
	// via crm_grpc_adapter. Same dial + OTel handler pattern as the other conns above.
	crmConn, err := grpc.NewClient(
		config.External.CrmSvc.GrpcTarget,
		grpc.WithTransportCredentials(insecure.NewCredentials()),
		grpc.WithStatsHandler(otelgrpc.NewClientHandler()),
	)
	if err != nil {
		logger.Error().
			Str("op", op).
			Str("scope", "Dial crm-svc gRPC").
			Str("target", config.External.CrmSvc.GrpcTarget).
			Err(err).
			Msg("")
		os.Exit(1)
	}
	defer func() {
		if err := crmConn.Close(); err != nil {
			logger.Error().Err(err).Msg("close crm gRPC conn")
		}
	}()
	crmGrpcAdapter := crm_grpc_adapter.NewAdapter(logger, tracer, crmConn)

	// --- Dial finance-svc gRPC (S5-E-01) ---
	//
	// Per ADR-0009, gateway's /v1/finance/summary and /v1/finance/journals
	// proxy to finance-svc.FinanceService over gRPC via finance_grpc_adapter.
	// Same dial + OTel handler pattern as the other conns above.
	financeConn, err := grpc.NewClient(
		config.External.FinanceSvc.GrpcTarget,
		grpc.WithTransportCredentials(insecure.NewCredentials()),
		grpc.WithStatsHandler(otelgrpc.NewClientHandler()),
	)
	if err != nil {
		logger.Error().
			Str("op", op).
			Str("scope", "Dial finance-svc gRPC").
			Str("target", config.External.FinanceSvc.GrpcTarget).
			Err(err).
			Msg("")
		os.Exit(1)
	}
	defer func() {
		if err := financeConn.Close(); err != nil {
			logger.Error().Err(err).Msg("close finance gRPC conn")
		}
	}()
	financeGrpcAdapter := finance_grpc_adapter.NewAdapter(logger, tracer, financeConn)

	// --- Dial jamaah-svc gRPC (Phase 6 / Wave 1A) ---
	//
	// Per ADR-0009, gateway's /v1/manifest/:departure_id proxies to
	// jamaah-svc.JamaahService over gRPC via jamaah_grpc_adapter.
	jamaahConn, err := grpc.NewClient(
		config.External.JamaahSvc.GrpcTarget,
		grpc.WithTransportCredentials(insecure.NewCredentials()),
		grpc.WithStatsHandler(otelgrpc.NewClientHandler()),
	)
	if err != nil {
		logger.Error().
			Str("op", op).
			Str("scope", "Dial jamaah-svc gRPC").
			Str("target", config.External.JamaahSvc.GrpcTarget).
			Err(err).
			Msg("")
		os.Exit(1)
	}
	defer func() {
		if err := jamaahConn.Close(); err != nil {
			logger.Error().Err(err).Msg("close jamaah gRPC conn")
		}
	}()
	jamaahGrpcAdapter := jamaah_grpc_adapter.NewAdapter(logger, tracer, jamaahConn)

	// --- Init REST adapters (interim; retire as each backend graduates to gRPC) ---
	// gateway-svc has no DB and no internal store; the service layer dispatches
	// to these per-backend adapters.
	//
	// - iam_rest_adapter: interim; retires with BL-IAM-018 / S1-E-12 once iam-svc
	//   HTTP surface is fully removed.
	// - finance_rest_adapter: still serves /v1/finance/system/live. Retires with
	//   BL-IAM-019 / S1-E-14.
	//
	// The seven pure-scaffold backends (booking/crm/jamaah/logistics/ops/
	// payment/visa) retired their REST surfaces in BL-REFACTOR-002..008 /
	// S1-E-13 — gateway has no adapters for them. catalog_rest_adapter removed
	// in S1-E-11 (catalog-svc is gRPC-only).
	iamAdapter := iam_rest_adapter.NewAdapter(logger, tracer, config.External.IamSvc.Address)
	financeAdapter := finance_rest_adapter.NewAdapter(logger, tracer, config.External.FinanceSvc.Address)

	// --- Init service layer ---
	svc := service.NewService(service.NewServiceParams{
		Logger:      logger,
		Tracer:      tracer,
		AppName:     config.App.Name,
		IamRest:     iamAdapter,
		IamGrpc:     iamGrpcAdapter,
		CatalogGrpc: catalogGrpcAdapter,
		FinanceRest: financeAdapter,
		FinanceGrpc: financeGrpcAdapter,
		BookingGrpc: bookingGrpcAdapter,
		CrmGrpc:     crmGrpcAdapter,
		JamaahGrpc:  jamaahGrpcAdapter,
	})

	// --- Init API layer (REST only — gateway is the edge proxy, no gRPC server) ---
	restServer := rest_oapi.NewServer(logger, tracer, svc, iamGrpcAdapter)

	// --- Run server ---
	// iamGrpcAdapter is passed explicitly so the router can wire RequireBearerToken
	// to the protected route group (BL-GTW-001 / S1-E-09). restServer still holds
	// the same reference via iamValidator for any per-handler use, but the router
	// needs it at group construction time to avoid a circular dependency.
	runRestServer(config.Api.Rest.Port, restServer, iamGrpcAdapter, config.Api.Metrics.Enabled, config.OtelTracer.Name)

	// --- Wait for signal ---
	ch := make(chan os.Signal, 1)
	signal.Notify(ch, os.Interrupt)

	<-ch

	logger.Info().Msg("end of program...")
}
