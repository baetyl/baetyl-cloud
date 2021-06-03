// Code generated by MockGen. DO NOT EDIT.
// Source: github.com/baetyl/baetyl-cloud/v2/plugin (interfaces: Auth)

// Package plugin is a generated GoMock package.
package plugin

import (
	common "github.com/baetyl/baetyl-cloud/v2/common"
	plugin "github.com/baetyl/baetyl-cloud/v2/plugin"
	gomock "github.com/golang/mock/gomock"
	reflect "reflect"
)

// MockAuth is a mock of Auth interface
type MockAuth struct {
	ctrl     *gomock.Controller
	recorder *MockAuthMockRecorder
}

// MockAuthMockRecorder is the mock recorder for MockAuth
type MockAuthMockRecorder struct {
	mock *MockAuth
}

// NewMockAuth creates a new mock instance
func NewMockAuth(ctrl *gomock.Controller) *MockAuth {
	mock := &MockAuth{ctrl: ctrl}
	mock.recorder = &MockAuthMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use
func (m *MockAuth) EXPECT() *MockAuthMockRecorder {
	return m.recorder
}

// AuthAndVerify mocks base method
func (m *MockAuth) AuthAndVerify(arg0 *common.Context, arg1 *plugin.PermissionRequest) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "AuthAndVerify", arg0, arg1)
	ret0, _ := ret[0].(error)
	return ret0
}

// AuthAndVerify indicates an expected call of AuthAndVerify
func (mr *MockAuthMockRecorder) AuthAndVerify(arg0, arg1 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "AuthAndVerify", reflect.TypeOf((*MockAuth)(nil).AuthAndVerify), arg0, arg1)
}

// Authenticate mocks base method
func (m *MockAuth) Authenticate(arg0 *common.Context) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Authenticate", arg0)
	ret0, _ := ret[0].(error)
	return ret0
}

// Authenticate indicates an expected call of Authenticate
func (mr *MockAuthMockRecorder) Authenticate(arg0 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Authenticate", reflect.TypeOf((*MockAuth)(nil).Authenticate), arg0)
}

// Close mocks base method
func (m *MockAuth) Close() error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Close")
	ret0, _ := ret[0].(error)
	return ret0
}

// Close indicates an expected call of Close
func (mr *MockAuthMockRecorder) Close() *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Close", reflect.TypeOf((*MockAuth)(nil).Close))
}

// Verify mocks base method
func (m *MockAuth) Verify(arg0 *common.Context, arg1 *plugin.PermissionRequest) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Verify", arg0, arg1)
	ret0, _ := ret[0].(error)
	return ret0
}

// Verify indicates an expected call of Verify
func (mr *MockAuthMockRecorder) Verify(arg0, arg1 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Verify", reflect.TypeOf((*MockAuth)(nil).Verify), arg0, arg1)
}
