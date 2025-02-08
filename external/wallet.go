package external

import (
	"bytes"
	"context"
	"encoding/json"
	"ewallet-wallet/helpers"
	"fmt"
	"net/http"

	"github.com/pkg/errors"
)

type UpdateBalance struct {
	Reference string  `json:"reference"`
	Amount    float64 `json:"amount"`
}

type UpdateBalanceResponse struct {
	Message string  `json:"reference"`
	Amount  float64 `json:"amount"`
}

type Wallet struct {
}

func (e *Wallet) CreditBalance(ctx context.Context, req UpdateBalance) (*UpdateBalanceResponse, error) {
	payload, err := json.Marshal(req)
	if err != nil {
		return nil, errors.Wrap(err, "failed to marshal json")
	}

	url := helpers.GetEnv("WALLET_HOST", "") + helpers.GetEnv("WALLET_ENDPOINT_CREDIT", "")

	httpReq, err := http.NewRequest(http.MethodPut, url, bytes.NewBuffer(payload))
	if err != nil {
		return nil, errors.Wrap(err, "failed to create wallet http request")
	}

	client := &http.Client{}
	resp, err := client.Do(httpReq)
	if err != nil {
		return nil, errors.Wrap(err, "failed to connect wallet service")
	}

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("got error response from wallet service : %d", resp.StatusCode)
	}

	result := &UpdateBalanceResponse{}
	err = json.NewDecoder(resp.Body).Decode(&result)
	if err != nil {
		return nil, errors.Wrap(err, "failed to read response body")
	}
	defer resp.Body.Close()

	return result, nil
}

func (e *Wallet) DebitBalance(ctx context.Context, req UpdateBalance) (*UpdateBalanceResponse, error) {
	payload, err := json.Marshal(req)
	if err != nil {
		return nil, errors.Wrap(err, "failed to marshal json")
	}

	url := helpers.GetEnv("WALLET_HOST", "") + helpers.GetEnv("WALLET_ENDPOINT_DEBIT", "")

	httpReq, err := http.NewRequest(http.MethodPut, url, bytes.NewBuffer(payload))
	if err != nil {
		return nil, errors.Wrap(err, "failed to create wallet http request")
	}

	client := &http.Client{}
	resp, err := client.Do(httpReq)
	if err != nil {
		return nil, errors.Wrap(err, "failed to connect wallet service")
	}

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("got error response from wallet service : %d", resp.StatusCode)
	}

	result := &UpdateBalanceResponse{}
	err = json.NewDecoder(resp.Body).Decode(&result)
	if err != nil {
		return nil, errors.Wrap(err, "failed to read response body")
	}
	defer resp.Body.Close()

	return result, nil
}
