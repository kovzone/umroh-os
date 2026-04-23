package grpc_api

// masters.go — gRPC handler implementations for master data CRUD RPCs.
//
// Implements pb.MastersHandler interface (hotel, airline, muthawwif, addon,
// departure pricing) on *Server.  Handlers gate on catalog.package.manage
// permission, delegate to s.svc, then map results to proto types.
//
// Wave-1A / S1-E-07 depth card.

import (
	"context"

	"catalog-svc/api/grpc_api/pb"
	"catalog-svc/service"
	"catalog-svc/util/apperrors"
	"catalog-svc/util/logging"

	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/codes"
	"google.golang.org/grpc/status"
)

// ---------------------------------------------------------------------------
// Hotel
// ---------------------------------------------------------------------------

func (s *Server) CreateHotel(ctx context.Context, req *pb.CreateHotelRequest) (*pb.CreateHotelResponse, error) {
	const op = "grpc_api.Server.CreateHotel"

	ctx, span := s.tracer.Start(ctx, op)
	defer span.End()
	span.SetAttributes(attribute.String("operation", op))

	logger := logging.LogWithTrace(ctx, s.logger)

	if req.GetUserId() == "" {
		return nil, status.Error(apperrors.GRPCCode(apperrors.ErrUnauthorized), "user_id is required")
	}
	if err := s.checkCatalogManagePermission(ctx, req.GetUserId()); err != nil {
		return nil, err
	}
	if req.GetName() == "" {
		return nil, status.Error(apperrors.GRPCCode(apperrors.ErrValidation), "name is required")
	}

	result, err := s.svc.CreateHotel(ctx, &service.CreateHotelParams{
		UserID:           req.GetUserId(),
		Name:             req.GetName(),
		City:             req.GetCity(),
		StarRating:       int(req.GetStarRating()),
		WalkingDistanceM: int(req.GetWalkingDistanceM()),
	})
	if err != nil {
		logger.Warn().Err(err).Msg("")
		span.RecordError(err)
		span.SetStatus(codes.Error, err.Error())
		return nil, status.Error(apperrors.GRPCCode(err), err.Error())
	}

	span.SetStatus(codes.Ok, "success")
	return &pb.CreateHotelResponse{Hotel: hotelResultToProto(result)}, nil
}

func (s *Server) UpdateHotel(ctx context.Context, req *pb.UpdateHotelRequest) (*pb.UpdateHotelResponse, error) {
	const op = "grpc_api.Server.UpdateHotel"

	ctx, span := s.tracer.Start(ctx, op)
	defer span.End()
	span.SetAttributes(attribute.String("operation", op))

	logger := logging.LogWithTrace(ctx, s.logger)

	if req.GetUserId() == "" {
		return nil, status.Error(apperrors.GRPCCode(apperrors.ErrUnauthorized), "user_id is required")
	}
	if req.GetId() == "" {
		return nil, status.Error(apperrors.GRPCCode(apperrors.ErrValidation), "id is required")
	}
	if err := s.checkCatalogManagePermission(ctx, req.GetUserId()); err != nil {
		return nil, err
	}

	result, err := s.svc.UpdateHotel(ctx, &service.UpdateHotelParams{
		UserID:           req.GetUserId(),
		ID:               req.GetId(),
		Name:             req.GetName(),
		City:             req.GetCity(),
		StarRating:       int(req.GetStarRating()),
		WalkingDistanceM: int(req.GetWalkingDistanceM()),
	})
	if err != nil {
		logger.Warn().Err(err).Msg("")
		span.RecordError(err)
		span.SetStatus(codes.Error, err.Error())
		return nil, status.Error(apperrors.GRPCCode(err), err.Error())
	}

	span.SetStatus(codes.Ok, "success")
	return &pb.UpdateHotelResponse{Hotel: hotelResultToProto(result)}, nil
}

