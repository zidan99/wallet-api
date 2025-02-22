package wallet_handler

import (
	"wallet-api/internal/constant/keys"
	"wallet-api/internal/models"
	"wallet-api/internal/pkg/logger"
	"wallet-api/internal/pkg/schemas"

	"github.com/gofiber/fiber/v2"
	"github.com/rotisserie/eris"
)

func (h *handler) HandleWalletHistory(c *fiber.Ctx) error {
	var user models.User
	if v, ok := c.Locals(keys.USER_AUTH_KEY).(models.User); ok {
		user = v
	}

	metadata := map[string]any{
		"original_body": c.Queries(),
	}

	var filter schemas.WalletHistoryFilter
	if err := c.QueryParser(&filter); err != nil {
		metadata["eris"] = eris.ToJSON(err, true)

		logger.Errorf(c.Context(), metadata, err.Error())

		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"status":  "fail",
			"message": "Bad Request ! The request is missing a required parameter",
			"error":   err.Error(),
		})
	}

	if filter.Page < 1 {
		filter.Page = 1
	}

	if filter.Limit < 1 {
		filter.Limit = 10
	}

	data, _, err := h.WalletService.GetWalletHistory(c.Context(), filter, user)
	if err != nil {
		metadata["eris"] = eris.ToJSON(err, true)

		logger.Errorf(c.Context(), metadata, err.Error())

		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"status":  "fail",
			"message": "Get data wallet transaction failed",
			"error":   err.Error(),
		})
	}

	return c.JSON(fiber.Map{
		"status":  "success",
		"message": "Get data wallet transaction success",
		"data":    data,
	})
}

func (h *handler) HandleGetBalanceWallet(c *fiber.Ctx) error {
	var user models.User

	if v, ok := c.Locals(keys.USER_AUTH_KEY).(models.User); ok {
		user = v
	}

	metadata := map[string]any{
		"original_body": c.Queries(),
	}

	var filter schemas.WalletHistoryFilter
	if err := c.QueryParser(&filter); err != nil {
		metadata["eris"] = eris.ToJSON(err, true)

		logger.Errorf(c.Context(), metadata, err.Error())

		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"status":  "fail",
			"message": "Bad Request ! The request is missing a required parameter",
			"error":   err.Error(),
		})
	}

	wallet, err := h.WalletService.GetWalletByID(c.UserContext(), user, filter.WalletID)
	if err != nil {
		metadata["eris"] = eris.ToJSON(err, true)

		logger.Errorf(c.Context(), metadata, err.Error())

		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"status":  "fail",
			"message": "Get data wallet failed",
			"error":   err.Error(),
		})
	}

	return c.JSON(fiber.Map{
		"status":  "success",
		"message": "Get data wallet success",
		"data":    wallet.Balance.String(),
	})
}
