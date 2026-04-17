package rest_oapi

import (
	"gateway-svc/service"
	"gateway-svc/util/apperrors"

	"github.com/gofiber/fiber/v2"
)

// Liveness verifies that the service process is alive and responsive.
// Must not access any external dependencies; trivial and deterministic.
//
// No logging or tracing — probe endpoints are hit at high frequency; kept cheap.
// Liveness implements ServerInterface.
func (s *Server) Liveness(c *fiber.Ctx) error {
	result, err := s.svc.Liveness(c.Context(), &service.LivenessParams{})
	if err != nil {
		return fiber.NewError(apperrors.HTTPStatus(err), err.Error())
	}

	response := LiveResponse{
		Data: struct {
			Ok bool `json:"ok"`
		}{Ok: result.OK},
	}

	return c.Status(fiber.StatusOK).JSON(response)
}

// Readiness verifies that the service instance can handle real traffic.
// gateway-svc has no required external dependency at scaffold time; future
// iterations may extend this to ping wired backend adapters.
//
// No logging or tracing — probe endpoints are hit at high frequency; kept cheap.
// Readiness implements ServerInterface.
func (s *Server) Readiness(c *fiber.Ctx) error {
	result, err := s.svc.Readiness(c.Context(), &service.ReadinessParams{})
	if err != nil {
		return fiber.NewError(apperrors.HTTPStatus(err), err.Error())
	}

	response := ReadyResponse{
		Data: struct {
			Ok bool `json:"ok"`
		}{Ok: result.OK},
	}

	return c.Status(fiber.StatusOK).JSON(response)
}
