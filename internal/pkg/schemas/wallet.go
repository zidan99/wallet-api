package schemas

import "github.com/greatcloak/decimal"

type WalletHistoryFilter struct {
	WalletID int64 `query:"wallet_id"`
	Page     int   `query:"page"`
	Limit    int   `query:"limit"`
}

type WalletDepositRequest struct {
	WalletID int64           `json:"wallet_id" validate:"required"`
	Nominal  decimal.Decimal `json:"nominal" validate:"required"`
	Note     *string         `json:"note,omitempty"`
	Desc     *string         `json:"desc,omitempty"`
}

type WalletWithdrawRequest struct {
	WalletID int64           `json:"wallet_id" validate:"required"`
	Nominal  decimal.Decimal `json:"nominal" validate:"required"`
	Note     *string         `json:"note,omitempty"`
	Desc     *string         `json:"desc,omitempty"`
}

type WalletTransferRequest struct {
	WalletOriginID      int64 `json:"wallet_origin_id" validate:"required"`
	WalletDestinationID int64 `json:"wallet_destination_id" validate:"required"`

	Nominal decimal.Decimal `json:"nominal" validate:"required"`
	Note    *string         `json:"note,omitempty"`
	Desc    *string         `json:"desc,omitempty"`
}
