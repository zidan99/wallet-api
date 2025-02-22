package health_check_handler

import "github.com/gofiber/fiber/v2"

type Contract interface {
	HandleHealthCheck(c *fiber.Ctx) error
}

type handler struct{}

func New() Contract {
	return &handler{}
}
