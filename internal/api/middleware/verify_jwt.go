package middleware

import (
	"strings"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/utils"

	"github.com/ashah360/fibertools"

	"github.com/ashah360/nyte-auth/internal/api/token"
)

func ExtractJWT(c *fiber.Ctx) string {
	bearer := utils.UnsafeString(c.Request().Header.Peek("Authorization"))

	data := strings.Split(bearer, " ")
	if len(data) == 2 {
		return data[1]
	}

	return ""
}

func ValidateJWT(c *fiber.Ctx) error {
	t := ExtractJWT(c)
	if t == "" {
		return fibertools.Message(c, fiber.StatusUnauthorized, "Unauthorized")
	}

	p, err := token.ValidateJWT(t)
	if err != nil {
		panic(err)
	}

	c.Locals("uid", p.ID)
	c.Locals("opaque_token", p.Token)

	return c.Next()
}
