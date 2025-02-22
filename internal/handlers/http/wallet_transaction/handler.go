package wallet_transaction_handler

import (
	"wallet-api/internal/services/wallet_transaction_service"

	"github.com/gofiber/fiber/v2"
)

type Contract interface {
	HandleDepositWallet(c *fiber.Ctx) error
	HandleWithdrawWallet(c *fiber.Ctx) error
	HandleTransferBalance(c *fiber.Ctx) error
}

type handler struct {
	WalletTransactionService wallet_transaction_service.Contract
}

func New(wallettransactionService wallet_transaction_service.Contract) Contract {
	return &handler{
		WalletTransactionService: wallettransactionService,
	}
}
