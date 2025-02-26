package wallet

import (
	"context"
	"ewallet-wallet/internal/models"
)

//go:generate mockgen -source=external.go -destination=external_mock_test.go -package=wallet
type External interface {
	ValidateToken(ctx context.Context, token string) (models.TokenData, error)
}
