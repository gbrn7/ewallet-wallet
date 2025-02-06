package interfaces

import (
	"context"
	"ewallet-wallet/internal/models"
)

type IExternal interface {
	ValidateToken(ctx context.Context, token string) (models.TokenData, error)
}
