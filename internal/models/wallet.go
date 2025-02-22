package models

import (
	"time"

	"github.com/greatcloak/decimal"
)

type Wallet struct {
	ID        uint            `gorm:"primaryKey;autoIncrement" json:"id"`
	UserID    uint            `gorm:"not null" json:"user_id"`
	Balance   decimal.Decimal `gorm:"not null;default:0;type:numeric(28,20)" json:"balance"`
	CreatedAt time.Time       `gorm:"not null;default:CURRENT_TIMESTAMP" json:"created_at"`
	UpdatedAt time.Time       `gorm:"not null;default:CURRENT_TIMESTAMP" json:"updated_at"`

	User User `gorm:"foreignKey:UserID" json:"user,omitempty"`
}