func (s *Server) DeleteHotel(ctx context.Context, req *pb.DeleteHotelRequest) (*pb.DeleteHotelResponse, error) {
	const op = "grpc_api.Server.DeleteHotel"

	ctx, span := s.tracer.Start(ctx, op)
	defer span.End()
	span.SetAttributes(attribute.String("operation", op))

	logger := logging.LogWithTrace(ctx, s.logger)

	if req.GetUserId() == "" {
		return nil, status.Error(apperrors.GRPCCode(apperrors.ErrUnauthorized), "user_id is required")
	}
	if req.GetId() == "" {
		return nil, status.Error(apperrors.GRPCCode(apperrors.ErrValidation), "id is required")
	}
	if err := s.checkCatalogManagePermission(ctx, req.GetUserId()); err != nil {
		return nil, err
	}

	err := s.svc.DeleteHotel(ctx, &service.DeleteMasterParams{
		UserID: req.GetUserId(),
		ID:     req.GetId(),
	})
	if err != nil {
		logger.Warn().Err(err).Msg("")
		span.RecordError(err)
		span.SetStatus(codes.Error, err.Error())
		return nil, status.Error(apperrors.GRPCCode(err), err.Error())
	}

	span.SetStatus(codes.Ok, "success")
	return &pb.DeleteHotelResponse{Ok: true}, nil
}

func (s *Server) ListHotels(ctx context.Context, req *pb.ListHotelsRequest) (*pb.ListHotelsResponse, error) {
	const op = "grpc_api.Server.ListHotels"

	ctx, span := s.tracer.Start(ctx, op)
	defer span.End()
	span.SetAttributes(attribute.String("operation", op))

	logger := logging.LogWithTrace(ctx, s.logger)

	result, err := s.svc.ListHotels(ctx, &service.ListMasterParams{
		UserID: req.GetUserId(),
		Cursor: req.GetCursor(),
		Limit:  int(req.GetLimit()),
	})
	if err != nil {
		logger.Warn().Err(err).Msg("")
		span.RecordError(err)
		span.SetStatus(codes.Error, err.Error())
		return nil, status.Error(apperrors.GRPCCode(err), err.Error())
	}

	hotels := make([]*pb.Hotel, 0, len(result.Hotels))
	for _, h := range result.Hotels {
		hotels = append(hotels, hotelResultToProto(h))
	}

	span.SetStatus(codes.Ok, "success")
	return &pb.ListHotelsResponse{
		Hotels:  hotels,
		HasMore: result.HasMore,
		Cursor:  result.Cursor,
	}, nil
}

// ---------------------------------------------------------------------------
// Airline
// ---------------------------------------------------------------------------

func (s *Server) CreateAirline(ctx context.Context, req *pb.CreateAirlineRequest) (*pb.CreateAirlineResponse, error) {
	const op = "grpc_api.Server.CreateAirline"

	ctx, span := s.tracer.Start(ctx, op)
	defer span.End()
	span.SetAttributes(attribute.String("operation", op))

	logger := logging.LogWithTrace(ctx, s.logger)

	if req.GetUserId() == "" {
		return nil, status.Error(apperrors.GRPCCode(apperrors.ErrUnauthorized), "user_id is required")
	}
	if req.GetCode() == "" || req.GetName() == "" {
		return nil, status.Error(apperrors.GRPCCode(apperrors.ErrValidation), "code and name are required")
	}
	if err := s.checkCatalogManagePermission(ctx, req.GetUserId()); err != nil {
		return nil, err
	}

	result, err := s.svc.CreateAirline(ctx, &service.CreateAirlineParams{
		UserID:       req.GetUserId(),
		Code:         req.GetCode(),
		Name:         req.GetName(),
		OperatorKind: req.GetOperatorKind(),
	})
	if err != nil {
		logger.Warn().Err(err).Msg("")
		span.RecordError(err)
		span.SetStatus(codes.Error, err.Error())
		return nil, status.Error(apperrors.GRPCCode(err), err.Error())
	}

	span.SetStatus(codes.Ok, "success")
	return &pb.CreateAirlineResponse{Airline: airlineResultToProto(result)}, nil
}

