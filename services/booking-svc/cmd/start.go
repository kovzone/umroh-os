package main

import (
	"context"
	"fmt"
	"os"
	"os/signal"

	"booking-svc/adapter/catalog_grpc_adapter"
	"booking-svc/adapter/crm_grpc_adapter"
	"booking-svc/adapter/finance_grpc_adapter"
	"booking-svc/adapter/iam_grpc_adapter"
	"booking-svc/adapter/logistics_grpc_adapter"
	"booking-svc/api/grpc_api"
	"booking-svc/service"
	"booking-svc/store/postgres_store"
	"booking-svc/util/config"
	"booking-svc/util/logging"
	"booking-svc/util/monitoring"
	"booking-svc/util/tracing"

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
	// Non-fatal: if otel-collector is unavailable at startup (e.g. ordering race
	// in docker-compose), log a warning and continue without metrics rather than
	// crashing the service. Booking functionality is not affected.
	cleanupMeter, err := monitoring.InitMeter(config.App.Name, config.OtelTracer.Endpoint)
	if err != nil {
		logger.Warn().
			Str("op", op).
			Str("scope", "Init otel meter").
			Err(err).
			Msg("metrics disabled — otel-collector unreachable at startup")
		cleanupMeter = func(_ context.Context) error { return nil }
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

	// --- Dial iam-svc gRPC (BL-IAM-004) ---
	//
	// Traffic stays inside the docker-compose network — insecure is fine today;
	// the gateway-svc hardening card adds mTLS. Unary stats handler propagates
	// the current trace context so iam-svc spans continue the booking-svc trace.
	// No booking-svc handler calls iam-svc yet; the dial is scaffolded so S1-E-03
	// (booking draft) and subsequent BL-BKG-* cards get a one-line adapter call.
	iamConn, err := grpc.NewClient(
		config.Iam.GrpcTarget,
		grpc.WithTransportCredentials(insecure.NewCredentials()),
		grpc.WithStatsHandler(otelgrpc.NewClientHandler()),
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

	// --- Dial catalog-svc gRPC (S1-E-03 / BL-BOOK-004) ---
	//
	// Used by CreateDraftBooking to validate departures (GetPackageDeparture) and
	// reserve seats atomically (ReserveSeats). Traffic stays inside the
	// docker-compose network — insecure is acceptable until gateway mTLS lands.
	catalogConn, err := grpc.NewClient(
		config.Catalog.GrpcTarget,
		grpc.WithTransportCredentials(insecure.NewCredentials()),
		grpc.WithStatsHandler(otelgrpc.NewClientHandler()),
	)
	if err != nil {
		logger.Error().
			Str("op", op).
			Str("scope", "Dial catalog-svc gRPC").
			Str("target", config.Catalog.GrpcTarget).
			Err(err).
			Msg("")
		os.Exit(1)
	}
	defer func() {
		if err := catalogConn.Close(); err != nil {
			logger.Error().Err(err).Msg("close catalog gRPC conn")
		}
	}()
	catalogAdapter := catalog_grpc_adapter.NewAdapter(logger, tracer, catalogConn)

	// --- Dial logistics-svc gRPC (S3-E-02) ---
	//
	// Used by MarkBookingPaid to trigger fulfillment task creation when a
	// booking transitions to paid_in_full. Fire-and-forget; failure is logged
	// but does not block payment-svc response.
	logisticsConn, err := grpc.NewClient(
		config.Logistics.GrpcTarget,
		grpc.WithTransportCredentials(insecure.NewCredentials()),
		grpc.WithStatsHandler(otelgrpc.NewClientHandler()),
	)
	if err != nil {
		logger.Error().
			Str("op", op).
			Str("scope", "Dial logistics-svc gRPC").
			Str("target", config.Logistics.GrpcTarget).
			Err(err).
			Msg("")
		os.Exit(1)
	}
	defer func() {
		if err := logisticsConn.Close(); err != nil {
			logger.Error().Err(err).Msg("close logistics gRPC conn")
		}
	}()
	logisticsAdapter := logistics_grpc_adapter.NewAdapter(logger, tracer, logisticsConn)

	// --- Dial finance-svc gRPC (S3-E-03) ---
	//
	// Used by MarkBookingPaid to trigger journal posting when a booking
	// transitions to paid_in_full. Fire-and-forget; failure is logged but
	// does not block payment-svc response.
	financeConn, err := grpc.NewClient(
		config.Finance.GrpcTarget,
		grpc.WithTransportCredentials(insecure.NewCredentials()),
		grpc.WithStatsHandler(otelgrpc.NewClientHandler()),
	)
	if err != nil {
		logger.Error().
			Str("op", op).
			Str("scope", "Dial finance-svc gRPC").
			Str("target", config.Finance.GrpcTarget).
			Err(err).
			Msg("")
		os.Exit(1)
	}
	defer func() {
		if err := financeConn.Close(); err != nil {
			logger.Error().Err(err).Msg("close finance gRPC conn")
		}
	}()
	financeAdapter := finance_grpc_adapter.NewAdapter(logger, tracer, financeConn)

	// --- Dial crm-svc gRPC (S4-E-02) ---
	//
	// Used by CreateDraftBooking and MarkBookingPaid to notify crm-svc of lead
	// lifecycle events. CRM calls are best-effort: failure does not block booking
	// operations. An empty crm.grpc_target in config disables the dial and the
	// service will log a warning that crmClient is nil.
	var crmAdapter *crm_grpc_adapter.Adapter
	if config.Crm.GrpcTarget != "" {
		crmConn, err := grpc.NewClient(
			config.Crm.GrpcTarget,
			grpc.WithTransportCredentials(insecure.NewCredentials()),
			grpc.WithStatsHandler(otelgrpc.NewClientHandler()),
		)
		if err != nil {
			logger.Error().
				Str("op", op).
				Str("scope", "Dial crm-svc gRPC").
				Str("target", config.Crm.GrpcTarget).
				Err(err).
				Msg("")
			os.Exit(1)
		}
		defer func() {
			if err := crmConn.Close(); err != nil {
				logger.Error().Err(err).Msg("close crm gRPC conn")
			}
		}()
		crmAdapter = crm_grpc_adapter.NewAdapter(logger, tracer, crmConn)
	} else {
		logger.Warn().
			Str("op", op).
			Msg("crm.grpc_target not configured — CRM fan-out disabled")
	}

	// --- Init service layer ---
	// Note: crmAdapter may be nil (when crm.grpc_target is not configured).
	// Passing a typed nil (*crm_grpc_adapter.Adapter) as the CrmClient interface
	// would produce a non-nil interface value, causing panics on nil dereference.
	// We use an explicit interface nil when the adapter is unset.
	var crmClient service.CrmClient
	if crmAdapter != nil {
		crmClient = crmAdapter
	}
	svc := service.NewService(logger, tracer, config.App.Name, store, iamAdapter, catalogAdapter, logisticsAdapter, financeAdapter, crmClient)

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
