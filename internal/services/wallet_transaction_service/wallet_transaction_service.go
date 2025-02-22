package wallet_transaction_service

import (
	"context"
	"errors"
	"time"
	"wallet-api/internal/constant/transaction_type"
	"wallet-api/internal/models"
	"wallet-api/internal/pkg/helpers"
	"wallet-api/internal/pkg/logger"
	"wallet-api/internal/pkg/schemas"

	"github.com/greatcloak/decimal"
	"github.com/rotisserie/eris"
	"gorm.io/gorm"
)

func (s *service) DepositWallet(
	ctx context.Context,
	user models.User,
	request schemas.WalletDepositRequest,
) error {
	metadata := map[string]any{
		"request": request,
	}

	txErr := s.TxRepository.StartTransaction(func(tx *gorm.DB) error {
		wallet, err := s.WalletRepository.GetWalletByID(ctx, tx, user.ID, request.WalletID)
		if err != nil {
			logger.Trace(ctx, metadata, err.Error())
			return eris.Wrap(err, "error getting wallet")
		}

		if wallet.ID == 0 {
			return eris.Wrap(errors.New("wallet not found"), "wallet not found")
		}

		wallet.Balance = wallet.Balance.Add(request.Nominal)
		updateColumn := []string{
			"balance",
		}

		err = s.WalletRepository.UpdateWalletBalance(ctx, tx, &wallet, updateColumn)
		if err != nil {
			logger.Trace(ctx, metadata, err.Error())
			return eris.Wrap(err, "error transaction")
		}

		transaction := models.TransactionLedger{
			WalletID:    wallet.ID,
			Type:        transaction_type.Deposit,
			Note:        request.Note,
			Description: request.Desc,
			Credit:      request.Nominal,
			Debit:       decimal.Zero,
			Status:      "completed",
			CreatedAt:   time.Now(),
		}

		err = s.TransactionLedgerRepository.InsertTransactionLedger(ctx, tx, transaction)
		if err != nil {
			logger.Trace(ctx, metadata, err.Error())
			return eris.Wrap(err, "error transaction")
		}

		return nil
	})

	if txErr != nil {
		logger.Trace(ctx, metadata, txErr.Error())
		return eris.Wrap(txErr, "error when performing database transaction")
	}

	metadata["data"] = string(helpers.ExtractStructToBytes(request))

	return nil
}

func (s *service) WithdrawWallet(
	ctx context.Context,
	user models.User,
	request schemas.WalletWithdrawRequest,
) (err error) {
	metadata := map[string]any{
		"request": request,
	}

	txErr := s.TxRepository.StartTransaction(func(tx *gorm.DB) (err error) {
		wallet, err := s.WalletRepository.GetWalletByID(ctx, tx, user.ID, request.WalletID)
		if err != nil {
			logger.Trace(ctx, metadata, err.Error())
			return eris.Wrap(err, "error getting wallet")
		}

		if wallet.ID == 0 {
			return eris.Wrap(errors.New("wallet not found"), "wallet not found")
		}

		wallet.Balance = wallet.Balance.Sub(request.Nominal)
		if wallet.Balance.LessThan(decimal.Zero) {
			err := errors.New("insufficient balance")

			logger.Trace(ctx, metadata, err.Error())
			return eris.Wrap(err, err.Error())
		}

		updateColumn := []string{
			"balance",
		}

		err = s.WalletRepository.UpdateWalletBalance(ctx, tx, &wallet, updateColumn)
		if err != nil {
			logger.Trace(ctx, metadata, err.Error())
			return eris.Wrap(err, "error transaction")
		}

		transaction := models.TransactionLedger{
			WalletID:    wallet.ID,
			Type:        transaction_type.Withdraw,
			Note:        request.Note,
			Description: request.Desc,
			Debit:       request.Nominal,
			Credit:      decimal.Zero,
			Status:      "completed",
			CreatedAt:   time.Now(),
		}

		err = s.TransactionLedgerRepository.InsertTransactionLedger(ctx, tx, transaction)
		if err != nil {
			logger.Trace(ctx, metadata, err.Error())
			return eris.Wrap(err, "error transaction")
		}

		return nil
	})

	if txErr != nil {
		logger.Trace(ctx, metadata, txErr.Error())
		return eris.Wrap(txErr, "error when performing database transaction")
	}

	metadata["data"] = string(helpers.ExtractStructToBytes(request))

	return err
}

