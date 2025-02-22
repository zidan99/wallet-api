package wallet_service

import (
	"context"
	"errors"
	"wallet-api/internal/models"
	"wallet-api/internal/pkg/helpers"
	"wallet-api/internal/pkg/logger"
	"wallet-api/internal/pkg/schemas"

	"github.com/rotisserie/eris"
)

func (s *service) GetWalletHistory(
	ctx context.Context,
	request schemas.WalletHistoryFilter,
	user models.User,
) (histories []models.TransactionLedger, total int64, err error) {
	metadata := map[string]any{
		"page":  request.Page,
		"limit": request.Limit,
	}

	wallet, err := s.WalletRepository.GetWalletByID(ctx, nil, user.ID, request.WalletID)
	if err != nil {
		logger.Trace(ctx, metadata, err.Error())
		return nil, 0, eris.Wrap(err, "error getting all history wallet")
	}

	if wallet.ID == 0 {
		return nil, 0, eris.Wrap(errors.New("wallet not found"), "wallet not found")
	}

	histories, total, err = s.WalletRepository.GetWalletHistory(ctx, nil, request, wallet)
	if err != nil {
		logger.Trace(ctx, metadata, err.Error())
		return nil, 0, eris.Wrap(err, "error getting all history wallet")
	}

	metadata["data"] = string(helpers.ExtractStructToBytes(histories))

	return histories, total, nil
}

func (s *service) GetWalletByID(
	ctx context.Context,
	user models.User,
	walletID int64,
) (wallet models.Wallet, err error) {
	metadata := map[string]any{
		"wallet_id": walletID,
	}

	wallet, err = s.WalletRepository.GetWalletByID(ctx, nil, user.ID, walletID)
	if err != nil {
		logger.Trace(ctx, metadata, err.Error())
		return wallet, eris.Wrap(err, "error getting wallet")
	}

	if wallet.ID == 0 {
		return wallet, eris.Wrap(errors.New("wallet not found"), "wallet not found")
	}

	metadata["data"] = string(helpers.ExtractStructToBytes(wallet))

	return wallet, nil
}