func (s *Server) UpdateAirline(ctx context.Context, req *pb.UpdateAirlineRequest) (*pb.UpdateAirlineResponse, error) {
	const op = "grpc_api.Server.UpdateAirline"

	ctx, span := s.tracer.Start(ctx, op)
	defer span.End()
	span.SetAttributes(attribute.String("operation", op))

	logger := logging.LogWithTrace(ctx, s.logger)

	if req.GetUserId() == "" {
		return nil, status.Error(apperrors.GRPCCode(apperrors.ErrUnauthorized), "user_id is required")
	}
	if req.GetId() == "" {
		return nil, status.Error(apperrors.GRPCCode(apperrors.ErrValidation), "id is required")
	}
	if err := s.checkCatalogManagePermission(ctx, req.GetUserId()); err != nil {
		return nil, err
	}

	result, err := s.svc.UpdateAirline(ctx, &service.UpdateAirlineParams{
		UserID: req.GetUserId(),
		ID:     req.GetId(),
		Code:   req.GetCode(),
		Name:   req.GetName(),
	})
	if err != nil {
		logger.Warn().Err(err).Msg("")
		span.RecordError(err)
		span.SetStatus(codes.Error, err.Error())
		return nil, status.Error(apperrors.GRPCCode(err), err.Error())
	}

	span.SetStatus(codes.Ok, "success")
	return &pb.UpdateAirlineResponse{Airline: airlineResultToProto(result)}, nil
}

func (s *Server) DeleteAirline(ctx context.Context, req *pb.DeleteAirlineRequest) (*pb.DeleteAirlineResponse, error) {
	const op = "grpc_api.Server.DeleteAirline"

	ctx, span := s.tracer.Start(ctx, op)
	defer span.End()
	span.SetAttributes(attribute.String("operation", op))

	logger := logging.LogWithTrace(ctx, s.logger)

	if req.GetUserId() == "" {
		return nil, status.Error(apperrors.GRPCCode(apperrors.ErrUnauthorized), "user_id is required")
	}
	if req.GetId() == "" {
		return nil, status.Error(apperrors.GRPCCode(apperrors.ErrValidation), "id is required")
	}
	if err := s.checkCatalogManagePermission(ctx, req.GetUserId()); err != nil {
		return nil, err
	}

	err := s.svc.DeleteAirline(ctx, &service.DeleteMasterParams{
		UserID: req.GetUserId(),
		ID:     req.GetId(),
	})
	if err != nil {
		logger.Warn().Err(err).Msg("")
		span.RecordError(err)
		span.SetStatus(codes.Error, err.Error())
		return nil, status.Error(apperrors.GRPCCode(err), err.Error())
	}

	span.SetStatus(codes.Ok, "success")
	return &pb.DeleteAirlineResponse{Ok: true}, nil
}

func (s *Server) ListAirlines(ctx context.Context, req *pb.ListAirlinesRequest) (*pb.ListAirlinesResponse, error) {
	const op = "grpc_api.Server.ListAirlines"

	ctx, span := s.tracer.Start(ctx, op)
	defer span.End()
	span.SetAttributes(attribute.String("operation", op))

	logger := logging.LogWithTrace(ctx, s.logger)

	result, err := s.svc.ListAirlines(ctx, &service.ListMasterParams{
		UserID: req.GetUserId(),
		Cursor: req.GetCursor(),
		Limit:  int(req.GetLimit()),
	})
	if err != nil {
		logger.Warn().Err(err).Msg("")
		span.RecordError(err)
		span.SetStatus(codes.Error, err.Error())
		return nil, status.Error(apperrors.GRPCCode(err), err.Error())
	}

	airlines := make([]*pb.Airline, 0, len(result.Airlines))
	for _, a := range result.Airlines {
		airlines = append(airlines, airlineResultToProto(a))
	}

	span.SetStatus(codes.Ok, "success")
	return &pb.ListAirlinesResponse{
		Airlines: airlines,
		HasMore:  result.HasMore,
		Cursor:   result.Cursor,
	}, nil
}

// ---------------------------------------------------------------------------
// Muthawwif
// ---------------------------------------------------------------------------

