package transaction_ledger_repository

import (
	"context"
	"wallet-api/internal/models"
	"wallet-api/internal/pkg/db"
	"wallet-api/internal/pkg/helpers"
	"wallet-api/internal/pkg/logger"

	"github.com/rotisserie/eris"
	"gorm.io/gorm"
)

type Contract interface {
	InsertTransactionLedger(
		ctx context.Context,
		tx *gorm.DB,
		transaction models.TransactionLedger,
	) (err error)
}

type repository struct {
	DB *gorm.DB
}

func New(db *gorm.DB) Contract {
	return &repository{DB: db}
}

func (r *repository) InsertTransactionLedger(
	ctx context.Context,
	tx *gorm.DB,
	transaction models.TransactionLedger,
) (err error) {
	metadata := map[string]any{
		"request": string(helpers.ExtractStructToBytes(transaction)),
	}

	query := db.
		Use(tx, r.DB).
		WithContext(ctx).
		Model(&models.TransactionLedger{})

	if result := query.Create(&transaction); result.Error != nil {
		logger.Trace(ctx, metadata, result.Error.Error())
		return eris.New(result.Error.Error())
	}

	return nil
}
