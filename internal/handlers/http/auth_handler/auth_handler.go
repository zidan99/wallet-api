package auth_handler

import (
	"wallet-api/internal/pkg/schemas"
	"wallet-api/internal/pkg/validator"

	"github.com/gofiber/fiber/v2"
)

func (h *handler) HandleLogin(c *fiber.Ctx) error {
	var request schemas.LoginRequest
	if err := c.BodyParser(&request); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"code":    fiber.StatusBadRequest,
			"message": "Failed Parsing Body Request",
			"data":    nil,
		})
	}

	stringValidation := validator.Validate(request)
	if len(stringValidation) > 0 {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"code":    fiber.StatusBadRequest,
			"message": "Bad Request ! The request is missing a required parameter",
			"data": fiber.Map{
				"errors": stringValidation,
			},
		})
	}

	_, token, err := h.UserService.Login(c.Context(), request)
	if err != nil {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
			"code":    fiber.StatusUnauthorized,
			"message": "Failed Login ! Invalid Credentials.",
			"data":    nil,
		})
	}

	response := fiber.Map{
		"code":    fiber.StatusOK,
		"message": "Login success",
		"data": fiber.Map{
			"token": token,
		},
	}
	return c.Status(fiber.StatusOK).JSON(response)
}
