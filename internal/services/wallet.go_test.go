package services

import (
	"context"
	"ewallet-wallet/internal/models"
	"reflect"
	"testing"
	"time"

	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/assert"
)

func TestWalletService_Create(t *testing.T) {
	ctrlMock := gomock.NewController(t)
	defer ctrlMock.Finish()

	mockRepo := NewMockIWalletRepo(ctrlMock)

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
			name: "success",
			args: args{
				ctx: context.Background(),
				wallet: &models.Wallet{
					UserID:  1,
					Balance: 200000,
				},
			},
			wantErr: false,
			mockFn: func(args args) {
				mockRepo.EXPECT().CreateWallet(args.ctx, args.wallet).Return(nil)
			},
		},
		{
			name: "error",
			args: args{
				ctx: context.Background(),
				wallet: &models.Wallet{
					UserID:  1,
					Balance: 200000,
				},
			},
			wantErr: true,
			mockFn: func(args args) {
				mockRepo.EXPECT().CreateWallet(args.ctx, args.wallet).Return(assert.AnError)
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.mockFn(tt.args)
			s := &WalletService{
				WalletRepo: mockRepo,
			}
			if err := s.Create(tt.args.ctx, tt.args.wallet); (err != nil) != tt.wantErr {
				t.Errorf("WalletService.Create() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestWalletService_CreditBalance(t *testing.T) {
	ctrlMock := gomock.NewController(t)
	defer ctrlMock.Finish()

	mockRepo := NewMockIWalletRepo(ctrlMock)

	now := time.Now()

	type args struct {
		ctx    context.Context
		userID uint64
		req    models.TransactionRequest
	}
	tests := []struct {
		name    string
		args    args
		want    models.BalanceResponse
		wantErr bool
		mockFn  func(args args)
	}{
		{
			name: "success",
			args: args{
				ctx:    context.Background(),
				userID: 1,
				req: models.TransactionRequest{
					Reference: "reference",
					Amount:    100000,
				},
			},
			want: models.BalanceResponse{
				Balance: 300000,
			},
			wantErr: false,
			mockFn: func(args args) {
				mockRepo.EXPECT().GetWalletTransactionByReference(args.ctx, args.req.Reference).Return(models.WalletTransaction{}, nil)

				mockRepo.EXPECT().UpdateBalance(args.ctx, args.userID, args.req.Amount).Return(models.Wallet{
					ID:        1,
					UserID:    1,
					Balance:   200000,
					CreatedAt: now,
					UpdatedAt: now,
				}, nil)

				mockRepo.EXPECT().CreateWalletTrx(args.ctx, &models.WalletTransaction{
					WalletID:              1,
					Amount:                args.req.Amount,
					WalletTransactionType: "CREDIT",
					Reference:             args.req.Reference,
				}).Return(nil)
			},
		},
		{
			name: "error, got wallet with duplicate reference",
			args: args{
				ctx:    context.Background(),
				userID: 1,
				req: models.TransactionRequest{
					Reference: "reference",
					Amount:    100000,
				},
			},
			want:    models.BalanceResponse{},
			wantErr: true,
			mockFn: func(args args) {
				mockRepo.EXPECT().GetWalletTransactionByReference(args.ctx, args.req.Reference).Return(models.WalletTransaction{
					ID:                    1,
					WalletID:              1,
					Amount:                100000,
					Reference:             "reference",
					WalletTransactionType: "CREDIT",
					CreatedAt:             now,
					UpdatedAt:             now,
				}, nil)
			},
		},
		{
			name: "error update balance",
			args: args{
				ctx:    context.Background(),
				userID: 1,
				req: models.TransactionRequest{
					Reference: "reference",
					Amount:    100000,
				},
			},
			want:    models.BalanceResponse{},
			wantErr: true,
			mockFn: func(args args) {
				mockRepo.EXPECT().GetWalletTransactionByReference(args.ctx, args.req.Reference).Return(models.WalletTransaction{}, nil)

				mockRepo.EXPECT().UpdateBalance(args.ctx, args.userID, args.req.Amount).Return(models.Wallet{}, assert.AnError)
			},
		},
		{
			name: "error create transaction",
			args: args{
				ctx:    context.Background(),
				userID: 1,
				req: models.TransactionRequest{
					Reference: "reference",
					Amount:    100000,
				},
			},
			want:    models.BalanceResponse{},
			wantErr: true,
			mockFn: func(args args) {
				mockRepo.EXPECT().GetWalletTransactionByReference(args.ctx, args.req.Reference).Return(models.WalletTransaction{}, nil)

				mockRepo.EXPECT().UpdateBalance(args.ctx, args.userID, args.req.Amount).Return(models.Wallet{
					ID:        1,
					UserID:    1,
					Balance:   200000,
					CreatedAt: now,
					UpdatedAt: now,
				}, nil)

				mockRepo.EXPECT().CreateWalletTrx(args.ctx, &models.WalletTransaction{
					WalletID:              1,
					Amount:                args.req.Amount,
					WalletTransactionType: "CREDIT",
					Reference:             args.req.Reference,
				}).Return(assert.AnError)
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.mockFn(tt.args)

			s := &WalletService{
				WalletRepo: mockRepo,
			}
			got, err := s.CreditBalance(tt.args.ctx, tt.args.userID, tt.args.req)
			if (err != nil) != tt.wantErr {
				t.Errorf("WalletService.CreditBalance() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("WalletService.CreditBalance() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestWalletService_DebitBalance(t *testing.T) {
	ctrlMock := gomock.NewController(t)
	defer ctrlMock.Finish()

	mockRepo := NewMockIWalletRepo(ctrlMock)

	now := time.Now()

	type args struct {
		ctx    context.Context
		userID uint64
		req    models.TransactionRequest
	}
	tests := []struct {
		name    string
		args    args
		want    models.BalanceResponse
		wantErr bool
		mockFn  func(args args)
	}{
		{
			name: "success",
			args: args{
				ctx:    context.Background(),
				userID: 1,
				req: models.TransactionRequest{
					Reference: "reference",
					Amount:    100000,
				},
			},
			want: models.BalanceResponse{
				Balance: 100000,
			},
			wantErr: false,
			mockFn: func(args args) {
				mockRepo.EXPECT().GetWalletTransactionByReference(args.ctx, args.req.Reference).Return(models.WalletTransaction{}, nil)

				mockRepo.EXPECT().UpdateBalance(args.ctx, args.userID, (-args.req.Amount)).Return(models.Wallet{
					ID:        1,
					UserID:    1,
					Balance:   200000,
					CreatedAt: now,
					UpdatedAt: now,
				}, nil)

				mockRepo.EXPECT().CreateWalletTrx(args.ctx, &models.WalletTransaction{
					WalletID:              1,
					Amount:                args.req.Amount,
					WalletTransactionType: "DEBIT",
					Reference:             args.req.Reference,
				}).Return(nil)
			},
		},
		{
			name: "error, got wallet with duplicate reference",
			args: args{
				ctx:    context.Background(),
				userID: 1,
				req: models.TransactionRequest{
					Reference: "reference",
					Amount:    100000,
				},
			},
			want:    models.BalanceResponse{},
			wantErr: true,
			mockFn: func(args args) {
				mockRepo.EXPECT().GetWalletTransactionByReference(args.ctx, args.req.Reference).Return(models.WalletTransaction{
					ID:                    1,
					WalletID:              1,
					Amount:                100000,
					Reference:             "reference",
					WalletTransactionType: "CREDIT",
					CreatedAt:             now,
					UpdatedAt:             now,
				}, nil)
			},
		},
		{
			name: "error update balance",
			args: args{
				ctx:    context.Background(),
				userID: 1,
				req: models.TransactionRequest{
					Reference: "reference",
					Amount:    500000,
				},
			},
			want:    models.BalanceResponse{},
			wantErr: true,
			mockFn: func(args args) {
				mockRepo.EXPECT().GetWalletTransactionByReference(args.ctx, args.req.Reference).Return(models.WalletTransaction{}, nil)

				mockRepo.EXPECT().UpdateBalance(args.ctx, args.userID, (-args.req.Amount)).Return(models.Wallet{}, assert.AnError)
			},
		},
		{
			name: "error create transaction",
			args: args{
				ctx:    context.Background(),
				userID: 1,
				req: models.TransactionRequest{
					Reference: "reference",
					Amount:    100000,
				},
			},
			want:    models.BalanceResponse{},
			wantErr: true,
			mockFn: func(args args) {
				mockRepo.EXPECT().GetWalletTransactionByReference(args.ctx, args.req.Reference).Return(models.WalletTransaction{}, nil)

				mockRepo.EXPECT().UpdateBalance(args.ctx, args.userID, (-args.req.Amount)).Return(models.Wallet{
					ID:        1,
					UserID:    1,
					Balance:   200000,
					CreatedAt: now,
					UpdatedAt: now,
				}, nil)

				mockRepo.EXPECT().CreateWalletTrx(args.ctx, &models.WalletTransaction{
					WalletID:              1,
					Amount:                args.req.Amount,
					WalletTransactionType: "DEBIT",
					Reference:             args.req.Reference,
				}).Return(assert.AnError)
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.mockFn(tt.args)
			s := &WalletService{
				WalletRepo: mockRepo,
			}
			got, err := s.DebitBalance(tt.args.ctx, tt.args.userID, tt.args.req)
			if (err != nil) != tt.wantErr {
				t.Errorf("WalletService.DebitBalance() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("WalletService.DebitBalance() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestWalletService_GetBalance(t *testing.T) {
	ctrlMock := gomock.NewController(t)
	defer ctrlMock.Finish()

	mockRepo := NewMockIWalletRepo(ctrlMock)

	now := time.Now()
	type args struct {
		ctx    context.Context
		userID uint64
	}
	tests := []struct {
		name    string
		args    args
		want    models.BalanceResponse
		wantErr bool
		mockFn  func(args args)
	}{
		{
			name: "success",
			args: args{
				ctx:    context.Background(),
				userID: 1,
			},
			want: models.BalanceResponse{
				Balance: 200000,
			},
			wantErr: false,
			mockFn: func(args args) {
				mockRepo.EXPECT().GetWalletByUserID(args.ctx, args.userID).Return(models.Wallet{
					ID:        1,
					UserID:    1,
					Balance:   200000,
					CreatedAt: now,
					UpdatedAt: now,
				}, nil)
			},
		},
		{
			name: "error",
			args: args{
				ctx:    context.Background(),
				userID: 1,
			},
			want:    models.BalanceResponse{},
			wantErr: true,
			mockFn: func(args args) {
				mockRepo.EXPECT().GetWalletByUserID(args.ctx, args.userID).Return(models.Wallet{}, assert.AnError)
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.mockFn(tt.args)
			s := &WalletService{
				WalletRepo: mockRepo,
			}
			got, err := s.GetBalance(tt.args.ctx, tt.args.userID)
			if (err != nil) != tt.wantErr {
				t.Errorf("WalletService.GetBalance() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("WalletService.GetBalance() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestWalletService_ExGetBalance(t *testing.T) {
	ctrlMock := gomock.NewController(t)
	defer ctrlMock.Finish()

	mockRepo := NewMockIWalletRepo(ctrlMock)
	type args struct {
		ctx      context.Context
		walletID int
	}

	now := time.Now()
	tests := []struct {
		name    string
		args    args
		want    models.BalanceResponse
		wantErr bool
		mockFn  func(args args)
	}{
		{
			name: "success",
			args: args{
				ctx:      context.Background(),
				walletID: 1,
			},
			want: models.BalanceResponse{
				Balance: 200000,
			},
			wantErr: false,
			mockFn: func(args args) {
				mockRepo.EXPECT().GetWalletByID(args.ctx, args.walletID).Return(models.Wallet{
					ID:        1,
					UserID:    1,
					Balance:   200000,
					CreatedAt: now,
					UpdatedAt: now,
				}, nil)
			},
		},
		{
			name: "error",
			args: args{
				ctx:      context.Background(),
				walletID: 1,
			},
			want:    models.BalanceResponse{},
			wantErr: true,
			mockFn: func(args args) {
				mockRepo.EXPECT().GetWalletByID(args.ctx, args.walletID).Return(models.Wallet{}, assert.AnError)
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.mockFn(tt.args)
			s := &WalletService{
				WalletRepo: mockRepo,
			}
			got, err := s.ExGetBalance(tt.args.ctx, tt.args.walletID)
			if (err != nil) != tt.wantErr {
				t.Errorf("WalletService.ExGetBalance() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("WalletService.ExGetBalance() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestWalletService_GetWalletHistory(t *testing.T) {
	ctrlMock := gomock.NewController(t)
	defer ctrlMock.Finish()

	mockRepo := NewMockIWalletRepo(ctrlMock)

	now := time.Now()
	type args struct {
		ctx    context.Context
		userID uint64
		param  models.WalletHistoryParam
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
				ctx:    context.Background(),
				userID: 1,
				param: models.WalletHistoryParam{
					Page:                  2,
					Limit:                 2,
					WalletTransactionType: "DEBIT",
				},
			},
			want: []models.WalletTransaction{
				{
					ID:                    2,
					WalletID:              2,
					Amount:                200000,
					WalletTransactionType: "DEBIT",
					Reference:             "reference2",
					CreatedAt:             now,
					UpdatedAt:             now,
				},
				{
					ID:                    3,
					WalletID:              3,
					Amount:                300000,
					WalletTransactionType: "DEBIT",
					Reference:             "reference2",
					CreatedAt:             now,
					UpdatedAt:             now,
				},
			},
			wantErr: false,
			mockFn: func(args args) {
				wallet := models.Wallet{
					ID:        1,
					UserID:    1,
					Balance:   200000,
					CreatedAt: now,
					UpdatedAt: now,
				}

				mockRepo.EXPECT().GetWalletByUserID(args.ctx, args.userID).Return(wallet, nil)

				offset := (args.param.Page - 1) * args.param.Limit

				mockRepo.EXPECT().GetWalletHistory(args.ctx, wallet.ID, offset, args.param.Limit, args.param.WalletTransactionType).Return([]models.WalletTransaction{
					{
						ID:                    2,
						WalletID:              2,
						Amount:                200000,
						WalletTransactionType: "DEBIT",
						Reference:             "reference2",
						CreatedAt:             now,
						UpdatedAt:             now,
					},
					{
						ID:                    3,
						WalletID:              3,
						Amount:                300000,
						WalletTransactionType: "DEBIT",
						Reference:             "reference2",
						CreatedAt:             now,
						UpdatedAt:             now,
					},
				}, nil)
			},
		},
		{
			name: "error get wallet",
			args: args{
				ctx:    context.Background(),
				userID: 1,
				param: models.WalletHistoryParam{
					Page:                  2,
					Limit:                 2,
					WalletTransactionType: "DEBIT",
				},
			},
			want:    nil,
			wantErr: true,
			mockFn: func(args args) {
				mockRepo.EXPECT().GetWalletByUserID(args.ctx, args.userID).Return(models.Wallet{}, assert.AnError)
			},
		},
		{
			name: "error get wallet history",
			args: args{
				ctx:    context.Background(),
				userID: 1,
				param: models.WalletHistoryParam{
					Page:                  2,
					Limit:                 2,
					WalletTransactionType: "DEBIT",
				},
			},
			want:    nil,
			wantErr: true,
			mockFn: func(args args) {
				wallet := models.Wallet{
					ID:        1,
					UserID:    1,
					Balance:   200000,
					CreatedAt: now,
					UpdatedAt: now,
				}

				mockRepo.EXPECT().GetWalletByUserID(args.ctx, args.userID).Return(wallet, nil)

				offset := (args.param.Page - 1) * args.param.Limit

				mockRepo.EXPECT().GetWalletHistory(args.ctx, wallet.ID, offset, args.param.Limit, args.param.WalletTransactionType).Return(nil, assert.AnError)
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.mockFn(tt.args)
			s := &WalletService{
				WalletRepo: mockRepo,
			}
			got, err := s.GetWalletHistory(tt.args.ctx, tt.args.userID, tt.args.param)
			if (err != nil) != tt.wantErr {
				t.Errorf("WalletService.GetWalletHistory() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("WalletService.GetWalletHistory() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestWalletService_CreateWalletLink(t *testing.T) {
	ctrlMock := gomock.NewController(t)
	defer ctrlMock.Finish()

	mockRepo := NewMockIWalletRepo(ctrlMock)

	clientSource := "fastcampus_wallet"

	type args struct {
		ctx          context.Context
		clientSource string
		req          *models.WalletLink
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
				clientSource: clientSource,
				req: &models.WalletLink{
					WalletID:     1,
					ClientSource: clientSource,
				},
			},
			wantErr: false,
			mockFn: func(args args) {
				mockRepo.EXPECT().InsertWalletLink(args.ctx, gomock.Any()).Return(nil)
			},
		},
		{
			name: "error",
			args: args{
				ctx:          context.Background(),
				clientSource: clientSource,
				req: &models.WalletLink{
					WalletID:     1,
					ClientSource: clientSource,
				},
			},
			wantErr: true,
			mockFn: func(args args) {
				mockRepo.EXPECT().InsertWalletLink(args.ctx, gomock.Any()).Return(assert.AnError)
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.mockFn(tt.args)
			s := &WalletService{
				WalletRepo: mockRepo,
			}
			got, err := s.CreateWalletLink(tt.args.ctx, tt.args.clientSource, tt.args.req)
			if (err != nil) != tt.wantErr {
				t.Errorf("WalletService.CreateWalletLink() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !tt.wantErr {
				assert.NotEmpty(t, got.OTP)
			} else {
				// fmt.Printf("test case: %s\n", tt.name)
				assert.Empty(t, got)
			}

		})
	}
}

func TestWalletService_WalletLinkConfirmation(t *testing.T) {
	ctrlMock := gomock.NewController(t)
	defer ctrlMock.Finish()

	mockRepo := NewMockIWalletRepo(ctrlMock)

	clientSource := "fastcampus_wallet"
	now := time.Now()
	type args struct {
		ctx          context.Context
		walletID     int
		clientSource string
		otp          string
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
				clientSource: clientSource,
				otp:          "121212",
			},
			wantErr: false,
			mockFn: func(args args) {
				mockRepo.EXPECT().GetWalletLink(args.ctx, args.walletID, clientSource).Return(models.WalletLink{
					ID:           1,
					WalletID:     1,
					ClientSource: clientSource,
					OTP:          args.otp,
					Status:       "pending",
					CreatedAt:    now,
					UpdatedAt:    now,
				}, nil)

				mockRepo.EXPECT().UpdateStatusWalletLink(args.ctx, args.walletID, clientSource, "linked").Return(nil)
			},
		},
		{
			name: "error get wallet",
			args: args{
				ctx:          context.Background(),
				walletID:     1,
				clientSource: clientSource,
				otp:          "121212",
			},
			wantErr: true,
			mockFn: func(args args) {
				mockRepo.EXPECT().GetWalletLink(args.ctx, args.walletID, clientSource).Return(models.WalletLink{}, assert.AnError)

			},
		},
		{
			name: "error when wallet link status not pending",
			args: args{
				ctx:          context.Background(),
				walletID:     1,
				clientSource: clientSource,
				otp:          "121212",
			},
			wantErr: true,
			mockFn: func(args args) {
				mockRepo.EXPECT().GetWalletLink(args.ctx, args.walletID, clientSource).Return(models.WalletLink{
					ID:           1,
					WalletID:     1,
					ClientSource: clientSource,
					OTP:          args.otp,
					Status:       "linked",
					CreatedAt:    now,
					UpdatedAt:    now,
				}, nil)
			},
		},
		{
			name: "error when otp invalid",
			args: args{
				ctx:          context.Background(),
				walletID:     1,
				clientSource: clientSource,
				otp:          "121212",
			},
			wantErr: true,
			mockFn: func(args args) {
				mockRepo.EXPECT().GetWalletLink(args.ctx, args.walletID, clientSource).Return(models.WalletLink{
					ID:           1,
					WalletID:     1,
					ClientSource: clientSource,
					OTP:          "3343",
					Status:       "linked",
					CreatedAt:    now,
					UpdatedAt:    now,
				}, nil)
			},
		},
		{
			name: "error when update status",
			args: args{
				ctx:          context.Background(),
				walletID:     1,
				clientSource: clientSource,
				otp:          "121212",
			},
			wantErr: true,
			mockFn: func(args args) {
				mockRepo.EXPECT().GetWalletLink(args.ctx, args.walletID, clientSource).Return(models.WalletLink{
					ID:           1,
					WalletID:     1,
					ClientSource: clientSource,
					OTP:          args.otp,
					Status:       "pending",
					CreatedAt:    now,
					UpdatedAt:    now,
				}, nil)

				mockRepo.EXPECT().UpdateStatusWalletLink(args.ctx, args.walletID, clientSource, "linked").Return(assert.AnError)
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.mockFn(tt.args)
			s := &WalletService{
				WalletRepo: mockRepo,
			}
			if err := s.WalletLinkConfirmation(tt.args.ctx, tt.args.walletID, tt.args.clientSource, tt.args.otp); (err != nil) != tt.wantErr {
				t.Errorf("WalletService.WalletLinkConfirmation() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestWalletService_WalletUnlink(t *testing.T) {
	ctrlMock := gomock.NewController(t)
	defer ctrlMock.Finish()

	mockRepo := NewMockIWalletRepo(ctrlMock)

	clientSource := "fastcampus_wallet"
	type args struct {
		ctx          context.Context
		walletID     int
		clientSource string
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
				clientSource: clientSource,
			},
			wantErr: false,
			mockFn: func(args args) {
				mockRepo.EXPECT().UpdateStatusWalletLink(args.ctx, args.walletID, clientSource, "unlinked").Return(nil)
			},
		},
		{
			name: "false",
			args: args{
				ctx:          context.Background(),
				walletID:     1,
				clientSource: clientSource,
			},
			wantErr: true,
			mockFn: func(args args) {
				mockRepo.EXPECT().UpdateStatusWalletLink(args.ctx, args.walletID, clientSource, "unlinked").Return(assert.AnError)
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.mockFn(tt.args)
			s := &WalletService{
				WalletRepo: mockRepo,
			}
			if err := s.WalletUnlink(tt.args.ctx, tt.args.walletID, tt.args.clientSource); (err != nil) != tt.wantErr {
				t.Errorf("WalletService.WalletUnlink() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestWalletService_ExternalTransaction(t *testing.T) {
	ctrlMock := gomock.NewController(t)
	defer ctrlMock.Finish()

	mockRepo := NewMockIWalletRepo(ctrlMock)

	now := time.Now()
	type args struct {
		ctx context.Context
		req models.ExternalTransactionRequest
	}
	tests := []struct {
		name    string
		args    args
		want    models.BalanceResponse
		wantErr bool
		mockFn  func(args args)
	}{
		{
			name: "success debit transaction",
			args: args{
				ctx: context.Background(),
				req: models.ExternalTransactionRequest{
					Amount:          50000,
					Reference:       "reference",
					TransactionType: "DEBIT",
					WalletID:        1,
				},
			},
			want: models.BalanceResponse{
				Balance: 150000,
			},
			wantErr: false,
			mockFn: func(args args) {
				mockRepo.EXPECT().GetWalletTransactionByReference(args.ctx, args.req.Reference).Return(models.WalletTransaction{}, nil)

				wallet := models.Wallet{
					ID:        1,
					UserID:    1,
					Balance:   200000,
					CreatedAt: now,
					UpdatedAt: now,
				}
				mockRepo.EXPECT().UpdateBalanceByID(args.ctx, args.req.WalletID, (-args.req.Amount)).Return(wallet, nil)

				walletTrx := &models.WalletTransaction{
					WalletID:              wallet.ID,
					Amount:                args.req.Amount,
					Reference:             args.req.Reference,
					WalletTransactionType: args.req.TransactionType,
				}
				mockRepo.EXPECT().CreateWalletTrx(args.ctx, walletTrx).Return(nil)

			},
		},
		{
			name: "success credit transaction",
			args: args{
				ctx: context.Background(),
				req: models.ExternalTransactionRequest{
					Amount:          50000,
					Reference:       "reference",
					TransactionType: "CREDIT",
					WalletID:        1,
				},
			},
			want: models.BalanceResponse{
				Balance: 250000,
			},
			wantErr: false,
			mockFn: func(args args) {
				mockRepo.EXPECT().GetWalletTransactionByReference(args.ctx, args.req.Reference).Return(models.WalletTransaction{}, nil)

				wallet := models.Wallet{
					ID:        1,
					UserID:    1,
					Balance:   200000,
					CreatedAt: now,
					UpdatedAt: now,
				}
				mockRepo.EXPECT().UpdateBalanceByID(args.ctx, args.req.WalletID, args.req.Amount).Return(wallet, nil)

				walletTrx := &models.WalletTransaction{
					WalletID:              wallet.ID,
					Amount:                args.req.Amount,
					Reference:             args.req.Reference,
					WalletTransactionType: args.req.TransactionType,
				}
				mockRepo.EXPECT().CreateWalletTrx(args.ctx, walletTrx).Return(nil)

			},
		},
		{
			name: "error, get duplicate wallet transaction",
			args: args{
				ctx: context.Background(),
				req: models.ExternalTransactionRequest{
					Amount:          50000,
					Reference:       "reference",
					TransactionType: "CREDIT",
					WalletID:        1,
				},
			},
			want:    models.BalanceResponse{},
			wantErr: true,
			mockFn: func(args args) {
				mockRepo.EXPECT().GetWalletTransactionByReference(args.ctx, args.req.Reference).Return(models.WalletTransaction{
					ID:                    1,
					WalletID:              1,
					Amount:                50000,
					WalletTransactionType: args.req.TransactionType,
					Reference:             args.req.Reference,
					CreatedAt:             now,
					UpdatedAt:             now,
				}, assert.AnError)

			},
		},
		{
			name: "error update balance",
			args: args{
				ctx: context.Background(),
				req: models.ExternalTransactionRequest{
					Amount:          50000,
					Reference:       "reference",
					TransactionType: "CREDIT",
					WalletID:        1,
				},
			},
			want:    models.BalanceResponse{},
			wantErr: true,
			mockFn: func(args args) {
				mockRepo.EXPECT().GetWalletTransactionByReference(args.ctx, args.req.Reference).Return(models.WalletTransaction{}, nil)

				mockRepo.EXPECT().UpdateBalanceByID(args.ctx, args.req.WalletID, args.req.Amount).Return(models.Wallet{}, assert.AnError)

			},
		},
		{
			name: "error create wallet transaction",
			args: args{
				ctx: context.Background(),
				req: models.ExternalTransactionRequest{
					Amount:          50000,
					Reference:       "reference",
					TransactionType: "DEBIT",
					WalletID:        1,
				},
			},
			want:    models.BalanceResponse{},
			wantErr: true,
			mockFn: func(args args) {
				mockRepo.EXPECT().GetWalletTransactionByReference(args.ctx, args.req.Reference).Return(models.WalletTransaction{}, nil)

				wallet := models.Wallet{
					ID:        1,
					UserID:    1,
					Balance:   200000,
					CreatedAt: now,
					UpdatedAt: now,
				}
				mockRepo.EXPECT().UpdateBalanceByID(args.ctx, args.req.WalletID, (-args.req.Amount)).Return(wallet, nil)

				walletTrx := &models.WalletTransaction{
					WalletID:              wallet.ID,
					Amount:                args.req.Amount,
					Reference:             args.req.Reference,
					WalletTransactionType: args.req.TransactionType,
				}
				mockRepo.EXPECT().CreateWalletTrx(args.ctx, walletTrx).Return(assert.AnError)

			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.mockFn(tt.args)
			s := &WalletService{
				WalletRepo: mockRepo,
			}
			got, err := s.ExternalTransaction(tt.args.ctx, tt.args.req)
			if (err != nil) != tt.wantErr {
				t.Errorf("WalletService.ExternalTransaction() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("WalletService.ExternalTransaction() = %v, want %v", got, tt.want)
			}
		})
	}
}
