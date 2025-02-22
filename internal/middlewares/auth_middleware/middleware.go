package auth_middleware

import (
	"wallet-api/internal/services/user_service"

	"github.com/gofiber/fiber/v2"
)

type Contract interface {
	ValidateUserToken(c *fiber.Ctx) error
}

type middleware struct {
	UserService user_service.Contract
}

func New(userService user_service.Contract) Contract {
	return &middleware{
		UserService: userService,
	}
}
