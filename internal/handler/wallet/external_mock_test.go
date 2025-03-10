// Code generated by MockGen. DO NOT EDIT.
// Source: external.go

// Package wallet is a generated GoMock package.
package wallet

import (
	context "context"
	models "ewallet-wallet/internal/models"
	reflect "reflect"

	gomock "github.com/golang/mock/gomock"
)

// MockExternal is a mock of External interface.
type MockExternal struct {
	ctrl     *gomock.Controller
	recorder *MockExternalMockRecorder
}

// MockExternalMockRecorder is the mock recorder for MockExternal.
type MockExternalMockRecorder struct {
	mock *MockExternal
}

// NewMockExternal creates a new mock instance.
func NewMockExternal(ctrl *gomock.Controller) *MockExternal {
	mock := &MockExternal{ctrl: ctrl}
	mock.recorder = &MockExternalMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use.
func (m *MockExternal) EXPECT() *MockExternalMockRecorder {
	return m.recorder
}

// ValidateToken mocks base method.
func (m *MockExternal) ValidateToken(ctx context.Context, token string) (models.TokenData, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "ValidateToken", ctx, token)
	ret0, _ := ret[0].(models.TokenData)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// ValidateToken indicates an expected call of ValidateToken.
func (mr *MockExternalMockRecorder) ValidateToken(ctx, token interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "ValidateToken", reflect.TypeOf((*MockExternal)(nil).ValidateToken), ctx, token)
}
