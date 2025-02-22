package auth_middleware

import (
	"context"
	"strings"
	"time"
	"wallet-api/internal/constant/keys"
	"wallet-api/internal/pkg/env"
	"wallet-api/internal/pkg/schemas"

	jwtware "github.com/gofiber/contrib/jwt"
	"github.com/gofiber/fiber/v2"
	"github.com/golang-jwt/jwt/v5"
)

func (m *middleware) ValidateUserToken(c *fiber.Ctx) error {
	tokenHeader := c.Get("Authorization")
	const bearerPrefix = "Bearer "
	if tokenHeader == "" {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
			"error": "Unauthorized",
		})
	}

	var token string
	if strings.HasPrefix(tokenHeader, bearerPrefix) {
		token = strings.TrimPrefix(tokenHeader, bearerPrefix)
	} else {
		token = tokenHeader
	}

	jwtSecret := env.GetWithDefault("JWT_SECRET", "defaultValue")
	parsedToken, err := jwt.Parse(token, func(token *jwt.Token) (interface{}, error) {
		return []byte(jwtSecret), nil
	})

	if err != nil {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
			"error": "Unauthorized",
		})
	}

	claims, ok := parsedToken.Claims.(jwt.MapClaims)
	if !ok {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
			"error": "Invalid token claims",
		})
	}

	if exp, ok := claims["exp"].(float64); ok {
		if time.Now().Unix() > int64(exp) {
			return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
				"error": "Token expired",
			})
		}
	}

	userID, ok := claims["entity_id"].(float64)
	if !ok {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
			"error": "Invalid user ID in token",
		})
	}

	user, err := m.UserService.GetUserByID(context.Background(), int64(userID))
	if err != nil {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
			"error": "User not found",
		})
	}

	// Add user to context
	c.Locals(keys.USER_AUTH_KEY, *user)

	return c.Next()
}

func (m *middleware) JWTAuth() fiber.Handler {
	JwtSecret := env.GetWithDefault("JWT_SECRET", "defaultValue")
	return jwtware.New(jwtware.Config{
		SigningKey:   jwtware.SigningKey{Key: []byte(JwtSecret)},
		ErrorHandler: jwtError,
	})
}

func (m *middleware) CheckTokenValidation() fiber.Handler {
	return func(ctx *fiber.Ctx) error {
		tokenHeader := ctx.Get("Authorization")
		const bearerPrefix = "Bearer "
		var token string

		if strings.HasPrefix(tokenHeader, bearerPrefix) {
			token = strings.TrimPrefix(tokenHeader, bearerPrefix)
		} else {
			token = tokenHeader
		}

		jwtSecret := env.GetWithDefault("JWT_SECRET", "defaultValue")
		parsedToken, err := jwt.Parse(token, func(token *jwt.Token) (interface{}, error) {
			return []byte(jwtSecret), nil
		})

		if err != nil {
			return ctx.Status(fiber.StatusUnauthorized).JSON(schemas.GeneralResponse{
				Code:    fiber.StatusUnauthorized,
				Message: "Invalid Token",
				Data:    nil,
			})
		}

		claims, ok := parsedToken.Claims.(jwt.MapClaims)
		if !ok {
			return ctx.Status(fiber.StatusUnauthorized).JSON(schemas.GeneralResponse{
				Code:    fiber.StatusUnauthorized,
				Message: "Invalid Token Claims",
				Data:    nil,
			})
		}

		expClaim := claims["exp"].(float64)
		id := claims["entity_id"].(int64)

		expirationTime := time.Unix(int64(expClaim), 0)
		return ctx.Status(fiber.StatusOK).JSON(schemas.GeneralResponse{
			Code:    fiber.StatusOK,
			Message: "Token Is Valid",
			Data: fiber.Map{
				"EntityID":         id,
				"Token Expired At": expirationTime,
			},
		})

	}
}

func jwtError(c *fiber.Ctx, err error) error {
	if err.Error() == "Missing or malformed JWT" {
		return c.Status(fiber.StatusBadRequest).
			JSON(schemas.GeneralResponse{
				Code:    fiber.StatusUnauthorized,
				Message: "Missing or Malformed JWT",
				Data:    err.Error(),
			})
	}
	return c.Status(fiber.StatusUnauthorized).
		JSON(schemas.GeneralResponse{
			Code:    fiber.StatusUnauthorized,
			Message: "Invalid or Expired JWT",
			Data:    err.Error(),
		})
}
