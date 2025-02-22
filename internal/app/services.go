package app

import (
	"wallet-api/internal/services/user_service"
	"wallet-api/internal/services/wallet_service"
	"wallet-api/internal/services/wallet_transaction_service"
)

type Services struct {
	UserService              user_service.Contract
	WalletService            wallet_service.Contract
	WalletTransactionService wallet_transaction_service.Contract
}

func RegisterServices(r Repositories) Services {
	return Services{
		UserService:              user_service.New(r.TxRepository, r.UserRepository),
		WalletService:            wallet_service.New(r.TxRepository, r.WalletRepository),
		WalletTransactionService: wallet_transaction_service.New(r.TxRepository, r.WalletRepository, r.TransactionLedgerRepository),
	}
}
