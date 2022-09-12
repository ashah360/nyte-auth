package app

import (
	"strings"

	"github.com/gofiber/fiber/v2"

	"github.com/ashah360/fibertools"

	"github.com/ashah360/nyte-auth/internal/api/token"
)

func (a *application) Verify(c *fiber.Ctx) error {
	t := c.Params("token")
	if t == "" {
		return fibertools.Message(c, fiber.StatusBadRequest, "Token is required")
	}

	d := c.Query("snapshot")
	if !strings.EqualFold(d, "show") {
		payload, err := token.ValidateJWT(t)
		if err != nil {
			panic(err)
		}

		return c.Status(fiber.StatusOK).JSON(payload)
	}

	us, err := a.authService.VerifyJWT(c.Context(), t)
	if err != nil {
		panic(err)
	}

	return c.Status(fiber.StatusOK).JSON(us)
}
