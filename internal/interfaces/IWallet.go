package interfaces

import (
	"context"
	"ewallet-wallet/internal/models"

	"github.com/gin-gonic/gin"
)

type IWalletRepo interface {
	CreateWallet(ctx context.Context, wallet *models.Wallet) error
	UpdateBalance(ctx context.Context, userID uint64, amount float64) (models.Wallet, error)
	CreateWalletTrx(ctx context.Context, walletHistory *models.WalletTransaction) error
	GetWalletTransactionByReference(ctx context.Context, reference string) (models.WalletTransaction, error)
	GetWalletByUserID(ctx context.Context, userID uint64) (models.Wallet, error)
}

type IWalletService interface {
	Create(ctx context.Context, wallet *models.Wallet) error
	CreditBalance(ctx context.Context, userID uint64, req models.TransactionRequest) (models.BalanceResponse, error)
	DebitBalance(ctx context.Context, userID uint64, req models.TransactionRequest) (models.BalanceResponse, error)
	GetBalance(ctx context.Context, userID uint64) (models.BalanceResponse, error)
}

type IWalletAPI interface {
	Create(*gin.Context)
	CreditBalance(c *gin.Context)
	DebitBalance(c *gin.Context)
	GetBalance(c *gin.Context)
}
