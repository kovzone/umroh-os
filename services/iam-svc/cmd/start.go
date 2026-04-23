package main

import (
	"context"
	"fmt"
	"os"
	"os/signal"
	"time"

	"iam-svc/api/grpc_api"
	"iam-svc/service"
	"iam-svc/store/postgres_store"
	"iam-svc/util/config"
	"iam-svc/util/logging"
	"iam-svc/util/monitoring"
	"iam-svc/util/token"
	"iam-svc/util/tracing"
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

	// --- Init token maker (PASETO / JWT per config) ---
	tokenMaker, err := token.NewMaker(config.Token.Type, config.Token.Key)
	if err != nil {
		logger.Error().
			Str("op", op).
			Str("scope", "Init token maker").
			Err(err).
			Msg("")
		os.Exit(1)
	}

	// --- Validate TOTP encryption key (AES-256 requires 32 bytes) ---
	totpKey := []byte(config.Totp.EncryptionKey)
	if len(totpKey) != 32 {
		logger.Error().
			Str("op", op).
			Str("scope", "Validate TOTP encryption key").
			Int("len", len(totpKey)).
			Msg("totp.encryption_key must be exactly 32 bytes (AES-256)")
		os.Exit(1)
	}

	// --- Init service layer ---
	svc := service.NewService(
		logger,
		tracer,
		config.App.Name,
		store,
		tokenMaker,
		config.Token.AccessDuration,
		config.Token.RefreshDuration,
		config.Totp.Issuer,
		totpKey,
	)

	// --- Init API layers ---
	// Per ADR 0009 / BL-IAM-018 (S1-E-12): iam-svc is gRPC-only.
	// REST surface removed; gateway-svc proxies all client-facing auth routes.
	grpcServer := grpc_api.NewServer(logger, tracer, svc)

	// --- Run servers ---
	runGrpcServer(config.Api.Grpc.Address, grpcServer)

	// --- Wait for signal ---
	ch := make(chan os.Signal, 1)
	signal.Notify(ch, os.Interrupt)

	<-ch

	logger.Info().Msg("end of program...")
}
