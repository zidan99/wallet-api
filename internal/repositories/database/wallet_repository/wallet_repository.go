package wallet_repository

import (
	"context"
	"wallet-api/internal/models"
	"wallet-api/internal/pkg/db"
	"wallet-api/internal/pkg/helpers"
	"wallet-api/internal/pkg/logger"
	"wallet-api/internal/pkg/schemas"

	"github.com/rotisserie/eris"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

type Contract interface {
	GetWalletByID(
		ctx context.Context,
		tx *gorm.DB,
		userID int64,
		walletID int64,
	) (wallet models.Wallet, err error)

	GetWalletByIDAnotherUser(
		ctx context.Context,
		tx *gorm.DB,
		walletID int64,
	) (wallet models.Wallet, err error)

	GetWalletHistory(
		ctx context.Context,
		tx *gorm.DB,
		searchParams schemas.WalletHistoryFilter,
		wallet models.Wallet,
	) ([]models.TransactionLedger, int64, error)

	UpdateWalletBalance(
		ctx context.Context,
		tx *gorm.DB,
		wallet *models.Wallet,
		updatedColumns []string,
	) (err error)
}

type repository struct {
	DB *gorm.DB
}

func New(db *gorm.DB) Contract {
	return &repository{DB: db}
}

func (r *repository) GetWalletByID(
	ctx context.Context,
	tx *gorm.DB,
	userID int64,
	walletID int64,
) (wallet models.Wallet, err error) {
	query := db.Use(tx, r.DB).
		Clauses(clause.Locking{Strength: "UPDATE"}).
		WithContext(ctx).
		Where("user_id = ?", userID).
		Where("id = ?", walletID)

	if err := query.Find(&wallet).Error; err != nil {
		logger.Trace(ctx, nil, err.Error())

		return wallet, err
	}

	return wallet, nil
}

func (r *repository) GetWalletByIDAnotherUser(
	ctx context.Context,
	tx *gorm.DB,
	walletID int64,
) (wallet models.Wallet, err error) {
	query := db.Use(tx, r.DB).
		Clauses(clause.Locking{Strength: "UPDATE"}).
		WithContext(ctx).
		Where("id = ?", walletID)

	if err := query.Find(&wallet).Error; err != nil {
		logger.Trace(ctx, nil, err.Error())

		return wallet, err
	}

	return wallet, nil
}

func (r *repository) GetWalletHistory(
	ctx context.Context,
	tx *gorm.DB,
	searchParams schemas.WalletHistoryFilter,
	wallet models.Wallet,
) (history []models.TransactionLedger, total int64, err error) {
	query := r.DB.Model(&models.TransactionLedger{}).
		Where("wallet_id = ?", wallet.ID)

	query.Count(&total)

	query = query.Order(clause.OrderByColumn{
		Column: clause.Column{
			Name: "created_at",
		}, Desc: true,
	})

	query = query.
		Offset((searchParams.Page - 1) * searchParams.Limit).
		Limit(searchParams.Limit)

	if err := query.Find(&history).Error; err != nil {
		errMeta := map[string]any{
			"page":  searchParams.Page,
			"limit": searchParams.Limit,
		}

		logger.Trace(ctx, errMeta, err.Error())

		return nil, 0, err
	}

	return history, total, nil
}

func (r *repository) UpdateWalletBalance(
	ctx context.Context,
	tx *gorm.DB,
	wallet *models.Wallet,
	updatedColumns []string,
) (err error) {
	metadata := map[string]any{
		"request":         string(helpers.ExtractStructToBytes(wallet)),
		"updated_columns": updatedColumns,
	}

	query := db.
		Use(tx, r.DB).
		WithContext(ctx).
		Select(updatedColumns)

	if result := query.Updates(wallet); result.Error != nil {
		logger.Trace(ctx, metadata, result.Error.Error())
		return eris.New(result.Error.Error())
	}

	return nil
}
