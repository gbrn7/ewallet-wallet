package wallet

import (
	"bytes"
	"context"
	"encoding/json"
	"ewallet-wallet/constants"
	"ewallet-wallet/helpers"
	"ewallet-wallet/internal/models"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/assert"
)

func TestHandler_Create(t *testing.T) {
	ctrlMock := gomock.NewController(t)
	defer ctrlMock.Finish()

	mockSvc := NewMockService(ctrlMock)
	mockMdw := NewMockMiddleware(ctrlMock)
	mockExt := NewMockExternal(ctrlMock)

	ctx := context.Background()
	now := time.Now()

	tests := []struct {
		name               string
		mockFn             func()
		expectedStatusCode int
		expectedBody       helpers.Response
		wantErr            bool
	}{
		{
			name: "success",
			mockFn: func() {
				wallet := &models.Wallet{
					UserID:  1,
					Balance: 200000,
				}
				mockSvc.EXPECT().Create(ctx, wallet).DoAndReturn(func(ctx context.Context, wallet *models.Wallet) error {
					*wallet = models.Wallet{
						ID:        1,
						UserID:    1,
						Balance:   200000,
						CreatedAt: now,
						UpdatedAt: now,
					}

					return nil
				})
			},
			expectedStatusCode: http.StatusOK,
			wantErr:            false,
			expectedBody: helpers.Response{
				Message: constants.SuccessMessage,
				Data: map[string]interface{}{
					"id":        float64(1),
					"user_id":   float64(1),
					"balance":   float64(200000),
					"CreatedAt": now.Format(time.RFC3339Nano),
					"UpdatedAt": now.Format(time.RFC3339Nano),
				},
			},
		},
		{
			name: "error",
			mockFn: func() {
				wallet := &models.Wallet{
					UserID:  1,
					Balance: 200000,
				}
				mockSvc.EXPECT().Create(ctx, wallet).Return(assert.AnError)
			},
			expectedStatusCode: http.StatusBadRequest,
			wantErr:            true,
			expectedBody: helpers.Response{
				Message: constants.ErrServerError,
				Data: map[string]interface{}{
					"user_id": float64(1),
					"balance": float64(200000),
				},
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.mockFn()
			api := gin.New()
			h := &Handler{
				Engine:     api,
				Service:    mockSvc,
				External:   mockExt,
				Middleware: mockMdw,
			}
			h.RegisterRoute()
			w := httptest.NewRecorder()

			endPoint := "/wallet/v1/"

			model := models.Wallet{
				UserID:  1,
				Balance: 200000,
			}

			val, err := json.Marshal(model)
			assert.NoError(t, err)

			body := bytes.NewReader(val)
			req, err := http.NewRequest(http.MethodPost, endPoint, body)
			assert.NoError(t, err)
			h.ServeHTTP(w, req)

			assert.Equal(t, tt.expectedStatusCode, w.Code)
			if !tt.wantErr {
				res := w.Result()
				defer res.Body.Close()

				response := helpers.Response{}
				err = json.Unmarshal(w.Body.Bytes(), &response)
				assert.NoError(t, err)

				assert.Equal(t, tt.expectedBody, response)
			}
		})
	}
}

