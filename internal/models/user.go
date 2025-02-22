package models

import (
	"time"
)

type User struct {
	ID       int64  `json:"id"`
	Email    string `json:"email"`
	Password string `json:"-"`

	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`

	Wallets []Wallet `gorm:"foreignKey:UserID" json:"wallets,omitempty"`
}
