package router

import (
	"context"
	"wallet-api/internal/app"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/compress"
	"github.com/gofiber/fiber/v2/middleware/cors"
	"github.com/gofiber/fiber/v2/middleware/helmet"
	"github.com/google/uuid"
)

func SetupRouter(f *fiber.App, m app.Middlewares, h app.Handlers) {
	f.Use(helmet.New())
	f.Use(cors.New())
	f.Use(compress.New())

	f.Use(func(c *fiber.Ctx) error {
		requestId, err := uuid.NewV7()
		if err != nil {
			requestId = uuid.New()
		}

		ctx := context.WithValue(c.UserContext(), "requestid", requestId.String())
		c.SetUserContext(ctx)

		return c.Next()
	})

	api := f.Group("/api")
	api.Get("/health-check", h.HealthCheckHandler.HandleHealthCheck)

	api.Post("/login", h.AuthHandler.HandleLogin)

	authValid := api.Use(m.AuthMiddleware.ValidateUserToken)

	wallet := authValid.Group("/wallet")
	wallet.Get("/history", h.WalletHandler.HandleWalletHistory)
	wallet.Get("/balance", h.WalletHandler.HandleGetBalanceWallet)

	wallet.Post("/deposit", h.WalletTransactionHandler.HandleDepositWallet)
	wallet.Post("/withdraw", h.WalletTransactionHandler.HandleWithdrawWallet)
	wallet.Post("/transfer", h.WalletTransactionHandler.HandleTransferBalance)
}
