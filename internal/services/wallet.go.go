package services

import (
	"context"
	"ewallet-wallet/internal/interfaces"
	"ewallet-wallet/internal/models"

	"github.com/pkg/errors"
	"gorm.io/gorm"
)

type WalletService struct {
	WalletRepo interfaces.IWalletRepo
}

func (s *WalletService) Create(ctx context.Context, wallet *models.Wallet) error {
	return s.WalletRepo.CreateWallet(ctx, wallet)
}

func (s *WalletService) CreditBalance(ctx context.Context, userID uint64, req models.TransactionRequest) (models.BalanceResponse, error) {
	var (
		resp models.BalanceResponse
	)

	history, err := s.WalletRepo.GetWalletTransactionByReference(ctx, req.Reference)
	if err != nil {
		if err != gorm.ErrRecordNotFound {
			return resp, errors.Wrap(err, "failed to check reference")
		}
	}

	if history.ID > 0 {
		return resp, errors.New("reference is duplicated")
	}

	wallet, err := s.WalletRepo.UpdateBalance(ctx, userID, req.Amount)
	if err != nil {
		return resp, errors.Wrap(err, "failed to updated balance")
	}

	walletTrx := &models.WalletTransaction{
		WalletID:              wallet.ID,
		Amount:                req.Amount,
		Reference:             req.Reference,
		WalletTransactionType: "CREDIT",
	}

	err = s.WalletRepo.CreateWalletTrx(ctx, walletTrx)
	if err != nil {
		return resp, errors.Wrap(err, "failed to insert wallet transaction")
	}

	resp.Balance = wallet.Balance + req.Amount

	return resp, nil
}

func (s *WalletService) DebitBalance(ctx context.Context, userID uint64, req models.TransactionRequest) (models.BalanceResponse, error) {
	var (
		resp models.BalanceResponse
	)

	history, err := s.WalletRepo.GetWalletTransactionByReference(ctx, req.Reference)
	if err != nil {
		if err != gorm.ErrRecordNotFound {
			return resp, errors.Wrap(err, "failed to check reference")
		}
	}

	if history.ID > 0 {
		return resp, errors.New("reference is duplicated")
	}

	wallet, err := s.WalletRepo.UpdateBalance(ctx, userID, -req.Amount)
	if err != nil {
		return resp, errors.Wrap(err, "failed to updated balance")
	}

	walletTrx := &models.WalletTransaction{
		WalletID:              wallet.ID,
		Amount:                req.Amount,
		Reference:             req.Reference,
		WalletTransactionType: "DEBIT",
	}

	err = s.WalletRepo.CreateWalletTrx(ctx, walletTrx)
	if err != nil {
		return resp, errors.Wrap(err, "failed to insert wallet transaction")
	}

	resp.Balance = wallet.Balance - req.Amount

	return resp, nil
}

func (s *WalletService) GetBalance(ctx context.Context, userID uint64) (models.BalanceResponse, error) {
	var (
		resp models.BalanceResponse
	)

	wallet, err := s.WalletRepo.GetWalletByUserID(ctx, userID)
	if err != nil {
		return resp, errors.Wrap(err, "failed to get wallet")
	}

	resp.Balance = wallet.Balance

	return resp, nil
}

func (s *WalletService) GetWalletHistory(ctx context.Context, userID uint64, param models.WalletHistoryParam) ([]models.WalletTransaction, error) {
	var (
		resp []models.WalletTransaction
	)

	wallet, err := s.WalletRepo.GetWalletByUserID(ctx, userID)
	if err != nil {
		return resp, errors.Wrap(err, "failed to get wallet")
	}

	offset := (param.Page - 1) * param.Limit

	resp, err = s.WalletRepo.GetWalletHistory(ctx, wallet.ID, offset, param.Limit, param.WalletTransactionType)

	if err != nil {
		return resp, errors.Wrap(err, "failed to get wallet history")
	}

	return resp, nil
}