func TestHandler_CreditBalance(t *testing.T) {
	ctrlMock := gomock.NewController(t)
	defer ctrlMock.Finish()

	mockSvc := NewMockService(ctrlMock)
	mockMdw := NewMockMiddleware(ctrlMock)
	mockExt := NewMockExternal(ctrlMock)

	reference := "reference"

	tests := []struct {
		name               string
		mockFn             func()
		expectedStatusCode int
		expectedBody       helpers.Response
		wantErr            bool
	}{
		{
			name: "success",
			mockFn: func() {
				mockMdw.EXPECT().MiddlewareValidateToken(gomock.Any()).Do(func(c *gin.Context) {
					tokenData := models.TokenData{
						UserID:   1,
						Username: "username",
						Fullname: "fullname",
						Email:    "email",
					}
					c.Set("token", tokenData)
				})

				transactionReq := models.TransactionRequest{
					Reference: reference,
					Amount:    100000,
				}
				mockSvc.EXPECT().CreditBalance(gomock.Any(), uint64(1), transactionReq).Return(models.BalanceResponse{
					Balance: 300000,
				}, nil)
			},
			expectedStatusCode: http.StatusOK,
			expectedBody: helpers.Response{
				Message: constants.SuccessMessage,
				Data: map[string]interface{}{
					"balance": float64(300000),
				},
			},
			wantErr: false,
		},
		{
			name: "error ",
			mockFn: func() {
				mockMdw.EXPECT().MiddlewareValidateToken(gomock.Any()).Do(func(c *gin.Context) {
					tokenData := models.TokenData{
						UserID:   1,
						Username: "username",
						Fullname: "fullname",
						Email:    "email",
					}
					c.Set("token", tokenData)
				})

				transactionReq := models.TransactionRequest{
					Reference: reference,
					Amount:    100000,
				}
				mockSvc.EXPECT().CreditBalance(gomock.Any(), uint64(1), transactionReq).Return(models.BalanceResponse{}, assert.AnError)
			},
			expectedStatusCode: http.StatusInternalServerError,
			expectedBody: helpers.Response{
				Message: constants.ErrServerError,
			},
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.mockFn()
			api := gin.New()

			h := &Handler{
				Engine:     api,
				Service:    mockSvc,
				External:   mockExt,
				Middleware: mockMdw,
			}
			h.RegisterRoute()
			w := httptest.NewRecorder()

			endPoint := "/wallet/v1/balance/credit"

			model := models.TransactionRequest{
				Reference: reference,
				Amount:    100000,
			}
			val, err := json.Marshal(model)
			assert.NoError(t, err)

			body := bytes.NewReader(val)
			req, err := http.NewRequest(http.MethodPut, endPoint, body)
			assert.NoError(t, err)
			req.Header.Set("Authorization", "authorization")

			h.ServeHTTP(w, req)

			assert.Equal(t, tt.expectedStatusCode, w.Code)
			if !tt.wantErr {
				res := w.Result()
				defer res.Body.Close()

				response := helpers.Response{}
				err = json.Unmarshal(w.Body.Bytes(), &response)
				assert.NoError(t, err)

				assert.Equal(t, tt.expectedBody, response)
			}
		})
	}
}

func TestHandler_DebitBalance(t *testing.T) {
	ctrlMock := gomock.NewController(t)
	defer ctrlMock.Finish()

	mockSvc := NewMockService(ctrlMock)
	mockMdw := NewMockMiddleware(ctrlMock)
	mockExt := NewMockExternal(ctrlMock)

	reference := "reference"

	tests := []struct {
		name               string
		mockFn             func()
		expectedStatusCode int
		expectedBody       helpers.Response
		wantErr            bool
	}{
		{
			name: "success",
			mockFn: func() {
				mockMdw.EXPECT().MiddlewareValidateToken(gomock.Any()).Do(func(c *gin.Context) {
					tokenData := models.TokenData{
						UserID:   1,
						Username: "username",
						Fullname: "fullname",
						Email:    "email",
					}
					c.Set("token", tokenData)
				})

				transactionReq := models.TransactionRequest{
					Reference: reference,
					Amount:    100000,
				}
				mockSvc.EXPECT().DebitBalance(gomock.Any(), uint64(1), transactionReq).Return(models.BalanceResponse{
					Balance: 100000,
				}, nil)
			},
			expectedStatusCode: http.StatusOK,
			expectedBody: helpers.Response{
				Message: constants.SuccessMessage,
				Data: map[string]interface{}{
					"balance": float64(100000),
				},
			},
			wantErr: false,
		},
		{
			name: "error ",
			mockFn: func() {
				mockMdw.EXPECT().MiddlewareValidateToken(gomock.Any()).Do(func(c *gin.Context) {
					tokenData := models.TokenData{
						UserID:   1,
						Username: "username",
						Fullname: "fullname",
						Email:    "email",
					}
					c.Set("token", tokenData)
				})

				transactionReq := models.TransactionRequest{
					Reference: reference,
					Amount:    100000,
				}
				mockSvc.EXPECT().DebitBalance(gomock.Any(), uint64(1), transactionReq).Return(models.BalanceResponse{}, assert.AnError)
			},
			expectedStatusCode: http.StatusInternalServerError,
			expectedBody: helpers.Response{
				Message: constants.ErrServerError,
			},
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.mockFn()
			api := gin.New()

			h := &Handler{
				Engine:     api,
				Service:    mockSvc,
				External:   mockExt,
				Middleware: mockMdw,
			}
			h.RegisterRoute()
			w := httptest.NewRecorder()

			endPoint := "/wallet/v1/balance/debit"

			model := models.TransactionRequest{
				Reference: reference,
				Amount:    100000,
			}
			val, err := json.Marshal(model)
			assert.NoError(t, err)

			body := bytes.NewReader(val)
			req, err := http.NewRequest(http.MethodPut, endPoint, body)
			assert.NoError(t, err)
			req.Header.Set("Authorization", "authorization")

			h.ServeHTTP(w, req)

			assert.Equal(t, tt.expectedStatusCode, w.Code)
			if !tt.wantErr {
				res := w.Result()
				defer res.Body.Close()

				response := helpers.Response{}
				err = json.Unmarshal(w.Body.Bytes(), &response)
				assert.NoError(t, err)

				assert.Equal(t, tt.expectedBody, response)
			}
		})
	}
}

