package repository

import (
	"context"
	"ewallet-wallet/internal/models"
	"fmt"

	"gorm.io/gorm"
)

type WalletRepo struct {
	DB *gorm.DB
}

func (r *WalletRepo) CreateWallet(ctx context.Context, wallet *models.Wallet) error {
	return r.DB.Create(wallet).Error
}

func (r *WalletRepo) UpdateBalance(ctx context.Context, userID uint64, amount float64) (models.Wallet, error) {
	var wallet models.Wallet
	err := r.DB.Transaction(func(tx *gorm.DB) error {
		err := tx.Raw("SELECT id, user_id, balance FROM wallets WHERE user_id = ? FOR UPDATE", userID).Scan(&wallet).Error
		if err != nil {
			return err
		}

		if (wallet.Balance + amount) < 0 {
			return fmt.Errorf("current balance is not enough to perform the transaction: %f - %f", wallet.Balance, amount)
		}

		err = tx.Exec("UPDATE wallets SET balance = balance + ? WHERE user_id = ?", amount, userID).Error
		if err != nil {
			return err
		}
		return nil
	})

	return wallet, err
}

func (r *WalletRepo) CreateWalletTrx(ctx context.Context, walletHistory *models.WalletTransaction) error {
	return r.DB.Create(walletHistory).Error
}

func (r *WalletRepo) GetWalletTransactionByReference(ctx context.Context, reference string) (models.WalletTransaction, error) {
	var (
		resp models.WalletTransaction
	)

	err := r.DB.Where("reference = ?", reference).Last(&resp).Error

	return resp, err
}
