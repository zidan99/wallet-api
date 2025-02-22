package wallet_service

import (
	"context"
	"wallet-api/internal/models"
	"wallet-api/internal/pkg/schemas"
	"wallet-api/internal/repositories/database/tx_repository"
	"wallet-api/internal/repositories/database/wallet_repository"
)

type Contract interface {
	GetWalletHistory(
		ctx context.Context,
		request schemas.WalletHistoryFilter,
		user models.User,
	) (histories []models.TransactionLedger, total int64, err error)

	GetWalletByID(
		ctx context.Context,
		user models.User,
		walletID int64,
	) (wallet models.Wallet, err error)
}

type service struct {
	txRepository     tx_repository.Contract
	WalletRepository wallet_repository.Contract
}

func New(txRepository tx_repository.Contract, wallet_repository wallet_repository.Contract) Contract {
	return &service{
		txRepository:     txRepository,
		WalletRepository: wallet_repository,
	}
}
