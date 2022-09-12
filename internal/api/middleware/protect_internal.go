package middleware

import (
	"os"

	"github.com/gofiber/fiber/v2"
)

const (
	InternalHeader       = "X-Internal-Access"
	InternalHeaderEnvVar = "X_INTERNAL_ACCESS"
)

func ProtectInternal(c *fiber.Ctx) error {
	if c.Get(InternalHeader) != os.Getenv(InternalHeaderEnvVar) {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
			"status":  fiber.StatusUnauthorized,
			"message": "Unauthorized",
		})
	}

	return c.Next()
}
