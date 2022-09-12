package app

import (
	"github.com/ashah360/nyte-auth/internal/captcha"
	"github.com/gofiber/fiber/v2"

	"github.com/ashah360/nyte-auth/internal/api/service"
)

type Application interface {
	Config() *Config

	RegisterUser(c *fiber.Ctx) error

	Login(c *fiber.Ctx) error
	Verify(c *fiber.Ctx) error
	Refresh(c *fiber.Ctx) error

	RefreshUserSnapshots(c *fiber.Ctx) error

	DeleteUserTokens(c *fiber.Ctx) error
}

type application struct {
	config       *Config
	authService  service.AuthService
	userService  service.UserService
	recapService captcha.RecaptchaService
}

func (a *application) Config() *Config {
	return a.config
}

func NewApplication(as service.AuthService, us service.UserService, rs captcha.RecaptchaService, cfg *Config) Application {
	return &application{
		config:       cfg,
		userService:  us,
		recapService: rs,
		authService:  as,
	}
}
