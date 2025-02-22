package user_service

import (
	"context"
	"errors"
	"time"
	"wallet-api/internal/models"
	"wallet-api/internal/pkg/helpers"
	jwtPkg "wallet-api/internal/pkg/jwt"
	"wallet-api/internal/pkg/logger"
	"wallet-api/internal/pkg/schemas"
)

func (s *service) Login(ctx context.Context, request schemas.LoginRequest) (user models.User, token string, err error) {
	metadata := map[string]any{}

	input := models.User{
		Email:    request.Email,
		Password: request.Password,
	}

	user, err = s.UserRepository.Login(ctx, input)
	if err != nil {
		metadata["error_message"] = err.Error()

		logger.Trace(ctx, metadata, "invalid credentials")
		return input, "", errors.New("invalid credentials")
	}

	checkPasswordErr := helpers.CheckHashedPasswordMatches(user.Password, input.Password)
	if checkPasswordErr != nil {
		return input, "", errors.New("invalid credentials")
	}

	token, err = jwtPkg.GenerateJWT(user.ID, 60*time.Minute)
	if err != nil {
		return user, "", err
	}

	return user, token, nil
}

func (s *service) GetUserByID(
	ctx context.Context,
	userID int64,
) (user *models.User, err error) {
	metadata := map[string]any{"user_id": userID}

	user, err = s.UserRepository.FindByID(ctx, nil, userID)
	if err != nil {
		logger.Trace(ctx, metadata, err.Error())
		return nil, err
	}

	return user, nil
}
