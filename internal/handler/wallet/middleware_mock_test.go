// Code generated by MockGen. DO NOT EDIT.
// Source: middleware.go

// Package wallet is a generated GoMock package.
package wallet

import (
	reflect "reflect"

	gin "github.com/gin-gonic/gin"
	gomock "github.com/golang/mock/gomock"
)

// MockMiddleware is a mock of Middleware interface.
type MockMiddleware struct {
	ctrl     *gomock.Controller
	recorder *MockMiddlewareMockRecorder
}

// MockMiddlewareMockRecorder is the mock recorder for MockMiddleware.
type MockMiddlewareMockRecorder struct {
	mock *MockMiddleware
}

// NewMockMiddleware creates a new mock instance.
func NewMockMiddleware(ctrl *gomock.Controller) *MockMiddleware {
	mock := &MockMiddleware{ctrl: ctrl}
	mock.recorder = &MockMiddlewareMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use.
func (m *MockMiddleware) EXPECT() *MockMiddlewareMockRecorder {
	return m.recorder
}

// MiddlewareSignatureValidation mocks base method.
func (m *MockMiddleware) MiddlewareSignatureValidation(c *gin.Context) {
	m.ctrl.T.Helper()
	m.ctrl.Call(m, "MiddlewareSignatureValidation", c)
}

// MiddlewareSignatureValidation indicates an expected call of MiddlewareSignatureValidation.
func (mr *MockMiddlewareMockRecorder) MiddlewareSignatureValidation(c interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "MiddlewareSignatureValidation", reflect.TypeOf((*MockMiddleware)(nil).MiddlewareSignatureValidation), c)
}

// MiddlewareValidateToken mocks base method.
func (m *MockMiddleware) MiddlewareValidateToken(c *gin.Context) {
	m.ctrl.T.Helper()
	m.ctrl.Call(m, "MiddlewareValidateToken", c)
}

// MiddlewareValidateToken indicates an expected call of MiddlewareValidateToken.
func (mr *MockMiddlewareMockRecorder) MiddlewareValidateToken(c interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "MiddlewareValidateToken", reflect.TypeOf((*MockMiddleware)(nil).MiddlewareValidateToken), c)
}
