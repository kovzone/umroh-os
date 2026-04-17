package rest_oapi

import (
	"gateway-svc/service"

	"github.com/rs/zerolog"
	"go.opentelemetry.io/otel/trace"
)

// Server implements rest_oapi.ServerInterface by delegating to the service layer.
type Server struct {
	logger *zerolog.Logger
	tracer trace.Tracer

	diagnosticSecret string

	svc service.IService
}

// NewServer returns a Server that handles OpenAPI routes using the given service.
func NewServer(logger *zerolog.Logger, tracer trace.Tracer, svc service.IService) *Server {
	return &Server{logger: logger, tracer: tracer, svc: svc}
}