func (s *Server) CreateMuthawwif(ctx context.Context, req *pb.CreateMuthawwifRequest) (*pb.CreateMuthawwifResponse, error) {
	const op = "grpc_api.Server.CreateMuthawwif"

	ctx, span := s.tracer.Start(ctx, op)
	defer span.End()
	span.SetAttributes(attribute.String("operation", op))

	logger := logging.LogWithTrace(ctx, s.logger)

	if req.GetUserId() == "" {
		return nil, status.Error(apperrors.GRPCCode(apperrors.ErrUnauthorized), "user_id is required")
	}
	if req.GetName() == "" {
		return nil, status.Error(apperrors.GRPCCode(apperrors.ErrValidation), "name is required")
	}
	if err := s.checkCatalogManagePermission(ctx, req.GetUserId()); err != nil {
		return nil, err
	}

	result, err := s.svc.CreateMuthawwif(ctx, &service.CreateMuthawwifParams{
		UserID:      req.GetUserId(),
		Name:        req.GetName(),
		PortraitUrl: req.GetPortraitUrl(),
	})
	if err != nil {
		logger.Warn().Err(err).Msg("")
		span.RecordError(err)
		span.SetStatus(codes.Error, err.Error())
		return nil, status.Error(apperrors.GRPCCode(err), err.Error())
	}

	span.SetStatus(codes.Ok, "success")
	return &pb.CreateMuthawwifResponse{Muthawwif: muthawwifResultToProto(result)}, nil
}

func (s *Server) UpdateMuthawwif(ctx context.Context, req *pb.UpdateMuthawwifRequest) (*pb.UpdateMuthawwifResponse, error) {
	const op = "grpc_api.Server.UpdateMuthawwif"

	ctx, span := s.tracer.Start(ctx, op)
	defer span.End()
	span.SetAttributes(attribute.String("operation", op))

	logger := logging.LogWithTrace(ctx, s.logger)

	if req.GetUserId() == "" {
		return nil, status.Error(apperrors.GRPCCode(apperrors.ErrUnauthorized), "user_id is required")
	}
	if req.GetId() == "" {
		return nil, status.Error(apperrors.GRPCCode(apperrors.ErrValidation), "id is required")
	}
	if err := s.checkCatalogManagePermission(ctx, req.GetUserId()); err != nil {
		return nil, err
	}

	result, err := s.svc.UpdateMuthawwif(ctx, &service.UpdateMuthawwifParams{
		UserID:      req.GetUserId(),
		ID:          req.GetId(),
		Name:        req.GetName(),
		PortraitUrl: req.GetPortraitUrl(),
	})
	if err != nil {
		logger.Warn().Err(err).Msg("")
		span.RecordError(err)
		span.SetStatus(codes.Error, err.Error())
		return nil, status.Error(apperrors.GRPCCode(err), err.Error())
	}

	span.SetStatus(codes.Ok, "success")
	return &pb.UpdateMuthawwifResponse{Muthawwif: muthawwifResultToProto(result)}, nil
}

func (s *Server) DeleteMuthawwif(ctx context.Context, req *pb.DeleteMuthawwifRequest) (*pb.DeleteMuthawwifResponse, error) {
	const op = "grpc_api.Server.DeleteMuthawwif"

	ctx, span := s.tracer.Start(ctx, op)
	defer span.End()
	span.SetAttributes(attribute.String("operation", op))

	logger := logging.LogWithTrace(ctx, s.logger)

	if req.GetUserId() == "" {
		return nil, status.Error(apperrors.GRPCCode(apperrors.ErrUnauthorized), "user_id is required")
	}
	if req.GetId() == "" {
		return nil, status.Error(apperrors.GRPCCode(apperrors.ErrValidation), "id is required")
	}
	if err := s.checkCatalogManagePermission(ctx, req.GetUserId()); err != nil {
		return nil, err
	}

	err := s.svc.DeleteMuthawwif(ctx, &service.DeleteMasterParams{
		UserID: req.GetUserId(),
		ID:     req.GetId(),
	})
	if err != nil {
		logger.Warn().Err(err).Msg("")
		span.RecordError(err)
		span.SetStatus(codes.Error, err.Error())
		return nil, status.Error(apperrors.GRPCCode(err), err.Error())
	}

	span.SetStatus(codes.Ok, "success")
	return &pb.DeleteMuthawwifResponse{Ok: true}, nil
}

