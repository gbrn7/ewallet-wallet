package models

import "time"

type Wallet struct {
	ID        int     `json:"id"`
	UserID    uint64  `json:"user_id" gorm:"column:user_id;unique"`
	Balance   float64 `json:"balance" gorm:"column:balance;type:decimal(15,2)"`
	CreatedAt time.Time
	UpdatedAt time.Time
}

func (*Wallet) TableName() string {
	return "wallets"
}

type WalletTransaction struct {
	ID                    int       `json:"-"`
	WalletID              int       `json:"Wallet_ID" gorm:"column:wallet_id"`
	Amount                float64   `json:"amount" gorm:"column:amount;type:decimal(15,2)"`
	WalletTransactionType string    `json:"wallet_transaction_type" gorm:"column:wallet_transaction_type;type:enum('CREDIT', 'DEBIT')"`
	Reference             string    `json:"reference" gorm:"column:reference;type:varchar(100);unique"`
	CreatedAt             time.Time `json:"date"`
	UpdatedAt             time.Time `json:"-"`
}

func (*WalletTransaction) TableName() string {
	return "wallet_transactions"
}

type WalletHistoryParam struct {
	Page                  int    `form:"page"`
	Limit                 int    `form:"limit"`
	WalletTransactionType string `form:"wallet_transaction_type"`
}
