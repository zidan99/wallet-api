package wallet_transaction_service

import (
	"context"
	"wallet-api/internal/models"
	"wallet-api/internal/pkg/schemas"
	"wallet-api/internal/repositories/database/transaction_ledger_repository"
	"wallet-api/internal/repositories/database/tx_repository"
	"wallet-api/internal/repositories/database/wallet_repository"
)

type Contract interface {
	DepositWallet(
		ctx context.Context,
		user models.User,
		request schemas.WalletDepositRequest,
	) error

	WithdrawWallet(
		ctx context.Context,
		user models.User,
		request schemas.WalletWithdrawRequest,
	) error

	WalletTransferBalance(
		ctx context.Context,
		user models.User,
		request schemas.WalletTransferRequest,
	) (err error)
}

type service struct {
	TxRepository                tx_repository.Contract
	WalletRepository            wallet_repository.Contract
	TransactionLedgerRepository transaction_ledger_repository.Contract
}

func New(
	txRepository tx_repository.Contract,
	wallet_repository wallet_repository.Contract,
	transaction_ledger_repository transaction_ledger_repository.Contract,
) Contract {
	return &service{
		TxRepository:                txRepository,
		WalletRepository:            wallet_repository,
		TransactionLedgerRepository: transaction_ledger_repository,
	}
}