func (s *Server) ListMuthawwif(ctx context.Context, req *pb.ListMuthawwifRequest) (*pb.ListMuthawwifResponse, error) {
	const op = "grpc_api.Server.ListMuthawwif"

	ctx, span := s.tracer.Start(ctx, op)
	defer span.End()
	span.SetAttributes(attribute.String("operation", op))

	logger := logging.LogWithTrace(ctx, s.logger)

	result, err := s.svc.ListMuthawwif(ctx, &service.ListMasterParams{
		UserID: req.GetUserId(),
		Cursor: req.GetCursor(),
		Limit:  int(req.GetLimit()),
	})
	if err != nil {
		logger.Warn().Err(err).Msg("")
		span.RecordError(err)
		span.SetStatus(codes.Error, err.Error())
		return nil, status.Error(apperrors.GRPCCode(err), err.Error())
	}

	muthawwifs := make([]*pb.Muthawwif, 0, len(result.Muthawwif))
	for _, m := range result.Muthawwif {
		muthawwifs = append(muthawwifs, muthawwifResultToProto(m))
	}

	span.SetStatus(codes.Ok, "success")
	return &pb.ListMuthawwifResponse{
		Muthawwif: muthawwifs,
		HasMore:   result.HasMore,
		Cursor:    result.Cursor,
	}, nil
}

// ---------------------------------------------------------------------------
// Addon
// ---------------------------------------------------------------------------

func (s *Server) CreateAddon(ctx context.Context, req *pb.CreateAddonRequest) (*pb.CreateAddonResponse, error) {
	const op = "grpc_api.Server.CreateAddon"

	ctx, span := s.tracer.Start(ctx, op)
	defer span.End()
	span.SetAttributes(attribute.String("operation", op))

	logger := logging.LogWithTrace(ctx, s.logger)

	if req.GetUserId() == "" {
		return nil, status.Error(apperrors.GRPCCode(apperrors.ErrUnauthorized), "user_id is required")
	}
	if req.GetName() == "" {
		return nil, status.Error(apperrors.GRPCCode(apperrors.ErrValidation), "name is required")
	}
	if err := s.checkCatalogManagePermission(ctx, req.GetUserId()); err != nil {
		return nil, err
	}

	result, err := s.svc.CreateAddon(ctx, &service.CreateAddonParams{
		UserID:        req.GetUserId(),
		Name:          req.GetName(),
		ListAmountIDR: req.GetListAmountIdr(),
	})
	if err != nil {
		logger.Warn().Err(err).Msg("")
		span.RecordError(err)
		span.SetStatus(codes.Error, err.Error())
		return nil, status.Error(apperrors.GRPCCode(err), err.Error())
	}

	span.SetStatus(codes.Ok, "success")
	return &pb.CreateAddonResponse{Addon: addonResultToProto(result)}, nil
}

func (s *Server) UpdateAddon(ctx context.Context, req *pb.UpdateAddonRequest) (*pb.UpdateAddonResponse, error) {
	const op = "grpc_api.Server.UpdateAddon"

	ctx, span := s.tracer.Start(ctx, op)
	defer span.End()
	span.SetAttributes(attribute.String("operation", op))

	logger := logging.LogWithTrace(ctx, s.logger)

	if req.GetUserId() == "" {
		return nil, status.Error(apperrors.GRPCCode(apperrors.ErrUnauthorized), "user_id is required")
	}
	if req.GetId() == "" {
		return nil, status.Error(apperrors.GRPCCode(apperrors.ErrValidation), "id is required")
	}
	if err := s.checkCatalogManagePermission(ctx, req.GetUserId()); err != nil {
		return nil, err
	}

	result, err := s.svc.UpdateAddon(ctx, &service.UpdateAddonParams{
		UserID:        req.GetUserId(),
		ID:            req.GetId(),
		Name:          req.GetName(),
		ListAmountIDR: req.GetListAmountIdr(),
	})
	if err != nil {
		logger.Warn().Err(err).Msg("")
		span.RecordError(err)
		span.SetStatus(codes.Error, err.Error())
		return nil, status.Error(apperrors.GRPCCode(err), err.Error())
	}

	span.SetStatus(codes.Ok, "success")
	return &pb.UpdateAddonResponse{Addon: addonResultToProto(result)}, nil
}