func TestHandler_GetBalance(t *testing.T) {
	ctrlMock := gomock.NewController(t)
	defer ctrlMock.Finish()

	mockSvc := NewMockService(ctrlMock)
	mockMdw := NewMockMiddleware(ctrlMock)
	mockExt := NewMockExternal(ctrlMock)

	tokenData := models.TokenData{
		UserID:   1,
		Username: "username",
		Fullname: "fullname",
		Email:    "email",
	}

	tests := []struct {
		name               string
		mockFn             func()
		expectedStatusCode int
		expectedBody       helpers.Response
		wantErr            bool
	}{
		{
			name: "success",
			mockFn: func() {
				mockMdw.EXPECT().MiddlewareValidateToken(gomock.Any()).Do(func(c *gin.Context) {
					c.Set("token", tokenData)
				})

				mockSvc.EXPECT().GetBalance(gomock.Any(), tokenData.UserID).Return(models.BalanceResponse{
					Balance: 200000,
				}, nil)
			},
			expectedStatusCode: http.StatusOK,
			expectedBody: helpers.Response{
				Message: constants.SuccessMessage,
				Data: map[string]interface{}{
					"balance": float64(200000),
				},
			},
			wantErr: false,
		},
		{
			name: "error",
			mockFn: func() {
				mockMdw.EXPECT().MiddlewareValidateToken(gomock.Any()).Do(func(c *gin.Context) {
					c.Set("token", tokenData)
				})

				mockSvc.EXPECT().GetBalance(gomock.Any(), tokenData.UserID).Return(models.BalanceResponse{}, assert.AnError)
			},
			expectedStatusCode: http.StatusInternalServerError,
			expectedBody: helpers.Response{
				Message: constants.SuccessMessage,
			},
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.mockFn()
			api := gin.New()
			h := &Handler{
				Engine:     api,
				Service:    mockSvc,
				External:   mockExt,
				Middleware: mockMdw,
			}
			h.RegisterRoute()
			w := httptest.NewRecorder()

			endPoint := "/wallet/v1/balance"
			req, err := http.NewRequest(http.MethodGet, endPoint, nil)
			assert.NoError(t, err)

			h.ServeHTTP(w, req)
			assert.Equal(t, tt.expectedStatusCode, w.Code)
			if !tt.wantErr {
				res := w.Result()
				defer res.Body.Close()

				response := helpers.Response{}
				err = json.Unmarshal(w.Body.Bytes(), &response)
				assert.NoError(t, err)

				assert.Equal(t, tt.expectedBody, response)
			}
		})
	}
}

