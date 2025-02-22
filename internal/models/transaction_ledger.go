package models

import (
	"time"

	"github.com/greatcloak/decimal"
)

type TransactionLedger struct {
	ID          uint            `gorm:"primaryKey;autoIncrement" json:"id"`
	WalletID    uint            `gorm:"not null" json:"wallet_id"`
	Type        string          `gorm:"not null;size:20" json:"type"`
	Note        *string         `gorm:"size:65535" json:"note,omitempty"`
	Status      string          `gorm:"not null;default:'completed';size:20" json:"status"`
	Credit      decimal.Decimal `gorm:"not null;default:0;type:numeric(28,20)" json:"credit"`
	Debit       decimal.Decimal `gorm:"not null;default:0;type:numeric(28,20)" json:"debit"`
	Description *string         `gorm:"size:65535" json:"description,omitempty"`
	CreatedAt   time.Time       `gorm:"not null;default:CURRENT_TIMESTAMP" json:"created_at"`

	Wallet Wallet `gorm:"foreignKey:WalletID" json:"wallet,omitempty"`
}

func (TransactionLedger) TableName() string {
	return "transaction_ledger"
}
