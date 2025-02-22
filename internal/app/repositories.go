package app

import (
	"wallet-api/internal/repositories/database/transaction_ledger_repository"
	"wallet-api/internal/repositories/database/tx_repository"
	"wallet-api/internal/repositories/database/user_repository"
	"wallet-api/internal/repositories/database/wallet_repository"

	"gorm.io/gorm"
)

type Repositories struct {
	TxRepository                tx_repository.Contract
	UserRepository              user_repository.Contract
	WalletRepository            wallet_repository.Contract
	TransactionLedgerRepository transaction_ledger_repository.Contract
}

func RegisterRepositories(db *gorm.DB) Repositories {
	return Repositories{
		TxRepository:                tx_repository.New(db),
		UserRepository:              user_repository.New(db),
		WalletRepository:            wallet_repository.New(db),
		TransactionLedgerRepository: transaction_ledger_repository.New(db),
	}
}