func (s *Server) DeleteAddon(ctx context.Context, req *pb.DeleteAddonRequest) (*pb.DeleteAddonResponse, error) {
	const op = "grpc_api.Server.DeleteAddon"

	ctx, span := s.tracer.Start(ctx, op)
	defer span.End()
	span.SetAttributes(attribute.String("operation", op))

	logger := logging.LogWithTrace(ctx, s.logger)

	if req.GetUserId() == "" {
		return nil, status.Error(apperrors.GRPCCode(apperrors.ErrUnauthorized), "user_id is required")
	}
	if req.GetId() == "" {
		return nil, status.Error(apperrors.GRPCCode(apperrors.ErrValidation), "id is required")
	}
	if err := s.checkCatalogManagePermission(ctx, req.GetUserId()); err != nil {
		return nil, err
	}

	err := s.svc.DeleteAddon(ctx, &service.DeleteMasterParams{
		UserID: req.GetUserId(),
		ID:     req.GetId(),
	})
	if err != nil {
		logger.Warn().Err(err).Msg("")
		span.RecordError(err)
		span.SetStatus(codes.Error, err.Error())
		return nil, status.Error(apperrors.GRPCCode(err), err.Error())
	}

	span.SetStatus(codes.Ok, "success")
	return &pb.DeleteAddonResponse{Ok: true}, nil
}

func (s *Server) ListAddons(ctx context.Context, req *pb.ListAddonsRequest) (*pb.ListAddonsResponse, error) {
	const op = "grpc_api.Server.ListAddons"

	ctx, span := s.tracer.Start(ctx, op)
	defer span.End()
	span.SetAttributes(attribute.String("operation", op))

	logger := logging.LogWithTrace(ctx, s.logger)

	result, err := s.svc.ListAddons(ctx, &service.ListMasterParams{
		UserID: req.GetUserId(),
		Cursor: req.GetCursor(),
		Limit:  int(req.GetLimit()),
	})
	if err != nil {
		logger.Warn().Err(err).Msg("")
		span.RecordError(err)
		span.SetStatus(codes.Error, err.Error())
		return nil, status.Error(apperrors.GRPCCode(err), err.Error())
	}

	addons := make([]*pb.Addon, 0, len(result.Addons))
	for _, a := range result.Addons {
		addons = append(addons, addonResultToProto(a))
	}

	span.SetStatus(codes.Ok, "success")
	return &pb.ListAddonsResponse{
		Addons:  addons,
		HasMore: result.HasMore,
		Cursor:  result.Cursor,
	}, nil
}

// ---------------------------------------------------------------------------
// Departure Pricing
// ---------------------------------------------------------------------------

func (s *Server) SetDeparturePricing(ctx context.Context, req *pb.SetDeparturePricingRequest) (*pb.SetDeparturePricingResponse, error) {
	const op = "grpc_api.Server.SetDeparturePricing"

	ctx, span := s.tracer.Start(ctx, op)
	defer span.End()
	span.SetAttributes(attribute.String("operation", op))

	logger := logging.LogWithTrace(ctx, s.logger)

	if req.GetUserId() == "" {
		return nil, status.Error(apperrors.GRPCCode(apperrors.ErrUnauthorized), "user_id is required")
	}
	if req.GetDepartureId() == "" {
		return nil, status.Error(apperrors.GRPCCode(apperrors.ErrValidation), "departure_id is required")
	}
	if len(req.GetPricings()) == 0 {
		return nil, status.Error(apperrors.GRPCCode(apperrors.ErrValidation), "pricings must not be empty")
	}
	if err := s.checkCatalogManagePermission(ctx, req.GetUserId()); err != nil {
		return nil, err
	}

	pricings := make([]service.PricingUpsertInput, 0, len(req.GetPricings()))
	for _, p := range req.GetPricings() {
		pricings = append(pricings, service.PricingUpsertInput{
			RoomType:      p.GetRoomType(),
			ListAmountIDR: p.GetListAmountIdr(),
		})
	}

	results, err := s.svc.SetDeparturePricing(ctx, &service.SetDeparturePricingParams{
		UserID:      req.GetUserId(),
		DepartureID: req.GetDepartureId(),
		Pricings:    pricings,
	})
	if err != nil {
		logger.Warn().Err(err).Msg("")
		span.RecordError(err)
		span.SetStatus(codes.Error, err.Error())
		return nil, status.Error(apperrors.GRPCCode(err), err.Error())
	}

	protoPricings := make([]*pb.DeparturePricing, 0, len(results))
	for _, r := range results {
		protoPricings = append(protoPricings, pricingResultToProto(r))
	}

	span.SetStatus(codes.Ok, "success")
	return &pb.SetDeparturePricingResponse{Pricings: protoPricings}, nil
}

