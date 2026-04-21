package main

import (
	"context"
	"fmt"
	"os"
	"os/signal"

	"gateway-svc/adapter/booking_rest_adapter"
	"gateway-svc/adapter/catalog_rest_adapter"
	"gateway-svc/adapter/crm_rest_adapter"
	"gateway-svc/adapter/finance_rest_adapter"
	"gateway-svc/adapter/iam_rest_adapter"
	"gateway-svc/adapter/jamaah_rest_adapter"
	"gateway-svc/adapter/logistics_rest_adapter"
	"gateway-svc/adapter/ops_rest_adapter"
	"gateway-svc/adapter/payment_rest_adapter"
	"gateway-svc/adapter/visa_rest_adapter"
	"gateway-svc/api/rest_oapi"
	"gateway-svc/service"
	"gateway-svc/util/config"
	"gateway-svc/util/logging"
	"gateway-svc/util/tracing"
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

	// --- Init REST adapters (one per backend service) ---
	// gateway-svc has no DB and no internal store; the service layer dispatches
	// to these per-backend adapters. As real per-route methods are added on
	// each adapter (e.g. iam.GetUser, catalog.GetPackage), the gateway proxies
	// them via the same dispatch shape.
	iamAdapter := iam_rest_adapter.NewAdapter(logger, tracer, config.External.IamSvc.Address)
	catalogAdapter := catalog_rest_adapter.NewAdapter(logger, tracer, config.External.CatalogSvc.Address)
	bookingAdapter := booking_rest_adapter.NewAdapter(logger, tracer, config.External.BookingSvc.Address)
	jamaahAdapter := jamaah_rest_adapter.NewAdapter(logger, tracer, config.External.JamaahSvc.Address)
	paymentAdapter := payment_rest_adapter.NewAdapter(logger, tracer, config.External.PaymentSvc.Address)
	visaAdapter := visa_rest_adapter.NewAdapter(logger, tracer, config.External.VisaSvc.Address)
	opsAdapter := ops_rest_adapter.NewAdapter(logger, tracer, config.External.OpsSvc.Address)
	logisticsAdapter := logistics_rest_adapter.NewAdapter(logger, tracer, config.External.LogisticsSvc.Address)
	financeAdapter := finance_rest_adapter.NewAdapter(logger, tracer, config.External.FinanceSvc.Address)
	crmAdapter := crm_rest_adapter.NewAdapter(logger, tracer, config.External.CrmSvc.Address)

	// --- Init service layer ---
	svc := service.NewService(service.NewServiceParams{
		Logger:        logger,
		Tracer:        tracer,
		AppName:       config.App.Name,
		IamRest:       iamAdapter,
		CatalogRest:   catalogAdapter,
		BookingRest:   bookingAdapter,
		JamaahRest:    jamaahAdapter,
		PaymentRest:   paymentAdapter,
		VisaRest:      visaAdapter,
		OpsRest:       opsAdapter,
		LogisticsRest: logisticsAdapter,
		FinanceRest:   financeAdapter,
		CrmRest:       crmAdapter,
	})

	// --- Init API layer (REST only — gateway is the edge proxy, no gRPC server) ---
	restServer := rest_oapi.NewServer(logger, tracer, svc)

	// --- Run server ---
	runRestServer(config.Api.Rest.Port, restServer, config.Api.Metrics.Enabled, config.OtelTracer.Name)

	// --- Wait for signal ---
	ch := make(chan os.Signal, 1)
	signal.Notify(ch, os.Interrupt)

	<-ch

	logger.Info().Msg("end of program...")
}
