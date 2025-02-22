package auth_handler

import (
	"wallet-api/internal/services/user_service"

	"github.com/gofiber/fiber/v2"
)

type Contract interface {
	HandleLogin(c *fiber.Ctx) error
}

type handler struct {
	UserService user_service.Contract
}

func New(userService user_service.Contract) Contract {
	return &handler{
		UserService: userService,
	}
}
