package main

import (
	"fmt"
	"log"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/cors"

	_ "github.com/lib/pq"

	"github.com/ashah360/fibertools"

	"github.com/ashah360/nyte-auth/internal/api/middleware"
)

func main() {

	app := fiber.New(fiber.Config{
		ErrorHandler: fibertools.ErrorHandler,
	})

	app.Use(fibertools.Recover())
	app.Use(cors.New())

	// GET handlers
	app.Get("/verify/:token", middleware.ProtectInternal, a.Verify)

	// POST handlers
	app.Post("/login" /* rateLimiter, */, a.Login)
	app.Post("/refresh", a.Refresh)
	app.Post("/refresh/:user_id", middleware.ProtectInternal, a.RefreshUserSnapshots)
	app.Post("/register", a.RegisterUser)

	// DELETE handlers
	app.Delete("/tokens/user/:user_id", middleware.ProtectInternal, a.DeleteUserTokens)

	if err := app.Listen(fmt.Sprintf(":%s", a.Config().Port)); err != nil {
		log.Println("An error occured, shutting down gracefully.", err)
		app.Shutdown()
	}
}
