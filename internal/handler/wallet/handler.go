package wallet

import (
	"context"
	"ewallet-wallet/internal/models"

	"github.com/gin-gonic/gin"
)

//go:generate mockgen -source=handler.go -destination=handler_mock_test.go -package=wallet
type Service interface {
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

type Handler struct {
	*gin.Engine
	Service    Service
	External   External
	Middleware Middleware
}

func NewHandler(api *gin.Engine, service Service, ext External, mdw Middleware) *Handler {
	return &Handler{
		api,
		service,
		ext,
		mdw,
	}
}

func (h *Handler) RegisterRoute() {
	walletV1 := h.Group("/wallet/v1")
	walletV1.POST("/", h.Create)
	walletV1.PUT("/balance/credit", h.Middleware.MiddlewareValidateToken, h.CreditBalance)
	walletV1.PUT("/balance/debit", h.Middleware.MiddlewareValidateToken, h.DebitBalance)
	walletV1.GET("/balance", h.Middleware.MiddlewareValidateToken, h.GetBalance)
	walletV1.GET("/history", h.Middleware.MiddlewareValidateToken, h.GetWalletHistory)

	exWalletv1 := walletV1.Group("/ex")
	exWalletv1.Use(h.Middleware.MiddlewareSignatureValidation)
	exWalletv1.POST("/link", h.CreateWalletLink)
	exWalletv1.PUT("/link/:wallet_id/confirmation", h.WalletLinkConfirmation)
	exWalletv1.DELETE("/:wallet_id/unlink", h.WalletUnlink)
	exWalletv1.GET("/:wallet_id/balance", h.ExGetBalance)
	exWalletv1.POST("/transaction", h.ExternalTransaction)
}
