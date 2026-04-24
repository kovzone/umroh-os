package main

import (
	"context"
	"fmt"
	"os"
	"os/signal"

	"gateway-svc/adapter/catalog_grpc_adapter"
	"gateway-svc/adapter/finance_grpc_adapter"
	"gateway-svc/adapter/health_check_adapter"
	"gateway-svc/adapter/iam_grpc_adapter"
	"gateway-svc/api/rest_oapi"
	"gateway-svc/service"
	"gateway-svc/util/config"
	"gateway-svc/util/logging"
	"gateway-svc/util/tracing"

	"github.com/rs/zerolog"
	"go.opentelemetry.io/contrib/instrumentation/google.golang.org/grpc/otelgrpc"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

// dialBackend opens a plaintext gRPC client to the given target and attaches
// the OTel stats handler so the current trace context propagates to the
// backend. Used for every dialed backend the gateway carries — the edge is
// responsible for this one-line plumbing (per ADR 0009).
func dialBackend(logger *zerolog.Logger, name, target string) *grpc.ClientConn {
	conn, err := grpc.NewClient(
		target,
		grpc.WithTransportCredentials(insecure.NewCredentials()),
		grpc.WithStatsHandler(otelgrpc.NewClientHandler()),
	)
	if err != nil {
		logger.Error().
			Str("scope", "Dial "+name+" gRPC").
			Str("target", target).
			Err(err).
			Msg("")
		os.Exit(1)
	}
	return conn
}

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

	// --- Dial every backend ---
	//
	// Per ADR 0009 the gateway is the only process that talks to downstream
	// backends. We hold one persistent gRPC conn per backend so (a) per-route
	// business adapters (iam_grpc_adapter, catalog_grpc_adapter,
	// finance_grpc_adapter) can multiplex over a warm conn and (b) the
	// aggregate /v1/system/backends health endpoint can fan
	// grpc.health.v1.Health.Check out to all of them without a dial on
	// every poll. Traffic stays inside the docker-compose network — insecure
	// is fine today; the trust-contract card (BL-GTW-100) adds mTLS.
	backends := []struct {
		name   string
		target string
	}{
		{"iam-svc", config.External.IamSvc.GrpcTarget},
		{"catalog-svc", config.External.CatalogSvc.GrpcTarget},
		{"booking-svc", config.External.BookingSvc.GrpcTarget},
		{"jamaah-svc", config.External.JamaahSvc.GrpcTarget},
		{"payment-svc", config.External.PaymentSvc.GrpcTarget},
		{"visa-svc", config.External.VisaSvc.GrpcTarget},
		{"ops-svc", config.External.OpsSvc.GrpcTarget},
		{"logistics-svc", config.External.LogisticsSvc.GrpcTarget},
		{"finance-svc", config.External.FinanceSvc.GrpcTarget},
		{"crm-svc", config.External.CrmSvc.GrpcTarget},
	}

	conns := make(map[string]*grpc.ClientConn, len(backends))
	healthBackends := make([]health_check_adapter.Backend, 0, len(backends))
	for _, b := range backends {
		conn := dialBackend(logger, b.name, b.target)
		conns[b.name] = conn
		healthBackends = append(healthBackends, health_check_adapter.Backend{Name: b.name, Conn: conn})
		defer func(name string, c *grpc.ClientConn) {
			if err := c.Close(); err != nil {
				logger.Error().Err(err).Msgf("close %s gRPC conn", name)
			}
		}(b.name, conn)
	}

	// --- Per-route business adapters over the shared conns ---
	iamGrpcAdapter := iam_grpc_adapter.NewAdapter(logger, tracer, conns["iam-svc"])
	catalogGrpcAdapter := catalog_grpc_adapter.NewAdapter(logger, tracer, conns["catalog-svc"])
	financeGrpcAdapter := finance_grpc_adapter.NewAdapter(logger, tracer, conns["finance-svc"])
	healthCheckAdapter := health_check_adapter.NewAdapter(logger, tracer, healthBackends)

	// --- Init service layer ---
	//
	// Every backend adapter is gRPC-based after BL-IAM-019 / S1-E-14; no REST
	// adapter remains on the gateway side. iam-svc retired its REST port in
	// BL-IAM-018 / S1-E-12; the seven pure-scaffold backends retired theirs in
	// BL-REFACTOR-002..008 / S1-E-13; catalog-svc did the same in G7; finance-svc
	// did the same in this card.
	svc := service.NewService(service.NewServiceParams{
		Logger:      logger,
		Tracer:      tracer,
		AppName:     config.App.Name,
		IamGrpc:     iamGrpcAdapter,
		CatalogGrpc: catalogGrpcAdapter,
		FinanceGrpc: financeGrpcAdapter,
		HealthCheck: healthCheckAdapter,
	})

	// --- Init API layer (REST only — gateway is the edge proxy, no gRPC server) ---
	// Same iam_grpc_adapter instance satisfies both surfaces: bearer validation
	// (ValidateToken) and permission checks (CheckPermission).
	restServer := rest_oapi.NewServer(logger, tracer, svc, iamGrpcAdapter, iamGrpcAdapter)

	// --- Run server ---
	runRestServer(config.Api.Rest.Port, restServer, config.Api.Metrics.Enabled, config.OtelTracer.Name)

	// --- Wait for signal ---
	ch := make(chan os.Signal, 1)
	signal.Notify(ch, os.Interrupt)

	<-ch

	logger.Info().Msg("end of program...")
}