func TestHandler_GetWalletHistory(t *testing.T) {
	ctrlMock := gomock.NewController(t)
	defer ctrlMock.Finish()

	mockSvc := NewMockService(ctrlMock)
	mockMdw := NewMockMiddleware(ctrlMock)
	mockExt := NewMockExternal(ctrlMock)

	reference := "REFERENCE"
	now := time.Now()
	tokenData := models.TokenData{
		UserID:   1,
		Username: "username",
		Fullname: "fullname",
		Email:    "email",
	}
	tests := []struct {
		name               string
		mockFn             func()
		expectedStatusCode int
		expectedBody       helpers.Response
		wantErr            bool
	}{
		{
			name: "success",
			mockFn: func() {
				mockMdw.EXPECT().MiddlewareValidateToken(gomock.Any()).Do(func(c *gin.Context) {
					c.Set("token", tokenData)
				})

				mockSvc.EXPECT().GetWalletHistory(gomock.Any(), tokenData.UserID, models.WalletHistoryParam{
					Page:                  2,
					Limit:                 2,
					WalletTransactionType: "DEBIT",
				}).Return([]models.WalletTransaction{
					{
						ID:                    1,
						WalletID:              1,
						Amount:                200000,
						WalletTransactionType: "DEBIT",
						Reference:             reference,
						CreatedAt:             now,
						UpdatedAt:             now,
					},
					{
						ID:                    2,
						WalletID:              1,
						Amount:                300000,
						WalletTransactionType: "DEBIT",
						Reference:             reference,
						CreatedAt:             now,
						UpdatedAt:             now,
					},
				}, nil)
			},
			expectedStatusCode: http.StatusOK,
			wantErr:            false,
			expectedBody: helpers.Response{
				Message: constants.SuccessMessage,
				Data: []interface{}{
					map[string]interface{}{
						"id":                      float64(1),
						"wallet_id":               float64(1),
						"amount":                  float64(200000),
						"wallet_transaction_type": "DEBIT",
						"reference":               reference,
						"created_at":              now.Format(time.RFC3339Nano),
						"updated_at":              now.Format(time.RFC3339Nano),
					},
					map[string]interface{}{
						"id":                      float64(2),
						"wallet_id":               float64(1),
						"amount":                  float64(300000),
						"wallet_transaction_type": "DEBIT",
						"reference":               reference,
						"created_at":              now.Format(time.RFC3339Nano),
						"updated_at":              now.Format(time.RFC3339Nano),
					},
				},
			},
		},
		{
			name: "error",
			mockFn: func() {
				mockMdw.EXPECT().MiddlewareValidateToken(gomock.Any()).Do(func(c *gin.Context) {
					c.Set("token", tokenData)
				})

				mockSvc.EXPECT().GetWalletHistory(gomock.Any(), tokenData.UserID, models.WalletHistoryParam{
					Page:                  2,
					Limit:                 2,
					WalletTransactionType: "DEBIT",
				}).Return([]models.WalletTransaction{}, assert.AnError)
			},
			expectedStatusCode: http.StatusInternalServerError,
			wantErr:            true,
			expectedBody: helpers.Response{
				Message: constants.ErrFailedBadRequest,
				Data: []interface{}{
					map[string]interface{}{
						"id":                      float64(1),
						"wallet_id":               float64(1),
						"amount":                  float64(200000),
						"wallet_transaction_type": "DEBIT",
						"reference":               reference,
						"created_at":              now.Format(time.RFC3339Nano),
						"updated_at":              now.Format(time.RFC3339Nano),
					},
					map[string]interface{}{
						"id":                      float64(2),
						"wallet_id":               float64(1),
						"amount":                  float64(300000),
						"wallet_transaction_type": "DEBIT",
						"reference":               reference,
						"created_at":              now.Format(time.RFC3339Nano),
						"updated_at":              now.Format(time.RFC3339Nano),
					},
				},
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.mockFn()
			api := gin.New()
			h := &Handler{
				Engine:     api,
				Service:    mockSvc,
				External:   mockExt,
				Middleware: mockMdw,
			}
			h.RegisterRoute()
			w := httptest.NewRecorder()

			endPoint := "/wallet/v1/history?page=2&limit=2&wallet_transaction_type=DEBIT"
			req, err := http.NewRequest(http.MethodGet, endPoint, nil)
			assert.NoError(t, err)

			h.ServeHTTP(w, req)
			assert.Equal(t, tt.expectedStatusCode, w.Code)
			if !tt.wantErr {
				res := w.Result()
				defer res.Body.Close()

				response := helpers.Response{}
				err = json.Unmarshal(w.Body.Bytes(), &response)
				assert.NoError(t, err)

				assert.Equal(t, tt.expectedBody, response)
			}
		})
	}
}

