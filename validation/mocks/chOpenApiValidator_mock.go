// Code generated by MockGen. DO NOT EDIT.
// Source: chOpenApiValidator.go

// Package mocks is a generated GoMock package.
package mocks

import (
	gomock "github.com/golang/mock/gomock"
	http "net/http"
	reflect "reflect"
)

// MockCHValidator is a mock of CHValidator interface
type MockCHValidator struct {
	ctrl     *gomock.Controller
	recorder *MockCHValidatorMockRecorder
}

// MockCHValidatorMockRecorder is the mock recorder for MockCHValidator
type MockCHValidatorMockRecorder struct {
	mock *MockCHValidator
}

// NewMockCHValidator creates a new mock instance
func NewMockCHValidator(ctrl *gomock.Controller) *MockCHValidator {
	mock := &MockCHValidator{ctrl: ctrl}
	mock.recorder = &MockCHValidatorMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use
func (m *MockCHValidator) EXPECT() *MockCHValidatorMockRecorder {
	return m.recorder
}

// ValidateRequestAgainstOpenApiSpec mocks base method
func (m *MockCHValidator) ValidateRequestAgainstOpenApiSpec(httpReq *http.Request, openApiSpec string) ([]byte, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "ValidateRequestAgainstOpenApiSpec", httpReq, openApiSpec)
	ret0, _ := ret[0].([]byte)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// ValidateRequestAgainstOpenApiSpec indicates an expected call of ValidateRequestAgainstOpenApiSpec
func (mr *MockCHValidatorMockRecorder) ValidateRequestAgainstOpenApiSpec(httpReq, openApiSpec interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "ValidateRequestAgainstOpenApiSpec", reflect.TypeOf((*MockCHValidator)(nil).ValidateRequestAgainstOpenApiSpec), httpReq, openApiSpec)
}