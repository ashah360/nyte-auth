package app

import (
	"log"
	"net/mail"
	"os"
	"strings"

	"github.com/ashah360/fibertools"
	"github.com/ashah360/nyte-auth/internal/api/cerror"
	"github.com/gofiber/fiber/v2"
)

type registerUserPayload struct {
	Email         string `json:"email"`
	Password      string `json:"password"`
	RecapResponse string `json:"g-recaptcha-response"`
}

var internalAccessHeader = os.Getenv("X_INTERNAL_ACCESS")

func (a *application) RegisterUser(c *fiber.Ctx) error {
	var p registerUserPayload

	if err := c.BodyParser(&p); err != nil {
		return fibertools.Message(c, fiber.StatusBadRequest, "Invalid registration format")
	}

	if a.config.RequireRecaptchaToRegister && fibertools.GetHeader(c, "X-Internal-Access") != internalAccessHeader {
		validRecap, err := a.recapService.ValidateRecaptcha(p.RecapResponse)
		if err != nil {
			log.Println(err)
		}

		if !validRecap {
			return fibertools.Message(c, fiber.StatusForbidden, "reCAPTCHA is not valid")
		}
	}

	// case insensitize email
	p.Email = strings.ToLower(p.Email)

	if _, err := mail.ParseAddress(p.Email); err != nil {
		return fibertools.Message(c, fiber.StatusBadRequest, "Invalid email address")
	}

	if u, err := a.userService.GetUserByEmail(c.Context(), p.Email); err != cerror.ErrUserDoesNotExist && u != nil {
		return fibertools.Message(c, fiber.StatusBadRequest, "User already exists")
	}

	if err := a.userService.CreateUser(c.Context(), p.Email, p.Password); err != nil {
		panic(err)
	}

	return fibertools.Message(c, fiber.StatusOK, "User created successfully")
}
