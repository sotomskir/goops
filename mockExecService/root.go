// Code generated by MockGen. DO NOT EDIT.
// Source: execService/root.go

// Package mock_execService is a generated GoMock package.
package mock_execService

import (
	gomock "github.com/golang/mock/gomock"
	reflect "reflect"
)

// MockIService is a mock of IService interface
type MockIService struct {
	ctrl     *gomock.Controller
	recorder *MockIServiceMockRecorder
}

// MockIServiceMockRecorder is the mock recorder for MockIService
type MockIServiceMockRecorder struct {
	mock *MockIService
}

// NewMockIService creates a new mock instance
func NewMockIService(ctrl *gomock.Controller) *MockIService {
	mock := &MockIService{ctrl: ctrl}
	mock.recorder = &MockIServiceMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use
func (m *MockIService) EXPECT() *MockIServiceMockRecorder {
	return m.recorder
}

// Exec mocks base method
func (m *MockIService) Exec(cmd string) (string, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Exec", cmd)
	ret0, _ := ret[0].(string)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// Exec indicates an expected call of Exec
func (mr *MockIServiceMockRecorder) Exec(cmd interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Exec", reflect.TypeOf((*MockIService)(nil).Exec), cmd)
}

// LogExec mocks base method
func (m *MockIService) LogExec(cmd string) {
	m.ctrl.T.Helper()
	m.ctrl.Call(m, "LogExec", cmd)
}

// LogExec indicates an expected call of LogExec
func (mr *MockIServiceMockRecorder) LogExec(cmd interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "LogExec", reflect.TypeOf((*MockIService)(nil).LogExec), cmd)
}