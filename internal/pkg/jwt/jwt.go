package jwt

import (
	"time"
	"wallet-api/internal/pkg/env"

	"github.com/golang-jwt/jwt/v5"
)

func GenerateJWT(userID int64, expiredTime time.Duration) (string, error) {
	expirationTime := time.Now().Add(expiredTime)

	jwtSecret := env.GetWithDefault("JWT_SECRET", "defaultValue")
	claims := jwt.MapClaims{
		"exp":       expirationTime.Unix(),
		"entity_id": userID,
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	tokenSigned, err := token.SignedString([]byte(jwtSecret))
	if err != nil {
		return tokenSigned, err
	}

	return tokenSigned, nil
}
