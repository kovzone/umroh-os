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

	// permChecker is the authorization surface for gateway-side permission gates
	// (BL-IAM-019 / S1-E-14). The same *iam_grpc_adapter.Adapter instance that
	// satisfies iamValidator also satisfies PermissionChecker — both are served
	// by iam-svc's gRPC surface.
	permChecker middleware.PermissionChecker
}

// NewServer returns a Server that handles OpenAPI routes using the given service.
// iamValidator produces identity envelopes for the bearer-auth middleware;
// permChecker enforces per-route (resource, action, scope) gates. Both are
// satisfied by a *iam_grpc_adapter.Adapter in production (same instance can
// be passed twice).
func NewServer(
	logger *zerolog.Logger,
	tracer trace.Tracer,
	svc service.IService,
	iamValidator middleware.IamValidator,
	permChecker middleware.PermissionChecker,
) *Server {
	return &Server{
		logger:       logger,
		tracer:       tracer,
		svc:          svc,
		iamValidator: iamValidator,
		permChecker:  permChecker,
	}
}

// RequireBearerToken returns the Fiber middleware that enforces edge auth
// per F1-W7. Apply it to protected route groups in cmd/server.go.
func (s *Server) RequireBearerToken() fiber.Handler {
	return middleware.RequireBearerToken(s.iamValidator)
}

// RequirePermission returns the Fiber middleware that enforces a specific
// (resource, action, scope) authorization tuple via iam-svc.CheckPermission.
// Compose AFTER RequireBearerToken on a protected route — the identity
// envelope produced by bearer auth is the input to this check.
func (s *Server) RequirePermission(resource, action, scope string) fiber.Handler {
	return middleware.RequirePermission(s.permChecker, resource, action, scope)
}