func (s *Server) GetDeparturePricing(ctx context.Context, req *pb.GetDeparturePricingRequest) (*pb.GetDeparturePricingResponse, error) {
	const op = "grpc_api.Server.GetDeparturePricing"

	ctx, span := s.tracer.Start(ctx, op)
	defer span.End()
	span.SetAttributes(attribute.String("operation", op))

	logger := logging.LogWithTrace(ctx, s.logger)

	if req.GetDepartureId() == "" {
		return nil, status.Error(apperrors.GRPCCode(apperrors.ErrValidation), "departure_id is required")
	}

	results, err := s.svc.GetDeparturePricing(ctx, &service.GetDeparturePricingParams{
		UserID:      req.GetUserId(),
		DepartureID: req.GetDepartureId(),
	})
	if err != nil {
		logger.Warn().Err(err).Msg("")
		span.RecordError(err)
		span.SetStatus(codes.Error, err.Error())
		return nil, status.Error(apperrors.GRPCCode(err), err.Error())
	}

	protoPricings := make([]*pb.DeparturePricing, 0, len(results))
	for _, r := range results {
		protoPricings = append(protoPricings, pricingResultToProto(r))
	}

	span.SetStatus(codes.Ok, "success")
	return &pb.GetDeparturePricingResponse{Pricings: protoPricings}, nil
}

// ---------------------------------------------------------------------------
// Proto conversion helpers
// ---------------------------------------------------------------------------

func hotelResultToProto(r *service.HotelResult) *pb.Hotel {
	if r == nil {
		return nil
	}
	return &pb.Hotel{
		Id:               r.ID,
		Name:             r.Name,
		City:             r.City,
		StarRating:       int32(r.StarRating),
		WalkingDistanceM: int32(r.WalkingDistanceM),
		CreatedAt:        r.CreatedAt,
		UpdatedAt:        r.UpdatedAt,
	}
}

func airlineResultToProto(r *service.AirlineResult) *pb.Airline {
	if r == nil {
		return nil
	}
	return &pb.Airline{
		Id:           r.ID,
		Code:         r.Code,
		Name:         r.Name,
		OperatorKind: r.OperatorKind,
		CreatedAt:    r.CreatedAt,
		UpdatedAt:    r.UpdatedAt,
	}
}

func muthawwifResultToProto(r *service.MuthawwifResult) *pb.Muthawwif {
	if r == nil {
		return nil
	}
	return &pb.Muthawwif{
		Id:          r.ID,
		Name:        r.Name,
		PortraitUrl: r.PortraitUrl,
		CreatedAt:   r.CreatedAt,
		UpdatedAt:   r.UpdatedAt,
	}
}

func addonResultToProto(r *service.AddonResult) *pb.Addon {
	if r == nil {
		return nil
	}
	return &pb.Addon{
		Id:            r.ID,
		Name:          r.Name,
		ListAmountIdr: r.ListAmountIDR,
		CreatedAt:     r.CreatedAt,
		UpdatedAt:     r.UpdatedAt,
	}
}

func pricingResultToProto(r *service.PricingResult) *pb.DeparturePricing {
	if r == nil {
		return nil
	}
	return &pb.DeparturePricing{
		Id:            r.ID,
		DepartureId:   r.DepartureID,
		RoomType:      r.RoomType,
		ListAmountIdr: r.ListAmount,
		CreatedAt:     r.CreatedAt,
		UpdatedAt:     r.UpdatedAt,
	}
}
