package repository

import (
	"context"
	"ewallet-wallet/internal/models"
	"reflect"
	"regexp"
	"testing"
	"time"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/stretchr/testify/assert"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

func TestWalletRepo_CreateWallet(t *testing.T) {
	db, mock, err := sqlmock.New()
	assert.NoError(t, err)
	defer db.Close()

	gormDB, err := gorm.Open(mysql.New(mysql.Config{
		Conn:                      db,
		SkipInitializeWithVersion: true,
	}), &gorm.Config{})

	assert.NoError(t, err)

	type args struct {
		ctx    context.Context
		wallet *models.Wallet
	}
	tests := []struct {
		name    string
		args    args
		wantErr bool
		mockFn  func(args args)
	}{
		{
			name:    "success",
			wantErr: false,
			args: args{
				ctx: context.Background(),
				wallet: &models.Wallet{
					UserID:  1,
					Balance: 20000000,
				},
			},
			mockFn: func(args args) {
				mock.ExpectBegin()
				mock.ExpectExec(regexp.QuoteMeta("INSERT INTO `wallets` (`user_id`,`balance`,`created_at`,`updated_at`) VALUES (?,?,?,?)")).WithArgs(
					args.wallet.UserID,
					args.wallet.Balance,
					sqlmock.AnyArg(),
					sqlmock.AnyArg(),
				).WillReturnResult(sqlmock.NewResult(1, 1))
				mock.ExpectCommit()
			},
		},
		{
			name:    "error",
			wantErr: true,
			args: args{
				ctx: context.Background(),
				wallet: &models.Wallet{
					UserID:  1,
					Balance: 20000000,
				},
			},
			mockFn: func(args args) {
				mock.ExpectBegin()
				mock.ExpectExec(regexp.QuoteMeta("INSERT INTO `wallets` (`user_id`,`balance`,`created_at`,`updated_at`) VALUES (?,?,?,?)")).WithArgs(
					args.wallet.UserID,
					args.wallet.Balance,
					sqlmock.AnyArg(),
					sqlmock.AnyArg(),
				).WillReturnError(assert.AnError)
				mock.ExpectRollback()
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.mockFn(tt.args)
			r := &WalletRepo{
				DB: gormDB,
			}
			if err := r.CreateWallet(tt.args.ctx, tt.args.wallet); (err != nil) != tt.wantErr {
				t.Errorf("WalletRepo.CreateWallet() error = %v, wantErr %v", err, tt.wantErr)
			}
			assert.NoError(t, mock.ExpectationsWereMet())
		})
	}
}

func TestWalletRepo_UpdateBalance(t *testing.T) {
	db, mock, err := sqlmock.New()
	assert.NoError(t, err)
	defer db.Close()

	gormDB, err := gorm.Open(mysql.New(mysql.Config{
		Conn:                      db,
		SkipInitializeWithVersion: true,
	}), &gorm.Config{})

	assert.NoError(t, err)

	type args struct {
		ctx    context.Context
		userID uint64
		amount float64
	}
	tests := []struct {
		name    string
		args    args
		want    models.Wallet
		wantErr bool
		mockFn  func(args args)
	}{
		{
			name: "success add credit",
			args: args{
				ctx:    context.Background(),
				userID: 1,
				amount: 20000,
			},
			want: models.Wallet{
				ID:      1,
				UserID:  1,
				Balance: 100000,
			},
			wantErr: false,
			mockFn: func(args args) {
				mock.ExpectBegin()

				mock.ExpectQuery(regexp.QuoteMeta("SELECT id, user_id, balance FROM wallets WHERE user_id = ? FOR UPDATE")).WithArgs(
					args.userID,
				).WillReturnRows(sqlmock.NewRows([]string{"id", "user_id", "balance"}).AddRow(1, 1, 100000))

				mock.ExpectExec(regexp.QuoteMeta("UPDATE wallets SET balance = balance + ? WHERE user_id = ?")).WithArgs(
					args.amount,
					args.userID,
				).WillReturnResult(sqlmock.NewResult(1, 1))

				mock.ExpectCommit()
			},
		},
		{
			name: "error get user",
			args: args{
				ctx:    context.Background(),
				userID: 1,
				amount: 20000,
			},
			want:    models.Wallet{},
			wantErr: true,
			mockFn: func(args args) {
				mock.ExpectBegin()

				mock.ExpectQuery(regexp.QuoteMeta("SELECT id, user_id, balance FROM wallets WHERE user_id = ? FOR UPDATE")).WithArgs(
					args.userID,
				).WillReturnError(assert.AnError)

				mock.ExpectRollback()
			},
		},
		{
			name: "error add credit",
			args: args{
				ctx:    context.Background(),
				userID: 1,
				amount: 20000,
			},
			want: models.Wallet{
				ID:      1,
				UserID:  1,
				Balance: 100000,
			},
			wantErr: true,
			mockFn: func(args args) {
				mock.ExpectBegin()

				mock.ExpectQuery(regexp.QuoteMeta("SELECT id, user_id, balance FROM wallets WHERE user_id = ? FOR UPDATE")).WithArgs(
					args.userID,
				).WillReturnRows(sqlmock.NewRows([]string{"id", "user_id", "balance"}).AddRow(1, 1, 100000))

				mock.ExpectExec(regexp.QuoteMeta("UPDATE wallets SET balance = balance + ? WHERE user_id = ?")).WithArgs(
					args.amount,
					args.userID,
				).WillReturnError(assert.AnError)

				mock.ExpectRollback()
			},
		},
		{
			name: "success decrease balance",
			args: args{
				ctx:    context.Background(),
				userID: 1,
				amount: -20000,
			},
			want: models.Wallet{
				ID:      1,
				UserID:  1,
				Balance: 100000,
			},
			wantErr: false,
			mockFn: func(args args) {
				mock.ExpectBegin()

				mock.ExpectQuery(regexp.QuoteMeta("SELECT id, user_id, balance FROM wallets WHERE user_id = ? FOR UPDATE")).WithArgs(
					args.userID,
				).WillReturnRows(sqlmock.NewRows([]string{"id", "user_id", "balance"}).AddRow(1, 1, 100000))

				mock.ExpectExec(regexp.QuoteMeta("UPDATE wallets SET balance = balance + ? WHERE user_id = ?")).WithArgs(
					args.amount,
					args.userID,
				).WillReturnResult(sqlmock.NewResult(1, 1))

				mock.ExpectCommit()
			},
		},
		{
			name: "error decrease balance",
			args: args{
				ctx:    context.Background(),
				userID: 1,
				amount: -120000,
			},
			want: models.Wallet{
				ID:      1,
				UserID:  1,
				Balance: 100000,
			},
			wantErr: true,
			mockFn: func(args args) {
				mock.ExpectBegin()

				mock.ExpectQuery(regexp.QuoteMeta("SELECT id, user_id, balance FROM wallets WHERE user_id = ? FOR UPDATE")).WithArgs(
					args.userID,
				).WillReturnRows(sqlmock.NewRows([]string{"id", "user_id", "balance"}).AddRow(1, 1, 100000))

				mock.ExpectRollback()
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.mockFn(tt.args)
			r := &WalletRepo{
				DB: gormDB,
			}
			got, err := r.UpdateBalance(tt.args.ctx, tt.args.userID, tt.args.amount)
			if (err != nil) != tt.wantErr {
				t.Errorf("WalletRepo.UpdateBalance() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("WalletRepo.UpdateBalance() = %v, want %v", got, tt.want)
			}
			assert.NoError(t, mock.ExpectationsWereMet())
		})
	}
}

func TestWalletRepo_CreateWalletTrx(t *testing.T) {
	db, mock, err := sqlmock.New()
	assert.NoError(t, err)
	defer db.Close()

	gormDB, err := gorm.Open(mysql.New(mysql.Config{
		Conn:                      db,
		SkipInitializeWithVersion: true,
	}), &gorm.Config{})

	assert.NoError(t, err)

	type args struct {
		ctx           context.Context
		walletHistory *models.WalletTransaction
	}
	tests := []struct {
		name    string
		args    args
		wantErr bool
		mockFn  func(args args)
	}{
		{
			name: "success",
			args: args{
				ctx: context.Background(),
				walletHistory: &models.WalletTransaction{
					WalletID:              1,
					Amount:                200000,
					WalletTransactionType: "DEBIT",
					Reference:             "reference",
				},
			},
			wantErr: false,
			mockFn: func(args args) {
				mock.ExpectBegin()
				mock.ExpectExec(regexp.QuoteMeta("INSERT INTO `wallet_transactions` (`wallet_id`,`amount`,`wallet_transaction_type`,`reference`,`created_at`,`updated_at`) VALUES (?,?,?,?,?,?)")).WithArgs(
					args.walletHistory.WalletID,
					args.walletHistory.Amount,
					args.walletHistory.WalletTransactionType,
					args.walletHistory.Reference,
					sqlmock.AnyArg(),
					sqlmock.AnyArg(),
				).WillReturnResult(sqlmock.NewResult(1, 1))
				mock.ExpectCommit()
			},
		},
		{
			name: "error",
			args: args{
				ctx: context.Background(),
				walletHistory: &models.WalletTransaction{
					WalletID:              1,
					Amount:                200000,
					WalletTransactionType: "DEBIT",
					Reference:             "reference",
				},
			},
			wantErr: true,
			mockFn: func(args args) {
				mock.ExpectBegin()
				mock.ExpectExec(regexp.QuoteMeta("INSERT INTO `wallet_transactions` (`wallet_id`,`amount`,`wallet_transaction_type`,`reference`,`created_at`,`updated_at`) VALUES (?,?,?,?,?,?)")).WithArgs(
					args.walletHistory.WalletID,
					args.walletHistory.Amount,
					args.walletHistory.WalletTransactionType,
					args.walletHistory.Reference,
					sqlmock.AnyArg(),
					sqlmock.AnyArg(),
				).WillReturnError(assert.AnError)
				mock.ExpectRollback()
			},
		},
		{
			name: "error transaction type",
			args: args{
				ctx: context.Background(),
				walletHistory: &models.WalletTransaction{
					WalletID:              1,
					Amount:                200000,
					WalletTransactionType: "TEST",
					Reference:             "reference",
				},
			},
			wantErr: true,
			mockFn: func(args args) {
				mock.ExpectBegin()
				mock.ExpectExec(regexp.QuoteMeta("INSERT INTO `wallet_transactions` (`wallet_id`,`amount`,`wallet_transaction_type`,`reference`,`created_at`,`updated_at`) VALUES (?,?,?,?,?,?)")).WithArgs(
					args.walletHistory.WalletID,
					args.walletHistory.Amount,
					args.walletHistory.WalletTransactionType,
					args.walletHistory.Reference,
					sqlmock.AnyArg(),
					sqlmock.AnyArg(),
				).WillReturnError(assert.AnError)
				mock.ExpectRollback()
			},
		},
	}
	for _, tt := range tests {
		tt.mockFn(tt.args)
		t.Run(tt.name, func(t *testing.T) {
			r := &WalletRepo{
				DB: gormDB,
			}
			if err := r.CreateWalletTrx(tt.args.ctx, tt.args.walletHistory); (err != nil) != tt.wantErr {
				t.Errorf("WalletRepo.CreateWalletTrx() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
		assert.NoError(t, mock.ExpectationsWereMet())
	}
}

func TestWalletRepo_GetWalletTransactionByReference(t *testing.T) {
	db, mock, err := sqlmock.New()
	assert.NoError(t, err)
	defer db.Close()

	gormDB, err := gorm.Open(mysql.New(mysql.Config{
		Conn:                      db,
		SkipInitializeWithVersion: true,
	}), &gorm.Config{})

	assert.NoError(t, err)

	now := time.Now()

	type args struct {
		ctx       context.Context
		reference string
	}
	tests := []struct {
		name    string
		args    args
		want    models.WalletTransaction
		wantErr bool
		mockFn  func(args args)
	}{
		{
			name: "success",
			args: args{
				reference: "reference",
			},
			want: models.WalletTransaction{
				ID:                    1,
				WalletID:              1,
				Amount:                200000,
				WalletTransactionType: "DEBIT",
				Reference:             "reference",
				CreatedAt:             now,
				UpdatedAt:             now,
			},
			wantErr: false,
			mockFn: func(args args) {
				mock.ExpectQuery(regexp.QuoteMeta("SELECT * FROM `wallet_transactions` WHERE reference = ? ORDER BY `wallet_transactions`.`id` DESC LIMIT ?")).WithArgs(
					args.reference,
					1,
				).
					WillReturnRows(sqlmock.NewRows([]string{"id", "wallet_id", "amount", "wallet_transaction_type", "reference", "created_at", "updated_at"}).
						AddRow(1, 1, 200000, "DEBIT", "reference", now, now))
			},
		},
		{
			name: "error",
			args: args{
				reference: "reference",
			},
			want:    models.WalletTransaction{},
			wantErr: true,
			mockFn: func(args args) {
				mock.ExpectQuery(regexp.QuoteMeta("SELECT * FROM `wallet_transactions` WHERE reference = ? ORDER BY `wallet_transactions`.`id` DESC LIMIT ?")).WithArgs(
					args.reference,
					1,
				).
					WillReturnError(assert.AnError)
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.mockFn(tt.args)
			r := &WalletRepo{
				DB: gormDB,
			}
			got, err := r.GetWalletTransactionByReference(tt.args.ctx, tt.args.reference)
			if (err != nil) != tt.wantErr {
				t.Errorf("WalletRepo.GetWalletTransactionByReference() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("WalletRepo.GetWalletTransactionByReference() = %v, want %v", got, tt.want)
			}
			assert.NoError(t, mock.ExpectationsWereMet())

		})
	}
}

func TestWalletRepo_GetWalletByUserID(t *testing.T) {
	db, mock, err := sqlmock.New()
	assert.NoError(t, err)
	defer db.Close()

	gormDB, err := gorm.Open(mysql.New(mysql.Config{
		Conn:                      db,
		SkipInitializeWithVersion: true,
	}), &gorm.Config{})

	assert.NoError(t, err)

	now := time.Now()

	type args struct {
		ctx    context.Context
		userID uint64
	}
	tests := []struct {
		name    string
		args    args
		want    models.Wallet
		wantErr bool
		mockFn  func(args args)
	}{
		{
			name: "success",
			args: args{
				ctx:    context.Background(),
				userID: 1,
			},
			want: models.Wallet{
				ID:        1,
				UserID:    1,
				Balance:   200000,
				CreatedAt: now,
				UpdatedAt: now,
			},
			wantErr: false,
			mockFn: func(args args) {
				mock.ExpectQuery(regexp.QuoteMeta("SELECT * FROM `wallets` WHERE user_id = ? ORDER BY `wallets`.`id` DESC LIMIT ?")).WithArgs(
					args.userID,
					1,
				).WillReturnRows(sqlmock.NewRows([]string{"id", "user_id", "balance", "created_at", "updated_at"}).AddRow(1, 1, 200000, now, now))
			},
		},
		{
			name: "false",
			args: args{
				ctx:    context.Background(),
				userID: 1,
			},
			want:    models.Wallet{},
			wantErr: true,
			mockFn: func(args args) {
				mock.ExpectQuery(regexp.QuoteMeta("SELECT * FROM `wallets` WHERE user_id = ? ORDER BY `wallets`.`id` DESC LIMIT ?")).WithArgs(
					args.userID,
					1,
				).WillReturnError(assert.AnError)
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.mockFn(tt.args)
			r := &WalletRepo{
				DB: gormDB,
			}
			got, err := r.GetWalletByUserID(tt.args.ctx, tt.args.userID)
			if (err != nil) != tt.wantErr {
				t.Errorf("WalletRepo.GetWalletByUserID() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("WalletRepo.GetWalletByUserID() = %v, want %v", got, tt.want)
			}
			assert.NoError(t, mock.ExpectationsWereMet())
		})
	}
}

func TestWalletRepo_GetWalletHistory(t *testing.T) {
	db, mock, err := sqlmock.New()
	assert.NoError(t, err)
	defer db.Close()

	gormDB, err := gorm.Open(mysql.New(mysql.Config{
		Conn:                      db,
		SkipInitializeWithVersion: true,
	}), &gorm.Config{})

	assert.NoError(t, err)

	now := time.Now()
	type args struct {
		ctx             context.Context
		walletID        int
		offset          int
		limit           int
		transactionType string
	}
	tests := []struct {
		name    string
		args    args
		want    []models.WalletTransaction
		wantErr bool
		mockFn  func(args args)
	}{
		{
			name: "success",
			args: args{
				ctx:             context.Background(),
				walletID:        1,
				offset:          2,
				limit:           2,
				transactionType: "DEBIT",
			},
			want: []models.WalletTransaction{
				{
					ID:                    3,
					WalletID:              3,
					Amount:                300000,
					WalletTransactionType: "DEBIT",
					Reference:             "reference",
					CreatedAt:             now,
					UpdatedAt:             now,
				},
				{
					ID:                    5,
					WalletID:              5,
					Amount:                500000,
					WalletTransactionType: "DEBIT",
					Reference:             "reference",
					CreatedAt:             now,
					UpdatedAt:             now,
				},
			},
			wantErr: false,
			mockFn: func(args args) {
				mock.ExpectQuery(regexp.QuoteMeta("SELECT * FROM `wallet_transactions` WHERE wallet_transaction_type = ? ORDER BY id DESC LIMIT ? OFFSET ?")).WithArgs(
					args.transactionType,
					args.limit,
					args.offset,
				).WillReturnRows(mock.NewRows([]string{"id", "wallet_id", "amount", "wallet_transaction_type", "reference", "created_at", "updated_at"}).
					AddRow(3, 3, 300000, "DEBIT", "reference", now, now).AddRow(5, 5, 500000, "DEBIT", "reference", now, now))
			},
		},
		{
			name: "error",
			args: args{
				ctx:             context.Background(),
				walletID:        1,
				offset:          2,
				limit:           2,
				transactionType: "DEBIT",
			},
			want:    nil,
			wantErr: true,
			mockFn: func(args args) {
				mock.ExpectQuery(regexp.QuoteMeta("SELECT * FROM `wallet_transactions` WHERE wallet_transaction_type = ? ORDER BY id DESC LIMIT ? OFFSET ?")).WithArgs(
					args.transactionType,
					args.limit,
					args.offset,
				).WillReturnError(assert.AnError)
			},
		},
		{
			name: "success with no transaction type",
			args: args{
				ctx:             context.Background(),
				walletID:        1,
				offset:          2,
				limit:           2,
				transactionType: "",
			},
			want: []models.WalletTransaction{
				{
					ID:                    3,
					WalletID:              3,
					Amount:                300000,
					WalletTransactionType: "DEBIT",
					Reference:             "reference",
					CreatedAt:             now,
					UpdatedAt:             now,
				},
				{
					ID:                    4,
					WalletID:              4,
					Amount:                400000,
					WalletTransactionType: "CREDIT",
					Reference:             "reference",
					CreatedAt:             now,
					UpdatedAt:             now,
				},
				{
					ID:                    5,
					WalletID:              5,
					Amount:                500000,
					WalletTransactionType: "DEBIT",
					Reference:             "reference",
					CreatedAt:             now,
					UpdatedAt:             now,
				},
			},
			wantErr: false,
			mockFn: func(args args) {
				mock.ExpectQuery(regexp.QuoteMeta("SELECT * FROM `wallet_transactions` ORDER BY id DESC LIMIT ? OFFSET ?")).WithArgs(
					args.limit,
					args.offset,
				).WillReturnRows(mock.NewRows([]string{"id", "wallet_id", "amount", "wallet_transaction_type", "reference", "created_at", "updated_at"}).
					AddRow(3, 3, 300000, "DEBIT", "reference", now, now).AddRow(4, 4, 400000, "CREDIT", "reference", now, now).
					AddRow(5, 5, 500000, "DEBIT", "reference", now, now))
			},
		},
		{
			name: "error with no transaction type",
			args: args{
				ctx:             context.Background(),
				walletID:        1,
				offset:          2,
				limit:           2,
				transactionType: "",
			},
			want:    nil,
			wantErr: true,
			mockFn: func(args args) {
				mock.ExpectQuery(regexp.QuoteMeta("SELECT * FROM `wallet_transactions` ORDER BY id DESC LIMIT ? OFFSET ?")).WithArgs(
					args.limit,
					args.offset,
				).WillReturnError(assert.AnError)
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.mockFn(tt.args)
			r := &WalletRepo{
				DB: gormDB,
			}
			got, err := r.GetWalletHistory(tt.args.ctx, tt.args.walletID, tt.args.offset, tt.args.limit, tt.args.transactionType)
			if (err != nil) != tt.wantErr {
				t.Errorf("WalletRepo.GetWalletHistory() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("WalletRepo.GetWalletHistory() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestWalletRepo_InsertWalletLink(t *testing.T) {
	db, mock, err := sqlmock.New()
	assert.NoError(t, err)
	defer db.Close()

	gormDB, err := gorm.Open(mysql.New(mysql.Config{
		Conn:                      db,
		SkipInitializeWithVersion: true,
	}), &gorm.Config{})

	assert.NoError(t, err)

	type args struct {
		ctx context.Context
		req *models.WalletLink
	}
	tests := []struct {
		name    string
		args    args
		wantErr bool
		mockFn  func(args args)
	}{
		{
			name: "success",
			args: args{
				ctx: context.Background(),
				req: &models.WalletLink{
					WalletID:     1,
					ClientSource: "fastcampus_wallet",
					OTP:          "098890",
					Status:       "pending",
				},
			},
			wantErr: false,
			mockFn: func(args args) {
				mock.ExpectBegin()
				mock.ExpectExec(regexp.QuoteMeta("INSERT INTO `wallet_links` (`wallet_id`,`client_source`,`otp`,`status`,`created_at`,`updated_at`) VALUES (?,?,?,?,?,?)")).WithArgs(
					args.req.WalletID,
					args.req.ClientSource,
					args.req.OTP,
					args.req.Status,
					sqlmock.AnyArg(),
					sqlmock.AnyArg(),
				).WillReturnResult(sqlmock.NewResult(1, 1))
				mock.ExpectCommit()
			},
		},
		{
			name: "error",
			args: args{
				ctx: context.Background(),
				req: &models.WalletLink{
					WalletID:     1,
					ClientSource: "fastcampus_wallet",
					OTP:          "098890",
					Status:       "pending",
				},
			},
			wantErr: true,
			mockFn: func(args args) {
				mock.ExpectBegin()
				mock.ExpectExec(regexp.QuoteMeta("INSERT INTO `wallet_links` (`wallet_id`,`client_source`,`otp`,`status`,`created_at`,`updated_at`) VALUES (?,?,?,?,?,?)")).WithArgs(
					args.req.WalletID,
					args.req.ClientSource,
					args.req.OTP,
					args.req.Status,
					sqlmock.AnyArg(),
					sqlmock.AnyArg(),
				).WillReturnError(assert.AnError)
				mock.ExpectRollback()
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.mockFn(tt.args)
			r := &WalletRepo{
				DB: gormDB,
			}
			if err := r.InsertWalletLink(tt.args.ctx, tt.args.req); (err != nil) != tt.wantErr {
				t.Errorf("WalletRepo.InsertWalletLink() error = %v, wantErr %v", err, tt.wantErr)
			}
			assert.NoError(t, mock.ExpectationsWereMet())
		})
	}
}

func TestWalletRepo_GetWalletLink(t *testing.T) {
	db, mock, err := sqlmock.New()
	assert.NoError(t, err)

	gormDB, err := gorm.Open(mysql.New(mysql.Config{
		Conn:                      db,
		SkipInitializeWithVersion: true,
	}), &gorm.Config{})

	assert.NoError(t, err)

	defer db.Close()
	now := time.Now()

	type args struct {
		ctx          context.Context
		walletID     int
		clientSource string
	}
	tests := []struct {
		name    string
		args    args
		want    models.WalletLink
		wantErr bool
		mockFn  func(args args)
	}{
		{
			name: "success",
			args: args{
				ctx:          context.Background(),
				walletID:     1,
				clientSource: "fastcampus_wallet",
			},
			want: models.WalletLink{
				ID:           1,
				WalletID:     1,
				ClientSource: "fastcampus_wallet",
				OTP:          "909029",
				Status:       "pending",
				CreatedAt:    now,
				UpdatedAt:    now,
			},
			wantErr: false,
			mockFn: func(args args) {
				mock.ExpectQuery(regexp.QuoteMeta("SELECT * FROM `wallet_links` WHERE wallet_id = ? AND client_source = ? ORDER BY `wallet_links`.`id` LIMIT ?")).WithArgs(
					args.walletID,
					args.clientSource,
					1,
				).WillReturnRows(sqlmock.NewRows([]string{"id", "wallet_id", "client_source", "otp", "status", "created_at", "updated_at"}).AddRow(1, 1, "fastcampus_wallet", "909029", "pending", now, now))
			},
		},
		{
			name: "error",
			args: args{
				ctx:          context.Background(),
				walletID:     1,
				clientSource: "fastcampus_wallet",
			},
			want:    models.WalletLink{},
			wantErr: true,
			mockFn: func(args args) {
				mock.ExpectQuery(regexp.QuoteMeta("SELECT * FROM `wallet_links` WHERE wallet_id = ? AND client_source = ? ORDER BY `wallet_links`.`id` LIMIT ?")).WithArgs(
					args.walletID,
					args.clientSource,
					1,
				).WillReturnError(assert.AnError)
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.mockFn(tt.args)
			r := &WalletRepo{
				DB: gormDB,
			}
			got, err := r.GetWalletLink(tt.args.ctx, tt.args.walletID, tt.args.clientSource)
			if (err != nil) != tt.wantErr {
				t.Errorf("WalletRepo.GetWalletLink() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("WalletRepo.GetWalletLink() = %v, want %v", got, tt.want)
			}
			assert.NoError(t, mock.ExpectationsWereMet())
		})
	}
}

func TestWalletRepo_UpdateStatusWalletLink(t *testing.T) {
	db, mock, err := sqlmock.New()
	assert.NoError(t, err)
	defer db.Close()

	gormDB, err := gorm.Open(mysql.New(mysql.Config{
		Conn:                      db,
		SkipInitializeWithVersion: true,
	}), &gorm.Config{})

	assert.NoError(t, err)

	type args struct {
		ctx          context.Context
		walletID     int
		clientSource string
		status       string
	}
	tests := []struct {
		name    string
		args    args
		wantErr bool
		mockFn  func(args args)
	}{
		{
			name: "success",
			args: args{
				ctx:          context.Background(),
				walletID:     1,
				clientSource: "fastcampus_wallet",
				status:       "success",
			},
			wantErr: false,
			mockFn: func(args args) {
				mock.ExpectExec(regexp.QuoteMeta("UPDATE wallet_links SET status = ? WHERE wallet_id = ? AND client_source = ?")).WithArgs(
					args.status,
					args.walletID,
					args.clientSource,
				).WillReturnResult(sqlmock.NewResult(1, 1))
			},
		},
		{
			name: "error",
			args: args{
				ctx:          context.Background(),
				walletID:     1,
				clientSource: "fastcampus_wallet",
				status:       "success",
			},
			wantErr: true,
			mockFn: func(args args) {
				mock.ExpectExec(regexp.QuoteMeta("UPDATE wallet_links SET status = ? WHERE wallet_id = ? AND client_source = ?")).WithArgs(
					args.status,
					args.walletID,
					args.clientSource,
				).WillReturnError(assert.AnError)
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.mockFn(tt.args)
			r := &WalletRepo{
				DB: gormDB,
			}
			if err := r.UpdateStatusWalletLink(tt.args.ctx, tt.args.walletID, tt.args.clientSource, tt.args.status); (err != nil) != tt.wantErr {
				t.Errorf("WalletRepo.UpdateStatusWalletLink() error = %v, wantErr %v", err, tt.wantErr)
			}
			assert.NoError(t, mock.ExpectationsWereMet())
		})
	}
}

func TestWalletRepo_GetWalletByID(t *testing.T) {
	db, mock, err := sqlmock.New()
	assert.NoError(t, err)
	defer db.Close()

	gormDB, err := gorm.Open(mysql.New(mysql.Config{
		Conn:                      db,
		SkipInitializeWithVersion: true,
	}), &gorm.Config{})

	assert.NoError(t, err)

	now := time.Now()

	type args struct {
		ctx      context.Context
		walletID int
	}
	tests := []struct {
		name    string
		args    args
		want    models.Wallet
		wantErr bool
		mockFn  func(args args)
	}{
		{
			name: "success",
			args: args{
				ctx:      context.Background(),
				walletID: 1,
			},
			want: models.Wallet{
				ID:        1,
				UserID:    1,
				Balance:   200000,
				CreatedAt: now,
				UpdatedAt: now,
			},
			wantErr: false,
			mockFn: func(args args) {
				mock.ExpectQuery(regexp.QuoteMeta("SELECT * FROM `wallets` WHERE id = ? ORDER BY `wallets`.`id` DESC LIMIT ?")).WithArgs(
					args.walletID,
					1,
				).WillReturnRows(sqlmock.NewRows([]string{"id", "user_id", "balance", "created_at", "updated_at"}).AddRow(1, 1, 200000, now, now))
			},
		},
		{
			name: "error",
			args: args{
				ctx:      context.Background(),
				walletID: 1,
			},
			want:    models.Wallet{},
			wantErr: true,
			mockFn: func(args args) {
				mock.ExpectQuery(regexp.QuoteMeta("SELECT * FROM `wallets` WHERE id = ? ORDER BY `wallets`.`id` DESC LIMIT ?")).WithArgs(
					args.walletID,
					1,
				).WillReturnError(assert.AnError)
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.mockFn(tt.args)
			r := &WalletRepo{
				DB: gormDB,
			}
			got, err := r.GetWalletByID(tt.args.ctx, tt.args.walletID)
			if (err != nil) != tt.wantErr {
				t.Errorf("WalletRepo.GetWalletByID() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("WalletRepo.GetWalletByID() = %v, want %v", got, tt.want)
			}
			assert.NoError(t, mock.ExpectationsWereMet())
		})
	}
}

func TestWalletRepo_UpdateBalanceByID(t *testing.T) {
	db, mock, err := sqlmock.New()
	assert.NoError(t, err)
	defer db.Close()

	gormDB, err := gorm.Open(mysql.New(mysql.Config{
		Conn:                      db,
		SkipInitializeWithVersion: true,
	}), &gorm.Config{})

	assert.NoError(t, err)

	type args struct {
		ctx      context.Context
		walletID int
		amount   float64
	}
	tests := []struct {
		name    string
		args    args
		want    models.Wallet
		wantErr bool
		mockFn  func(args args)
	}{
		{
			name: "success add balance",
			args: args{
				ctx:      context.Background(),
				walletID: 1,
				amount:   50000,
			},
			want: models.Wallet{
				ID:      1,
				Balance: 200000,
			},
			wantErr: false,
			mockFn: func(args args) {
				mock.ExpectBegin()

				mock.ExpectQuery(regexp.QuoteMeta("SELECT id, balance FROM wallets WHERE id = ? FOR UPDATE")).WillReturnRows(sqlmock.NewRows([]string{"id", "balance"}).AddRow(1, 200000))

				mock.ExpectExec(regexp.QuoteMeta("UPDATE wallets SET balance = balance + ? WHERE id = ?")).WithArgs(
					args.amount,
					args.walletID,
				).WillReturnResult(sqlmock.NewResult(1, 1))

				mock.ExpectCommit()
			},
		},
		{
			name: "success decrease balance",
			args: args{
				ctx:      context.Background(),
				walletID: 1,
				amount:   -50000,
			},
			want: models.Wallet{
				ID:      1,
				Balance: 200000,
			},
			wantErr: false,
			mockFn: func(args args) {
				mock.ExpectBegin()

				mock.ExpectQuery(regexp.QuoteMeta("SELECT id, balance FROM wallets WHERE id = ? FOR UPDATE")).WillReturnRows(sqlmock.NewRows([]string{"id", "balance"}).AddRow(1, 200000))

				mock.ExpectExec(regexp.QuoteMeta("UPDATE wallets SET balance = balance + ? WHERE id = ?")).WithArgs(
					args.amount,
					args.walletID,
				).WillReturnResult(sqlmock.NewResult(1, 1))

				mock.ExpectCommit()
			},
		},
		{
			name: "error get wallet",
			args: args{
				ctx:      context.Background(),
				walletID: 1,
				amount:   250000,
			},
			want:    models.Wallet{},
			wantErr: true,
			mockFn: func(args args) {
				mock.ExpectBegin()

				mock.ExpectQuery(regexp.QuoteMeta("SELECT id, balance FROM wallets WHERE id = ? FOR UPDATE")).WillReturnError(assert.AnError)

				mock.ExpectRollback()
			},
		},
		{
			name: "err updated balance",
			args: args{
				ctx:      context.Background(),
				walletID: 1,
				amount:   50000,
			},
			want: models.Wallet{
				ID:      1,
				Balance: 200000,
			},
			wantErr: true,
			mockFn: func(args args) {
				mock.ExpectBegin()

				mock.ExpectQuery(regexp.QuoteMeta("SELECT id, balance FROM wallets WHERE id = ? FOR UPDATE")).WillReturnRows(sqlmock.NewRows([]string{"id", "balance"}).AddRow(1, 200000))

				mock.ExpectExec(regexp.QuoteMeta("UPDATE wallets SET balance = balance + ? WHERE id = ?")).WithArgs(
					args.amount,
					args.walletID,
				).WillReturnError(assert.AnError)

				mock.ExpectRollback()
			},
		},
		{
			name: "error decrease balance",
			args: args{
				ctx:      context.Background(),
				walletID: 1,
				amount:   -250000,
			},
			want: models.Wallet{
				ID:      1,
				Balance: 200000,
			},
			wantErr: true,
			mockFn: func(args args) {
				mock.ExpectBegin()

				mock.ExpectQuery(regexp.QuoteMeta("SELECT id, balance FROM wallets WHERE id = ? FOR UPDATE")).WithArgs(
					args.walletID,
				).WillReturnRows(sqlmock.NewRows([]string{"id", "balance"}).AddRow(1, 200000))

				mock.ExpectRollback()
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.mockFn(tt.args)
			r := &WalletRepo{
				DB: gormDB,
			}
			got, err := r.UpdateBalanceByID(tt.args.ctx, tt.args.walletID, tt.args.amount)
			if (err != nil) != tt.wantErr {
				t.Errorf("WalletRepo.UpdateBalanceByID() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("WalletRepo.UpdateBalanceByID() = %v, want %v", got, tt.want)
			}
			assert.NoError(t, mock.ExpectationsWereMet())
		})
	}
}
