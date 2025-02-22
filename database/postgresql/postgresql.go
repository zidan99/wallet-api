package postgresql

import (
	"errors"
	"wallet-api/internal/pkg/env"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

func Open() (*gorm.DB, error) {
	conn := env.Get("DB_CONNECTION")
	if conn == "" {
		return nil, errors.New("connection is not provided")
	}

	db, err := gorm.Open(postgres.Open(conn), &gorm.Config{
		SkipDefaultTransaction: true,
	})

	if err != nil {
		return nil, err
	}

	sql, err := db.DB()
	if err != nil {
		return nil, err
	}

	sql.SetMaxOpenConns(500)

	return db, nil
}