func TestHandler_CreateWalletLink(t *testing.T) {
	ctrlMock := gomock.NewController(t)
	defer ctrlMock.Finish()

	mockSvc := NewMockService(ctrlMock)
	mockMdw := NewMockMiddleware(ctrlMock)
	mockExt := NewMockExternal(ctrlMock)

	clientSource := "fastcampus_ecommerce"
	otp := "121212"
	clientID := "fastcampus_ecommerce"

	tests := []struct {
		name               string
		mockFn             func()
		expectedStatusCode int
		expectedBody       helpers.Response
		wantErr            bool
	}{
		{
			name: "success",
			mockFn: func() {
				mockMdw.EXPECT().MiddlewareSignatureValidation(gomock.Any()).Do(func(c *gin.Context) {
					c.Set("client_id", clientID)
					c.Next()
				})

				mockSvc.EXPECT().CreateWalletLink(gomock.Any(), clientSource, &models.WalletLink{
					WalletID:     1,
					ClientSource: clientSource,
					OTP:          otp,
					Status:       "pending",
				}).Return(&models.WalletStructOTP{
					OTP: otp,
				}, nil)
			},
			expectedStatusCode: http.StatusOK,
			expectedBody: helpers.Response{
				Message: constants.SuccessMessage,
				Data: map[string]interface{}{
					"otp": otp,
				},
			},
			wantErr: false,
		},
		{
			name: "error",
			mockFn: func() {
				mockMdw.EXPECT().MiddlewareSignatureValidation(gomock.Any()).Do(func(c *gin.Context) {
					c.Set("client_id", clientID)
					c.Next()
				})

				mockSvc.EXPECT().CreateWalletLink(gomock.Any(), clientSource, &models.WalletLink{
					WalletID:     1,
					ClientSource: clientSource,
					OTP:          otp,
					Status:       "pending",
				}).Return(&models.WalletStructOTP{}, assert.AnError)
			},
			expectedStatusCode: http.StatusInternalServerError,
			expectedBody: helpers.Response{
				Message: constants.ErrServerError,
			},
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.mockFn()
			api := gin.New()
			h := &Handler{
				Engine:     api,
				Service:    mockSvc,
				External:   mockExt,
				Middleware: mockMdw,
			}
			h.RegisterRoute()
			w := httptest.NewRecorder()

			endPoint := "/wallet/v1/ex/link"
			model := models.WalletLink{
				WalletID:     1,
				ClientSource: clientSource,
				OTP:          otp,
				Status:       "pending",
			}
			val, err := json.Marshal(model)
			assert.NoError(t, err)

			body := bytes.NewReader(val)
			req, err := http.NewRequest(http.MethodPost, endPoint, body)
			assert.NoError(t, err)

			h.ServeHTTP(w, req)

			assert.Equal(t, tt.expectedStatusCode, w.Code)
			if !tt.wantErr {
				res := w.Result()
				defer res.Body.Close()

				response := helpers.Response{}
				err = json.Unmarshal(w.Body.Bytes(), &response)
				assert.NoError(t, err)

				assert.Equal(t, tt.expectedBody, response)
			}
		})
	}
}

func TestHandler_WalletLinkConfirmation(t *testing.T) {
	ctrlMock := gomock.NewController(t)
	defer ctrlMock.Finish()

	mockSvc := NewMockService(ctrlMock)
	mockMdw := NewMockMiddleware(ctrlMock)
	mockExt := NewMockExternal(ctrlMock)

	clientSource := "fastcampus_ecommerce"
	otp := "121212"
	clientID := "fastcampus_ecommerce"

	tests := []struct {
		name               string
		mockFn             func()
		expectedStatusCode int
		expectedBody       helpers.Response
		wantErr            bool
	}{
		{
			name: "success",
			mockFn: func() {
				mockMdw.EXPECT().MiddlewareSignatureValidation(gomock.Any()).Do(func(c *gin.Context) {
					c.Set("client_id", clientID)
					c.Next()
				})

				mockSvc.EXPECT().WalletLinkConfirmation(gomock.Any(), 1, clientSource, otp).Return(nil)
			},
			expectedStatusCode: http.StatusOK,
			expectedBody: helpers.Response{
				Message: constants.SuccessMessage,
			},
		},
		{
			name: "error",
			mockFn: func() {
				mockMdw.EXPECT().MiddlewareSignatureValidation(gomock.Any()).Do(func(c *gin.Context) {
					c.Set("client_id", clientID)
					c.Next()
				})

				mockSvc.EXPECT().WalletLinkConfirmation(gomock.Any(), 1, clientSource, otp).Return(assert.AnError)
			},
			expectedStatusCode: http.StatusInternalServerError,
			expectedBody: helpers.Response{
				Message: constants.ErrServerError,
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.mockFn()
			api := gin.New()
			h := &Handler{
				Engine:     api,
				Service:    mockSvc,
				External:   mockExt,
				Middleware: mockMdw,
			}
			h.RegisterRoute()
			w := httptest.NewRecorder()

			endPoint := "/wallet/v1/ex/link/1/confirmation"
			model := models.WalletLink{
				OTP: otp,
			}

			val, err := json.Marshal(model)
			assert.NoError(t, err)

			body := bytes.NewReader(val)
			req, err := http.NewRequest(http.MethodPut, endPoint, body)
			assert.NoError(t, err)
			h.ServeHTTP(w, req)

			assert.Equal(t, tt.expectedStatusCode, w.Code)
			if !tt.wantErr {
				res := w.Result()
				defer res.Body.Close()

				response := helpers.Response{}
				err = json.Unmarshal(w.Body.Bytes(), &response)
				assert.NoError(t, err)

				assert.Equal(t, tt.expectedBody, response)
			}
		})
	}
}

