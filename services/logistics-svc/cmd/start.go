package main

import (
	"context"
	"fmt"
	"os"
	"os/signal"

	"logistics-svc/api/grpc_api"
	"logistics-svc/service"
	"logistics-svc/store/postgres_store"
	"logistics-svc/util/config"
	"logistics-svc/util/logging"
	"logistics-svc/util/monitoring"
	"logistics-svc/util/tracing"
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

	// --- Init service layer ---
	svc := service.NewService(logger, tracer, config.App.Name, store)

	// --- Init API layer (gRPC only per ADR 0009) ---
	grpcServer := grpc_api.NewServer(logger, tracer, svc)

	// --- Run server ---
	runGrpcServer(config.Api.Grpc.Address, grpcServer)

	// --- Wait for signal ---
	ch := make(chan os.Signal, 1)
	signal.Notify(ch, os.Interrupt)

	<-ch

	logger.Info().Msg("end of program...")
}
