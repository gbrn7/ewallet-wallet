package i_service

import (
	"context"
	"ewallet-wallet/internal/models"
)

type IWalletService interface {
	Create(ctx context.Context, wallet *models.Wallet) error
	CreditBalance(ctx context.Context, userID uint64, req models.TransactionRequest) (models.BalanceResponse, error)
	DebitBalance(ctx context.Context, userID uint64, req models.TransactionRequest) (models.BalanceResponse, error)
	GetBalance(ctx context.Context, userID uint64) (models.BalanceResponse, error)
	GetWalletHistory(ctx context.Context, userID uint64, param models.WalletHistoryParam) ([]models.WalletTransaction, error)
	ExGetBalance(ctx context.Context, walletID int) (models.BalanceResponse, error)

	CreateWalletLink(ctx context.Context, clientSource string, req *models.WalletLink) (*models.WalletStructOTP, error)
	WalletLinkConfirmation(ctx context.Context, walletID int, clientSource string, otp string) error
	WalletUnlink(ctx context.Context, walletID int, clientSource string) error
	ExternalTransaction(ctx context.Context, req models.ExternalTransactionRequest) (models.BalanceResponse, error)
}
