package main

import (
	"context"
	"fmt"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"payment-svc/adapter/booking_grpc_adapter"
	"payment-svc/adapter/gateway"
	"payment-svc/api/grpc_api"
	"payment-svc/api/http_api"
	"payment-svc/service"
	"payment-svc/store/postgres_store"
	"payment-svc/util/config"
	"payment-svc/util/logging"
	"payment-svc/util/monitoring"
	"payment-svc/util/tracing"

	"github.com/rs/zerolog"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

func start() {
	const op = "main.start"

	// --- Init logger ---
	rootLogger := logging.NewLogger(logging.Options{
		Level:    logging.LevelDebug,
		TimeZone: logging.TZWIB,
	})
	logger := &rootLogger

	// --- Load config ---
	cfg, err := config.LoadConfig(".")
	if err != nil {
		logger.Error().
			Str("op", op).
			Str("scope", "Load config").
			Err(err).
			Msg("")
		os.Exit(1)
	}

	// --- Guard: MOCK_GATEWAY must be false in production ---
	// Per S2-J-04 contract: refuse to start if ENV=production && MOCK_GATEWAY=true.
	if cfg.Gateway.MockGateway {
		envName := os.Getenv("ENV")
		if envName == "production" || envName == "prod" {
			// Per S2-J-04: payment-svc MUST refuse to start if ENV=production && MOCK_GATEWAY=true.
			logger.Error().
				Str("op", op).
				Msg("MOCK_GATEWAY=true is forbidden in production (ENV=production). Refusing to start.")
			os.Exit(1)
		}
		logger.Warn().
			Str("op", op).
			Msg("MOCK_GATEWAY=true — using mock gateway adapter. DO NOT use this in production.")
	}

	// --- Init OTel tracer ---
	cleanupTracer, err := tracing.InitTracer(cfg.OtelTracer.Name, cfg.OtelTracer.Endpoint)
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

	tracer := tracing.GetTracer(cfg.OtelTracer.Name)

	// --- Init OTel meter ---
	cleanupMeter, err := monitoring.InitMeter(cfg.App.Name, cfg.OtelTracer.Endpoint)
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
		Str("app_name", cfg.App.Name).
		Str("grpc_addr", cfg.Api.Grpc.Address).
		Str("webhook_addr", cfg.Api.Http.WebhookAddress).
		Bool("mock_gateway", cfg.Gateway.MockGateway).
		Msg(fmt.Sprintf("Starting '%s' service ...", cfg.App.Name))

	// --- Create postgres pool ---
	postgresPool, err := createPostgresPoolWithRetry(cfg.Store.Postgres)
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

	// --- Register DB pool metrics ---
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

	// --- Init gateway adapters ---
	// Gateway selection per Q013: Midtrans primary, Xendit fallback.
	// MOCK_GATEWAY=true → mock adapter used for all VA issuance (dev only).
	// TODO(REAL_API_KEY): populate MIDTRANS_SERVER_KEY and XENDIT_CALLBACK_TOKEN + XENDIT_SECRET_KEY
	//   in the deployment environment to enable real payment processing.
	var (
		primaryGateway gateway.GatewayAdapter // Midtrans
		xenditGateway  gateway.GatewayAdapter // Xendit
		mockGw         gateway.GatewayAdapter // mock (MOCK_GATEWAY=true)
	)

	mockGw = gateway.NewMockAdapter()

	if cfg.Gateway.MockGateway {
		// In mock mode, primary and xendit adapters are set to mock so
		// selectGateway() always returns the mock regardless of preference.
		primaryGateway = mockGw
		xenditGateway = mockGw
		logger.Warn().Str("op", op).Msg("gateway adapters: ALL set to mock (MOCK_GATEWAY=true)")
	} else {
		// TODO(REAL_API_KEY): non-mock production path.
		// These stubs will return errors until real keys are configured.
		primaryGateway = gateway.NewMidtransAdapter(
			cfg.Gateway.MidtransServerKey,
			cfg.Gateway.MidtransBaseURL,
		)
		xenditGateway = gateway.NewXenditAdapter(
			cfg.Gateway.XenditCallbackToken,
			cfg.Gateway.XenditSecretKey,
			cfg.Gateway.XenditBaseURL,
		)
		logger.Info().Str("op", op).Msg("gateway adapters: Midtrans (primary) + Xendit (fallback)")
	}

	// --- Dial booking-svc gRPC ---
	// payment-svc calls booking-svc.MarkBookingPaid after each successful webhook.
	// Per ADR-0006: direct gRPC, no event bus.
	bookingAddr := cfg.Services.BookingSvcAddr
	if bookingAddr == "" {
		bookingAddr = "booking-svc:50051"
	}
	bookingConn, bookingDialErr := dialGRPC(logger, bookingAddr)
	var bookingAdapter *booking_grpc_adapter.Adapter
	if bookingDialErr != nil {
		logger.Warn().
			Str("op", op).
			Str("addr", bookingAddr).
			Err(bookingDialErr).
			Msg("booking-svc dial failed — MarkBookingPaid will be skipped (reconcile will catch up)")
		// Non-fatal: reconciliation cron will catch orphaned invoice updates.
	} else {
		bookingAdapter = booking_grpc_adapter.NewAdapter(logger, tracer, bookingConn)
		logger.Info().Str("op", op).Str("addr", bookingAddr).Msg("booking-svc gRPC connected")
	}

	// --- Init service layer (PaymentService — full F5 implementation) ---
	paymentSvc := service.NewPaymentService(service.PaymentServiceConfig{
		Logger:         logger,
		Tracer:         tracer,
		AppName:        cfg.App.Name,
		Store:          store,
		PrimaryGateway: primaryGateway,
		XenditGateway:  xenditGateway,
		MockGateway:    mockGw,
		BookingAdapter: bookingAdapter,
		IAMAudit:       nil, // wired in S3 when iam_grpc_adapter is ready
	})

	// --- Start reconciliation cron (S2-E-03) ---
	// Runs every 1h in a goroutine. Catches missed webhooks + expired VAs.
	// Cancelled when the root context is cancelled (on SIGTERM).
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	// StartReconcileCron starts the internal goroutine and returns immediately.
	paymentSvc.StartReconcileCron(ctx)
	logger.Info().Str("op", op).Msg("reconciliation cron started (interval=1h)")

	// --- Start HTTP webhook listener (S2-E-02) ---
	// Separate port from gRPC. Internal only — NOT exposed to the internet directly.
	// nginx/LB forwards POST /v1/webhooks/* from gateway's IP range to this port.
	webhookAddr := cfg.Api.Http.WebhookAddress
	if webhookAddr == "" {
		webhookAddr = "0.0.0.0:50065"
	}
	webhookHandler := http_api.NewWebhookHandler(logger, paymentSvc, cfg.Gateway.MockGateway)
	webhookMux := http.NewServeMux()
	webhookHandler.RegisterRoutes(webhookMux)
	webhookServer := &http.Server{
		Addr:         webhookAddr,
		Handler:      webhookMux,
		ReadTimeout:  10 * time.Second,
		WriteTimeout: 30 * time.Second,
		IdleTimeout:  60 * time.Second,
	}
	go func() {
		logger.Info().Str("op", op).Str("addr", webhookAddr).Msg("HTTP webhook listener starting")
		if err := webhookServer.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			logger.Error().Err(err).Str("op", op).Msg("HTTP webhook server failed")
		}
	}()

	// --- Init gRPC API layer ---
	grpcServer := grpc_api.NewServer(logger, tracer, paymentSvc)

	// --- Run gRPC server ---
	runGrpcServer(cfg.Api.Grpc.Address, grpcServer)

	// --- Wait for signal ---
	ch := make(chan os.Signal, 1)
	signal.Notify(ch, os.Interrupt, syscall.SIGTERM)

	sig := <-ch
	logger.Info().Str("signal", sig.String()).Msg("shutdown signal received")

	// Graceful shutdown of HTTP webhook server.
	shutdownCtx, shutdownCancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer shutdownCancel()
	if err := webhookServer.Shutdown(shutdownCtx); err != nil {
		logger.Error().Err(err).Str("op", op).Msg("webhook server shutdown error")
	}

	logger.Info().Msg("payment-svc stopped")
}

// dialGRPC dials a gRPC server with insecure credentials (mTLS to be added in S3).
func dialGRPC(logger *zerolog.Logger, addr string) (*grpc.ClientConn, error) {
	conn, err := grpc.NewClient(addr,
		grpc.WithTransportCredentials(insecure.NewCredentials()),
	)
	if err != nil {
		return nil, fmt.Errorf("grpc.NewClient(%q): %w", addr, err)
	}
	logger.Debug().Str("addr", addr).Msg("gRPC client created")
	return conn, nil
}
