package wallet_handler

import (
	"wallet-api/internal/services/wallet_service"

	"github.com/gofiber/fiber/v2"
)

type Contract interface {
	HandleWalletHistory(c *fiber.Ctx) error
	HandleGetBalanceWallet(c *fiber.Ctx) error
}

type handler struct {
	WalletService wallet_service.Contract
}

func New(walletService wallet_service.Contract) Contract {
	return &handler{
		WalletService: walletService,
	}
}