func TestHandler_WalletUnlink(t *testing.T) {
	ctrlMock := gomock.NewController(t)
	defer ctrlMock.Finish()

	mockSvc := NewMockService(ctrlMock)
	mockMdw := NewMockMiddleware(ctrlMock)
	mockExt := NewMockExternal(ctrlMock)

	clientSource := "fastcampus_ecommerce"
	clientID := "fastcampus_ecommerce"

	tests := []struct {
		name               string
		mockFn             func()
		expectedStatusCode int
		expectedBody       helpers.Response
		wantErr            bool
	}{
		{
			name: "success",
			mockFn: func() {
				mockMdw.EXPECT().MiddlewareSignatureValidation(gomock.Any()).Do(func(c *gin.Context) {
					c.Set("client_id", clientID)
					c.Next()
				})

				mockSvc.EXPECT().WalletUnlink(gomock.Any(), 1, clientSource).Return(nil)
			},
			expectedStatusCode: http.StatusOK,
			expectedBody: helpers.Response{
				Message: constants.SuccessMessage,
			},
			wantErr: false,
		},
		{
			name: "error",
			mockFn: func() {
				mockMdw.EXPECT().MiddlewareSignatureValidation(gomock.Any()).Do(func(c *gin.Context) {
					c.Set("client_id", clientID)
					c.Next()
				})

				mockSvc.EXPECT().WalletUnlink(gomock.Any(), 1, clientSource).Return(assert.AnError)
			},
			expectedStatusCode: http.StatusInternalServerError,
			expectedBody: helpers.Response{
				Message: constants.ErrServerError,
			},
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.mockFn()
			api := gin.New()
			h := &Handler{
				Engine:     api,
				Service:    mockSvc,
				External:   mockExt,
				Middleware: mockMdw,
			}
			h.RegisterRoute()
			w := httptest.NewRecorder()

			endPoint := "/wallet/v1/ex/1/unlink"
			req, err := http.NewRequest(http.MethodDelete, endPoint, nil)
			assert.NoError(t, err)
			h.ServeHTTP(w, req)

			assert.Equal(t, tt.expectedStatusCode, w.Code)
			if !tt.wantErr {
				res := w.Result()
				defer res.Body.Close()

				response := helpers.Response{}
				err = json.Unmarshal(w.Body.Bytes(), &response)
				assert.NoError(t, err)

				assert.Equal(t, tt.expectedBody, response)
			}
		})
	}
}

func TestHandler_ExGetBalance(t *testing.T) {
	ctrlMock := gomock.NewController(t)
	defer ctrlMock.Finish()

	mockSvc := NewMockService(ctrlMock)
	mockMdw := NewMockMiddleware(ctrlMock)
	mockExt := NewMockExternal(ctrlMock)

	clientID := "fastcampus_ecommerce"

	tests := []struct {
		name               string
		mockFn             func()
		expectedStatusCode int
		expectedBody       helpers.Response
		wantErr            bool
	}{
		{
			name: "success",
			mockFn: func() {
				mockMdw.EXPECT().MiddlewareSignatureValidation(gomock.Any()).Do(func(c *gin.Context) {
					c.Set("client_id", clientID)
					c.Next()
				})

				mockSvc.EXPECT().ExGetBalance(gomock.Any(), 1).Return(models.BalanceResponse{
					Balance: 200000,
				}, nil)
			},
			expectedStatusCode: http.StatusOK,
			expectedBody: helpers.Response{
				Message: constants.SuccessMessage,
				Data: map[string]interface{}{
					"balance": float64(200000),
				},
			},
			wantErr: false,
		},
		{
			name: "error",
			mockFn: func() {
				mockMdw.EXPECT().MiddlewareSignatureValidation(gomock.Any()).Do(func(c *gin.Context) {
					c.Set("client_id", clientID)
					c.Next()
				})

				mockSvc.EXPECT().ExGetBalance(gomock.Any(), 1).Return(models.BalanceResponse{}, assert.AnError)
			},
			expectedStatusCode: http.StatusInternalServerError,
			expectedBody: helpers.Response{
				Message: constants.ErrServerError,
				Data: map[string]interface{}{
					"balance": float64(200000),
				},
			},
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.mockFn()
			api := gin.New()
			h := &Handler{
				Engine:     api,
				Service:    mockSvc,
				External:   mockExt,
				Middleware: mockMdw,
			}
			h.RegisterRoute()
			w := httptest.NewRecorder()

			endPoint := "/wallet/v1/ex/1/balance"
			req, err := http.NewRequest(http.MethodGet, endPoint, nil)
			assert.NoError(t, err)

			h.ServeHTTP(w, req)
			assert.Equal(t, tt.expectedStatusCode, w.Code)
			if !tt.wantErr {
				res := w.Result()
				defer res.Body.Close()

				response := helpers.Response{}
				err = json.Unmarshal(w.Body.Bytes(), &response)
				assert.NoError(t, err)

				assert.Equal(t, tt.expectedBody, response)
			}
		})
	}
}

