package app

import (
	"wallet-api/internal/handlers/http/auth_handler"
	"wallet-api/internal/handlers/http/health_check_handler"
	"wallet-api/internal/handlers/http/wallet_handler"
	wallet_transaction_handler "wallet-api/internal/handlers/http/wallet_transaction"
)

type Handlers struct {
	HealthCheckHandler       health_check_handler.Contract
	AuthHandler              auth_handler.Contract
	WalletHandler            wallet_handler.Contract
	WalletTransactionHandler wallet_transaction_handler.Contract
}

func RegisterHandlers(s Services) Handlers {
	return Handlers{
		HealthCheckHandler:       health_check_handler.New(),
		AuthHandler:              auth_handler.New(s.UserService),
		WalletHandler:            wallet_handler.New(s.WalletService),
		WalletTransactionHandler: wallet_transaction_handler.New(s.WalletTransactionService),
	}
}
