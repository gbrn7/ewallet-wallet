package models

import "github.com/go-playground/validator"

type TransactionRequest struct {
	Reference string  `json:"reference" valid:"required"`
	Amount    float64 `json:"amount" valid:"required"`
}

func (l TransactionRequest) Validate() error {
	v := validator.New()
	return v.Struct(l)
}

type BalanceResponse struct {
	Balance float64 `json:"balance"`
}

type ExternalTransactionRequest struct {
	Amount          float64 `json:"amount" valid:"required"`
	Reference       string  `json:"reference" valid:"required"`
	TransactionType string  `json:"transaction_type" valid:"required"`
	WalletID        int     `json:"wallet_id" valid:"required"`
}

func (l ExternalTransactionRequest) Validate() error {
	v := validator.New()
	return v.Struct(l)
}