func TestHandler_ExternalTransaction(t *testing.T) {
	ctrlMock := gomock.NewController(t)
	defer ctrlMock.Finish()

	mockSvc := NewMockService(ctrlMock)
	mockMdw := NewMockMiddleware(ctrlMock)
	mockExt := NewMockExternal(ctrlMock)

	clientID := "fastcampus_ecommerce"
	reference := "reference"
	transactionType := "DEBIT"

	tests := []struct {
		name               string
		mockFn             func()
		expectedStatusCode int
		expectedBody       helpers.Response
		wantErr            bool
	}{
		{
			name: "success",
			mockFn: func() {
				mockMdw.EXPECT().MiddlewareSignatureValidation(gomock.Any()).Do(func(c *gin.Context) {
					c.Set("client_id", clientID)
					c.Next()
				})

				mockSvc.EXPECT().ExternalTransaction(gomock.Any(), models.ExternalTransactionRequest{
					Amount:          100000,
					Reference:       reference,
					TransactionType: transactionType,
					WalletID:        1,
				}).Return(models.BalanceResponse{
					Balance: 100000,
				}, nil)
			},
			expectedStatusCode: http.StatusOK,
			expectedBody: helpers.Response{
				Message: constants.SuccessMessage,
				Data: map[string]interface{}{
					"balance": float64(100000),
				},
			},
			wantErr: false,
		},
		{
			name: "error",
			mockFn: func() {
				mockMdw.EXPECT().MiddlewareSignatureValidation(gomock.Any()).Do(func(c *gin.Context) {
					c.Set("client_id", clientID)
					c.Next()
				})

				mockSvc.EXPECT().ExternalTransaction(gomock.Any(), models.ExternalTransactionRequest{
					Amount:          100000,
					Reference:       reference,
					TransactionType: transactionType,
					WalletID:        1,
				}).Return(models.BalanceResponse{}, assert.AnError)
			},
			expectedStatusCode: http.StatusInternalServerError,
			expectedBody: helpers.Response{
				Message: constants.ErrServerError,
			},
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.mockFn()
			api := gin.New()
			h := &Handler{
				Engine:     api,
				Service:    mockSvc,
				External:   mockExt,
				Middleware: mockMdw,
			}
			h.RegisterRoute()
			w := httptest.NewRecorder()

			endPoint := "/wallet/v1/ex/transaction"

			model := models.ExternalTransactionRequest{
				Amount:          100000,
				Reference:       reference,
				TransactionType: transactionType,
				WalletID:        1,
			}
			val, err := json.Marshal(model)
			assert.NoError(t, err)

			body := bytes.NewReader(val)
			req, err := http.NewRequest(http.MethodPost, endPoint, body)
			assert.NoError(t, err)

			h.ServeHTTP(w, req)
			assert.Equal(t, tt.expectedStatusCode, w.Code)
			if !tt.wantErr {
				res := w.Result()
				defer res.Body.Close()

				response := helpers.Response{}
				err = json.Unmarshal(w.Body.Bytes(), &response)
				assert.NoError(t, err)

				assert.Equal(t, tt.expectedBody, response)
			}
		})
	}
}
