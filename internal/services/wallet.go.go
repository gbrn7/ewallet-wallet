package services

import (
	"context"
	"ewallet-wallet/internal/interfaces/i_repository"
	"ewallet-wallet/internal/models"
	"fmt"
	"math/rand"
	"strconv"

	"github.com/pkg/errors"
	"gorm.io/gorm"
)

type WalletService struct {
	WalletRepo i_repository.IWalletRepo
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

func (s *WalletService) ExGetBalance(ctx context.Context, walletID int) (models.BalanceResponse, error) {
	var (
		resp models.BalanceResponse
	)

	wallet, err := s.WalletRepo.GetWalletByID(ctx, walletID)
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

func (s *WalletService) CreateWalletLink(ctx context.Context, clientSource string, req *models.WalletLink) (*models.WalletStructOTP, error) {
	req.ClientSource = clientSource
	req.Status = "pending"
	req.OTP = strconv.Itoa(rand.Intn(999999))

	err := s.WalletRepo.InsertWalletLink(ctx, req)
	if err != nil {
		return nil, errors.Wrap(err, "failed to insert wallet link")
	}

	resp := models.WalletStructOTP{
		OTP: req.OTP,
	}

	return &resp, nil
}

func (s *WalletService) WalletLinkConfirmation(ctx context.Context, walletID int, clientSource string, otp string) error {
	walletLink, err := s.WalletRepo.GetWalletLink(ctx, walletID, clientSource)
	if err != nil {
		return errors.Wrap(err, "failed to get wallet link")
	}

	if walletLink.Status != "pending" {
		return errors.New("wallet status is not pending")
	}

	if walletLink.OTP != otp {
		return fmt.Errorf("invalid otp. requested = %s, store = %s", otp, walletLink.OTP)
	}

	return s.WalletRepo.UpdateStatusWalletLink(ctx, walletID, clientSource, "linked")
}

func (s *WalletService) WalletUnlink(ctx context.Context, walletID int, clientSource string) error {

	return s.WalletRepo.UpdateStatusWalletLink(ctx, walletID, clientSource, "unlinked")
}

func (s *WalletService) ExternalTransaction(ctx context.Context, req models.ExternalTransactionRequest) (models.BalanceResponse, error) {
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

	amount := req.Amount
	if req.TransactionType == "DEBIT" {
		amount = -req.Amount
	}

	wallet, err := s.WalletRepo.UpdateBalanceByID(ctx, req.WalletID, amount)
	if err != nil {
		return resp, errors.Wrap(err, "failed to updated balance")
	}

	walletTrx := &models.WalletTransaction{
		WalletID:              wallet.ID,
		Amount:                req.Amount,
		Reference:             req.Reference,
		WalletTransactionType: req.TransactionType,
	}

	err = s.WalletRepo.CreateWalletTrx(ctx, walletTrx)
	if err != nil {
		return resp, errors.Wrap(err, "failed to insert wallet transaction")
	}

	resp.Balance = wallet.Balance + amount

	return resp, nil
}
