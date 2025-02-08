package external

import (
	"context"
	"ewallet-wallet/constants"
	"ewallet-wallet/external/proto/tokenvalidation"
	"ewallet-wallet/internal/models"
	"fmt"

	"github.com/pkg/errors"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

type External struct {
}

func (*External) ValidateToken(ctx context.Context, token string) (models.TokenData, error) {
	var (
		resp models.TokenData
	)

	conn, err := grpc.NewClient("localhost:7000", grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		return resp, errors.Wrap(err, "failed to dial grpc")
	}

	defer conn.Close()

	client := tokenvalidation.NewTokenValidationClient(conn)
	req := &tokenvalidation.TokenRequest{
		Token: token,
	}
	response, err := client.ValidateToken(ctx, req)
	if err != nil {
		return resp, errors.Wrap(err, "failed to validate token")
	}

	if response.Message != constants.SuccessMessage {
		return resp, fmt.Errorf("got response error from ums: %s", response.Message)
	}

	resp.UserID = response.Data.UserId
	resp.Username = response.Data.Username
	resp.Fullname = response.Data.FullName
	resp.Email = response.Data.Email

	return resp, nil
}
