package app

import (
	"github.com/ashah360/fibertools"
	"github.com/gofiber/fiber/v2"
)

func (a *application) DeleteUserTokens(c *fiber.Ctx) error {
	userId := c.Params("user_id")
	if userId == "" {
		return fibertools.Message(c, fiber.StatusBadRequest, "User ID required")
	}

	err := a.authService.DeleteTokensByUserID(c.Context(), userId)
	if err != nil {
		panic(err)
	}

	return c.SendStatus(fiber.StatusNoContent)
}
