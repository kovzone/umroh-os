// Package middleware provides HTTP middleware for the REST API.
package middleware

import (
	"errors"
	"log"

	"visa-svc/util/apperrors"

	"github.com/gofiber/fiber/v2"
)

// ErrorHandler creates a middleware for centralized error handling.
// Uses apperrors.HTTPStatus and ErrorCode for consistent error envelope.
func ErrorHandler() fiber.Handler {
	return func(ctx *fiber.Ctx) error {
		// Forward to next handler
		err := ctx.Next()

		// Check if response was written
		if len(ctx.Response().Body()) == 0 {
			if err == nil {
				// No error but no response sent - this is a handler bug
				log.Printf("Warning: Handler didn't send any response for %s %s\n",
					ctx.Method(), ctx.Path())

				return ctx.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
					"error": fiber.Map{
						"code":    "INTERNAL_ERROR",
						"message": "Internal server error: no response sent",
					},
				})
			}
		} else if err == nil {
			// Response was sent and no error - all good
			return nil
		}

		// Handle fiber errors
		var fiberErr *fiber.Error
		if errors.As(err, &fiberErr) {
			code := codeFromHTTPStatus(fiberErr.Code)
			return ctx.Status(fiberErr.Code).JSON(fiber.Map{
				"error": fiber.Map{
					"code":    code,
					"message": fiberErr.Message,
				},
			})
		}

		// Default: use apperrors for domain errors or fallback to 500
		status := apperrors.HTTPStatus(err)
		code := apperrors.ErrorCode(err)
		message := err.Error()
		if message == "" {
			message = "Internal server error"
		}
		return ctx.Status(status).JSON(fiber.Map{
			"error": fiber.Map{
				"code":    code,
				"message": message,
			},
		})
	}
}

func codeFromHTTPStatus(status int) string {
	switch status {
	case 400:
		return "VALIDATION_ERROR"
	case 401:
		return "UNAUTHORIZED"
	case 403:
		return "FORBIDDEN"
	case 404:
		return "NOT_FOUND"
	case 409:
		return "CONFLICT"
	default:
		return "INTERNAL_ERROR"
	}
}
