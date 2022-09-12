package app

import (
	"github.com/gofiber/fiber/v2"

	"github.com/ashah360/fibertools"

	"github.com/ashah360/nyte-auth/internal/api/token"
)

type RefreshPayload struct {
	Token string `json:"token"`
}

func (a *application) Refresh(c *fiber.Ctx) error {
	rp := new(RefreshPayload)
	if err := c.BodyParser(rp); err != nil {
		return fibertools.Message(c, fiber.StatusBadRequest, "Token required")
	}

	tp, err := token.ValidateJWT(rp.Token)
	if err != nil {
		panic(err)
	}

	rt, err := a.authService.RefreshToken(c.Context(), tp.ID, tp.Token)
	if err != nil {
		panic(err)
	}

	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"token": rt,
	})
}

func (a *application) RefreshUserSnapshots(c *fiber.Ctx) error {
	uid := c.Params("user_id")
	if uid == "" {
		return fibertools.Message(c, fiber.StatusBadRequest, "User ID is required")
	}

	c.Context().SetUserValue("TOKEN_TTL_SECONDS", a.config.TokenTTL)

	if err := a.authService.RefreshUserSnapshots(c.Context(), uid); err != nil {
		panic(err)
	}

	return c.SendStatus(fiber.StatusOK)
}
