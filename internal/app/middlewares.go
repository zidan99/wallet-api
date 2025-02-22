package app

import "wallet-api/internal/middlewares/auth_middleware"

type Middlewares struct {
	AuthMiddleware auth_middleware.Contract
}

func RegisterMiddlewares(s Services) Middlewares {
	return Middlewares{
		AuthMiddleware: auth_middleware.New(s.UserService),
	}
}
