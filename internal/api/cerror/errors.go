package cerror

import (
	"github.com/gofiber/fiber/v2"
)

var (
	ErrUserDoesNotExist         = fiber.NewError(fiber.StatusNotFound, "User does not exist")
	ErrInvalidAccountDetails    = fiber.NewError(fiber.StatusUnauthorized, "Invalid account details")
	ErrUserSnapshotDoesNotExist = fiber.NewError(fiber.StatusNotFound, "User token expired or not found")
	ErrMalformedToken           = fiber.NewError(fiber.StatusBadRequest, "Malformed token")
	ErrInvalidToken             = fiber.NewError(fiber.StatusBadRequest, "Invalid token")
	ErrMissingToken             = fiber.NewError(fiber.StatusBadRequest, "Missing token")
	ErrUserAlreadyExists        = fiber.NewError(fiber.StatusBadRequest, "User already exists")
)
