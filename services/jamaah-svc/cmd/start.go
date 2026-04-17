package main

import (
	"context"
	"fmt"
	"os"
	"os/signal"
	"time"

	"jamaah-svc/api/grpc_api"
	"jamaah-svc/api/rest_oapi"
	"jamaah-svc/service"
	"jamaah-svc/store/postgres_store"
	"jamaah-svc/util/config"
	"jamaah-svc/util/logging"
	"jamaah-svc/util/monitoring"
	"jamaah-svc/util/tracing"
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

	// --- Start DB pool metrics collector (when metrics enabled) ---
	if config.Api.Metrics.Enabled {
		go monitoring.RegisterDBPoolStats(context.Background(), func() monitoring.DBPoolStats {
			s := postgresPool.Stat()
			return monitoring.DBPoolStats{
				Acquired: s.AcquiredConns(),
				Idle:     s.IdleConns(),
				Total:    s.TotalConns(),
			}
		}, 10*time.Second)
	}

	// --- Init service layer ---
	svc := service.NewService(logger, tracer, config.App.Name, store)

	// --- Init API layers ---
	restServer := rest_oapi.NewServer(logger, tracer, svc)
	grpcServer := grpc_api.NewServer(logger, tracer, svc)

	// --- Run servers ---
	runRestServer(config.Api.Rest.Port, restServer, config.Api.Metrics.Enabled)
	runGrpcServer(config.Api.Grpc.Address, grpcServer)

	// --- Wait for signal ---
	ch := make(chan os.Signal, 1)
	signal.Notify(ch, os.Interrupt)

	<-ch

	logger.Info().Msg("end of program...")
}
