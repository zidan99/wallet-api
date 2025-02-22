package wallet_transaction_handler

import (
	"errors"
	"strings"
	"wallet-api/internal/constant/keys"
	"wallet-api/internal/models"
	"wallet-api/internal/pkg/helpers"
	"wallet-api/internal/pkg/logger"
	"wallet-api/internal/pkg/schemas"
	"wallet-api/internal/pkg/validator"

	"github.com/gofiber/fiber/v2"
	"github.com/greatcloak/decimal"
	"github.com/rotisserie/eris"
)

func (h *handler) HandleDepositWallet(c *fiber.Ctx) error {
	var user models.User
	if v, ok := c.Locals(keys.USER_AUTH_KEY).(models.User); ok {
		user = v
	}

	metadata := map[string]any{
		"original_body": c.Queries(),
	}

	var request schemas.WalletDepositRequest
	if err := c.BodyParser(&request); err != nil {
		metadata["eris"] = eris.ToJSON(err, true)

		logger.Errorf(c.Context(), metadata, err.Error())

		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"status":  "fail",
			"message": "Bad Request ! The request is missing a required parameter",
			"error":   err.Error(),
		})
	}

	metadata["request_body"] = string(helpers.ExtractStructToBytes(request))

	if errs := validator.Validate(request); len(errs) > 0 {
		err := errors.New(strings.Join(errs, ", "))

		metadata["eris"] = eris.ToJSON(err, true)

		logger.Errorf(c.UserContext(), metadata, err.Error())

		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"status":  "fail",
			"message": err.Error(),
			"error":   err.Error(),
		})
	}

	if !request.Nominal.GreaterThan(decimal.Zero) {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"status":  "fail",
			"message": "Bad Request ! The request is missing a required parameter",
			"error":   "Nominal must greater than zero",
		})
	}

	err := h.WalletTransactionService.DepositWallet(c.UserContext(), user, request)
	if err != nil {
		metadata["eris"] = eris.ToJSON(err, true)

		logger.Errorf(c.Context(), metadata, err.Error())

		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"status":  "fail",
			"message": "Deposit failed",
			"error":   err.Error(),
		})
	}

	return c.JSON(fiber.Map{
		"status":  "success",
		"message": "Deposit success",
		"data":    nil,
	})
}

func (h *handler) HandleWithdrawWallet(c *fiber.Ctx) error {
	var user models.User
	if v, ok := c.Locals(keys.USER_AUTH_KEY).(models.User); ok {
		user = v
	}

	metadata := map[string]any{
		"original_body": c.Queries(),
	}

	var request schemas.WalletWithdrawRequest
	if err := c.BodyParser(&request); err != nil {
		metadata["eris"] = eris.ToJSON(err, true)

		logger.Errorf(c.Context(), metadata, err.Error())

		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"status":  "fail",
			"message": "Bad Request ! The request is missing a required parameter",
			"error":   err.Error(),
		})
	}

	metadata["request_body"] = string(helpers.ExtractStructToBytes(request))

	if errs := validator.Validate(request); len(errs) > 0 {
		err := errors.New(strings.Join(errs, ", "))

		metadata["eris"] = eris.ToJSON(err, true)

		logger.Errorf(c.UserContext(), metadata, err.Error())

		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"status":  "fail",
			"message": err.Error(),
			"error":   err.Error(),
		})
	}

	if request.Nominal.LessThan(decimal.Zero) {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"status":  "fail",
			"message": "Bad Request ! The request is missing a required parameter",
			"error":   "Nominal must greater than zero",
		})
	}

	err := h.WalletTransactionService.WithdrawWallet(c.UserContext(), user, request)
	if err != nil {
		metadata["eris"] = eris.ToJSON(err, true)

		logger.Errorf(c.Context(), metadata, err.Error())

		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"status":  "fail",
			"message": "Withdraw failed",
			"error":   err.Error(),
		})
	}

	return c.JSON(fiber.Map{
		"status":  "success",
		"message": "Withdraw success",
		"data":    nil,
	})
}

func (h *handler) HandleTransferBalance(c *fiber.Ctx) error {
	var user models.User
	if v, ok := c.Locals(keys.USER_AUTH_KEY).(models.User); ok {
		user = v
	}

	metadata := map[string]any{
		"original_body": c.Queries(),
	}

	var request schemas.WalletTransferRequest
	if err := c.BodyParser(&request); err != nil {
		metadata["eris"] = eris.ToJSON(err, true)

		logger.Errorf(c.Context(), metadata, err.Error())

		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"status":  "fail",
			"message": "Bad Request ! The request is missing a required parameter",
			"error":   err.Error(),
		})
	}

	metadata["request_body"] = string(helpers.ExtractStructToBytes(request))

	if errs := validator.Validate(request); len(errs) > 0 {
		err := errors.New(strings.Join(errs, ", "))

		metadata["eris"] = eris.ToJSON(err, true)

		logger.Errorf(c.UserContext(), metadata, err.Error())

		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"status":  "fail",
			"message": err.Error(),
			"error":   err.Error(),
		})
	}

	if !request.Nominal.GreaterThan(decimal.Zero) {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"status":  "fail",
			"message": "Bad Request ! The request is missing a required parameter",
			"error":   "Nominal must greater than zero",
		})
	}

	err := h.WalletTransactionService.WalletTransferBalance(c.UserContext(), user, request)
	if err != nil {
		metadata["eris"] = eris.ToJSON(err, true)

		logger.Errorf(c.Context(), metadata, err.Error())

		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"status":  "fail",
			"message": "Transfer failed",
			"error":   err.Error(),
		})
	}

	return c.JSON(fiber.Map{
		"status":  "success",
		"message": "Transfer success",
		"data":    nil,
	})
}