func (s *service) WalletTransferBalance(
	ctx context.Context,
	user models.User,
	request schemas.WalletTransferRequest,
) (err error) {
	metadata := map[string]any{
		"request": request,
	}

	txErr := s.TxRepository.StartTransaction(func(tx *gorm.DB) (err error) {
		walletOrigin, err := s.WalletRepository.GetWalletByID(ctx, tx, user.ID, request.WalletOriginID)
		if err != nil {
			logger.Trace(ctx, metadata, err.Error())
			return eris.Wrap(err, "error getting wallet")
		}

		if walletOrigin.ID == 0 {
			return eris.Wrap(errors.New("wallet origin not found"), "wallet origin not found")
		}

		walletOrigin.Balance = walletOrigin.Balance.Sub(request.Nominal)
		if walletOrigin.Balance.LessThan(decimal.Zero) {
			err := errors.New("insufficient balance")

			logger.Trace(ctx, metadata, err.Error())
			return eris.Wrap(err, err.Error())
		}

		walletDestination, err := s.WalletRepository.GetWalletByIDAnotherUser(ctx, tx, request.WalletDestinationID)
		if err != nil {
			logger.Trace(ctx, metadata, err.Error())
			return eris.Wrap(err, "error getting wallet")
		}

		if walletDestination.ID == 0 {
			return eris.Wrap(errors.New("wallet destination not found"), "wallet destination not found")
		}

		updateColumn := []string{
			"balance",
		}

		err = s.WalletRepository.UpdateWalletBalance(ctx, tx, &walletOrigin, updateColumn)
		if err != nil {
			logger.Trace(ctx, metadata, err.Error())
			return eris.Wrap(err, "error transaction")
		}

		transaction := models.TransactionLedger{
			WalletID:    walletOrigin.ID,
			Type:        transaction_type.TransferOut,
			Note:        request.Note,
			Description: request.Desc,
			Debit:       request.Nominal,
			Credit:      decimal.Zero,
			Status:      "completed",
			CreatedAt:   time.Now(),
		}

		err = s.TransactionLedgerRepository.InsertTransactionLedger(ctx, tx, transaction)
		if err != nil {
			logger.Trace(ctx, metadata, err.Error())
			return eris.Wrap(err, "error transaction")
		}

		err = s.WalletRepository.UpdateWalletBalance(ctx, tx, &walletDestination, updateColumn)
		if err != nil {
			logger.Trace(ctx, metadata, err.Error())
			return eris.Wrap(err, "error transaction")
		}

		transaction = models.TransactionLedger{
			WalletID:    walletDestination.ID,
			Type:        transaction_type.TransferIn,
			Note:        request.Note,
			Description: request.Desc,
			Credit:      request.Nominal,
			Debit:       decimal.Zero,
			Status:      "completed",
			CreatedAt:   time.Now(),
		}

		err = s.TransactionLedgerRepository.InsertTransactionLedger(ctx, tx, transaction)
		if err != nil {
			logger.Trace(ctx, metadata, err.Error())
			return eris.Wrap(err, "error transaction")
		}

		return nil
	})

	if txErr != nil {
		logger.Trace(ctx, metadata, txErr.Error())
		return eris.Wrap(txErr, "error when performing database transaction")
	}

	metadata["data"] = string(helpers.ExtractStructToBytes(request))

	return err
}
