package middleware

import (
	"bytes"
	"encoding/json"
	"ewallet-wallet/generate_signature"
	"ewallet-wallet/internal/models"
	"io"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/assert"
)

func TestExternalDependency_MiddlewareValidateToken(t *testing.T) {
	ctrlMock := gomock.NewController(t)
	defer ctrlMock.Finish()

	mockExt := NewMockExternal(ctrlMock)

	auth := "Authorization"

	tests := []struct {
		name               string
		wantErr            bool
		mockFn             func()
		expectedStatusCode int
	}{
		{
			name:    "success",
			wantErr: false,
			mockFn: func() {
				mockExt.EXPECT().ValidateToken(gomock.Any(), auth).Return(models.TokenData{
					UserID:   1,
					Username: "username",
					Fullname: "fullname",
					Email:    "email@gmail.com",
				}, nil)
			},
			expectedStatusCode: http.StatusOK,
		},
		{
			name:    "error",
			wantErr: true,
			mockFn: func() {
				mockExt.EXPECT().ValidateToken(gomock.Any(), auth).Return(models.TokenData{}, assert.AnError)
			},
			expectedStatusCode: http.StatusUnauthorized,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.mockFn()
			api := gin.New()

			d := &ExternalDependency{
				External: mockExt,
			}

			w := httptest.NewRecorder()
			endPoint := "/validate-token"
			api.GET(endPoint, d.MiddlewareValidateToken)

			req, err := http.NewRequest(http.MethodGet, endPoint, nil)
			assert.NoError(t, err)
			req.Header.Set("Authorization", auth)

			api.ServeHTTP(w, req)
			assert.Equal(t, tt.expectedStatusCode, w.Code)

		})
	}
}

func TestExternalDependency_MiddlewareSignatureValidation(t *testing.T) {
	ctrlMock := gomock.NewController(t)
	defer ctrlMock.Finish()

	mockExt := NewMockExternal(ctrlMock)

	clientID := "fastcampus_ecommerce"
	endPoint := "/signature-validation"

	now := time.Now()

	model := models.WalletLink{
		WalletID:     1,
		ClientSource: "clientSource",
		OTP:          "otp",
		Status:       "pending",
	}
	val, err := json.Marshal(model)
	assert.NoError(t, err)

	tests := []struct {
		name               string
		wantErr            bool
		expectedStatusCode int
		signature          string
		method             string
	}{
		{
			name:               "success with get method",
			wantErr:            false,
			expectedStatusCode: http.StatusOK,
			signature:          generate_signature.GenerateSignature(endPoint, clientID, now, http.MethodGet, ``),
			method:             http.MethodGet,
		},
		{
			name:               "error with get method",
			wantErr:            true,
			expectedStatusCode: http.StatusUnauthorized,
			signature:          "signature",
			method:             http.MethodGet,
		},
		{
			name:               "success with non get method",
			wantErr:            false,
			expectedStatusCode: http.StatusOK,
			signature:          generate_signature.GenerateSignature(endPoint, clientID, now, http.MethodPost, string(val)),
			method:             http.MethodPost,
		},
		{
			name:               "error with non get method",
			wantErr:            true,
			expectedStatusCode: http.StatusUnauthorized,
			signature:          "signature",
			method:             http.MethodPost,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			api := gin.New()

			d := &ExternalDependency{
				External: mockExt,
			}

			w := httptest.NewRecorder()

			api.GET(endPoint, d.MiddlewareSignatureValidation)
			api.POST(endPoint, d.MiddlewareSignatureValidation)

			var body io.Reader

			if tt.method != http.MethodGet {

				body = bytes.NewReader(val)
			}

			req, err := http.NewRequest(tt.method, endPoint, body)
			assert.NoError(t, err)

			req.Header.Set("Client-id", clientID)
			req.Header.Set("Timestamp", now.Format(time.RFC3339))
			req.Header.Set("Signature", tt.signature)

			api.ServeHTTP(w, req)

			assert.Equal(t, tt.expectedStatusCode, w.Code)
		})
	}
}
