package i_repository

import (
	"context"
	"ewallet-wallet/internal/models"
)

//go:generate mockgen -source=i_wallet_repository.go -destination=../../services/wallet_mock_test.go -package=services
type IWalletRepo interface {
	CreateWallet(ctx context.Context, wallet *models.Wallet) error
	UpdateBalance(ctx context.Context, userID uint64, amount float64) (models.Wallet, error)
	CreateWalletTrx(ctx context.Context, walletHistory *models.WalletTransaction) error
	GetWalletTransactionByReference(ctx context.Context, reference string) (models.WalletTransaction, error)
	GetWalletByUserID(ctx context.Context, userID uint64) (models.Wallet, error)
	GetWalletHistory(ctx context.Context, walletID int, offset int, limit int, transactionType string) ([]models.WalletTransaction, error)
	GetWalletByID(ctx context.Context, walletID int) (models.Wallet, error)

	InsertWalletLink(ctx context.Context, req *models.WalletLink) error
	GetWalletLink(ctx context.Context, walletID int, clientSource string) (models.WalletLink, error)
	UpdateStatusWalletLink(ctx context.Context, walletID int, clientSource string, status string) error
	UpdateBalanceByID(ctx context.Context, walletID int, amount float64) (models.Wallet, error)
}
