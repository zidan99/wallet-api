package user_repository

import (
	"context"
	"wallet-api/internal/models"
	"wallet-api/internal/pkg/db"
	"wallet-api/internal/pkg/logger"

	"github.com/rotisserie/eris"
	"gorm.io/gorm"
)

type Contract interface {
	Login(ctx context.Context, user models.User) (models.User, error)
	FindByID(
		ctx context.Context,
		tx *gorm.DB,
		id int64,
	) (user *models.User, err error)
}

type repository struct {
	DB *gorm.DB
}

func New(db *gorm.DB) Contract {
	return &repository{DB: db}
}

func (r *repository) Login(ctx context.Context, user models.User) (models.User, error) {
	var result models.User
	err := r.DB.
		WithContext(ctx).
		Where("email = ?", user.Email).
		First(&result).
		Error

	if err != nil {
		return models.User{}, err
	}

	return result, nil
}

func (r *repository) FindByID(
	ctx context.Context,
	tx *gorm.DB,
	id int64,
) (user *models.User, err error) {
	metadata := map[string]any{
		"user_id": id,
	}

	query := db.
		Use(tx, r.DB).
		WithContext(ctx).
		Limit(1).Preload("Wallets")

	if result := query.Find(&user); result.Error != nil {
		logger.Trace(ctx, metadata, result.Error.Error())
		return nil, eris.New(result.Error.Error())
	}

	return user, nil
}
