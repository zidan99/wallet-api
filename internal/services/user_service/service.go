package user_service

import (
	"context"
	"wallet-api/internal/models"
	"wallet-api/internal/pkg/schemas"
	"wallet-api/internal/repositories/database/tx_repository"
	"wallet-api/internal/repositories/database/user_repository"
)

type Contract interface {
	Login(ctx context.Context, request schemas.LoginRequest) (user models.User, token string, err error)
	GetUserByID(
		ctx context.Context,
		userID int64,
	) (user *models.User, err error)
}

type service struct {
	txRepository   tx_repository.Contract
	UserRepository user_repository.Contract
}

func New(txRepository tx_repository.Contract, user_repository user_repository.Contract) Contract {
	return &service{
		txRepository:   txRepository,
		UserRepository: user_repository,
	}
}
