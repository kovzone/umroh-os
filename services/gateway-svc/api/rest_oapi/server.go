package rest_oapi

import (
	"gateway-svc/api/rest_oapi/middleware"
	"gateway-svc/service"

	"github.com/gofiber/fiber/v2"
	"github.com/rs/zerolog"
	"go.opentelemetry.io/otel/trace"
)

// Server implements rest_oapi.ServerInterface by delegating to the service layer.
type Server struct {
	logger *zerolog.Logger
	tracer trace.Tracer

	diagnosticSecret string

	svc service.IService

	// iamValidator is the bearer-validation gateway keeps per ADR 0009 / F1-W7.
	// Wrapped by the bearer-auth middleware; attached to protected route groups
	// as they land (first consumer: BL-IAM-018 / S1-E-12 when iam client-facing
	// auth routes move to gateway).
	iamValidator middleware.IamValidator
}

// NewServer returns a Server that handles OpenAPI routes using the given service.
// iamValidator is the producer of identity envelopes used by the bearer-auth
// middleware; pass a *iam_grpc_adapter.Adapter in production.
func NewServer(logger *zerolog.Logger, tracer trace.Tracer, svc service.IService, iamValidator middleware.IamValidator) *Server {
	return &Server{
		logger:       logger,
		tracer:       tracer,
		svc:          svc,
		iamValidator: iamValidator,
	}
}

// RequireBearerToken returns the Fiber middleware that enforces edge auth
// per F1-W7. Apply it to protected route groups in cmd/server.go.
func (s *Server) RequireBearerToken() fiber.Handler {
	return middleware.RequireBearerToken(s.iamValidator)
}
