package app

import (
	"github.com/ashah360/fibertools"

	"github.com/gofiber/fiber/v2"
)

type LoginPayload struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

func (a *application) Login(c *fiber.Ctx) error {
	p := new(LoginPayload)

	if err := c.BodyParser(p); err != nil {
		return fibertools.Message(c, fiber.StatusBadRequest, "Invalid login format")
	}

	c.Context().SetUserValue("TOKEN_TTL_SECONDS", a.config.TokenTTL)

	token, err := a.authService.AuthenticateUser(c.Context(), p.Email, p.Password)
	if err != nil {
		panic(err)
	}

	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"token": token,
	})
}
