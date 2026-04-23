package main

import (
	"context"
	"fmt"
	"os"
	"os/signal"

	"catalog-svc/adapter/iam_grpc_adapter"
	"catalog-svc/api/grpc_api"
	"catalog-svc/service"
	"catalog-svc/store/postgres_store"
	"catalog-svc/util/config"
	"catalog-svc/util/logging"
	"catalog-svc/util/monitoring"
	"catalog-svc/util/tracing"

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
	cleanupTracer, err := tracing.InitTracer(config.OtelTracer.Name, config.OtelTracer.Endpoint)
	if err != nil {
		logger.Error().
			Str("op", op).
			Str("scope", "Init otel tracer").
			Err(err).
			Msg("")
	}
	defer func() {
		if err := cleanupTracer(context.Background()); err != nil {
			logger.Error().
				Str("op", op).
				Str("scope", "Cleanup otel tracer").
				Err(err).
				Msg("")
		}
	}()

	tracer := tracing.GetTracer(config.OtelTracer.Name)

	// --- Init otel meter (OTLP push → otel-collector → Prometheus exporter) ---
	cleanupMeter, err := monitoring.InitMeter(config.App.Name, config.OtelTracer.Endpoint)
	if err != nil {
		logger.Error().
			Str("op", op).
			Str("scope", "Init otel meter").
			Err(err).
			Msg("")
		os.Exit(1)
	}
	defer func() {
		if err := cleanupMeter(context.Background()); err != nil {
			logger.Error().
				Str("op", op).
				Str("scope", "Cleanup otel meter").
				Err(err).
				Msg("")
		}
	}()

	logger.Info().
		Str("op", op).
		Str("config", fmt.Sprintf("%+v", config)).
		Msg(fmt.Sprintf("Starting '%s' service ...", config.App.Name))

	// --- Create postgres pool ---
	postgresPool, err := createPostgresPoolWithRetry(config.Store.Postgres)
	if err != nil {
		logger.Error().
			Str("op", op).
			Str("scope", "Create postgres pool").
			Err(err).
			Msg("")
		os.Exit(1)
	}

	// --- Init store layer ---
	store := postgres_store.NewStore(logger, tracer, postgresPool)

	// --- Register DB pool metrics (OTel observable gauges, pulled on each scrape) ---
	if err := monitoring.RegisterDBPoolStats(func() monitoring.DBPoolStats {
		s := postgresPool.Stat()
		return monitoring.DBPoolStats{
			Acquired: s.AcquiredConns(),
			Idle:     s.IdleConns(),
			Total:    s.TotalConns(),
		}
	}); err != nil {
		logger.Error().
			Str("op", op).
			Str("scope", "Register DB pool stats").
			Err(err).
			Msg("")
		os.Exit(1)
	}

	// --- Dial iam-svc gRPC (S1-E-07 / BL-CAT-014) ---
	//
	// catalog-svc needs to call iam-svc to ValidateToken + CheckPermission for
	// every staff catalog write RPC. Traffic stays inside the docker-compose
	// network — insecure credentials are fine for the pilot.
	iamConn, err := grpc.NewClient(
		config.Iam.GrpcTarget,
		grpc.WithTransportCredentials(insecure.NewCredentials()),
	)
	if err != nil {
		logger.Error().
			Str("op", op).
			Str("scope", "Dial iam-svc gRPC").
			Str("target", config.Iam.GrpcTarget).
			Err(err).
			Msg("")
		os.Exit(1)
	}
	defer func() {
		if err := iamConn.Close(); err != nil {
			logger.Error().Err(err).Msg("close iam gRPC conn")
		}
	}()
	iamAdapter := iam_grpc_adapter.NewAdapter(logger, tracer, iamConn)

	// --- Init service layer ---
	svc := service.NewService(logger, tracer, config.App.Name, store)

	// --- Init API layer (gRPC only per ADR 0009 / BL-REFACTOR-001 / S1-E-11) ---
	// The REST server has been retired; all external reads go through gateway-svc.
	// The iamAdapter is passed to the gRPC server so catalog write RPCs can gate
	// on catalog.package.manage permission (S1-E-07 / BL-CAT-014).
	grpcServer := grpc_api.NewServer(logger, tracer, svc, iamAdapter)

	// --- Run server (gRPC only) ---
	runGrpcServer(config.Api.Grpc.Address, grpcServer)

	// --- Wait for signal ---
	ch := make(chan os.Signal, 1)
	signal.Notify(ch, os.Interrupt)

	<-ch

	logger.Info().Msg("end of program...")
}
