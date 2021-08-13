// Code generated by MockGen. DO NOT EDIT.
// Source: github.com/baetyl/baetyl-cloud/v2/plugin (interfaces: JWT)

// Package plugin is a generated GoMock package.
package plugin

import (
	common "github.com/baetyl/baetyl-cloud/v2/common"
	plugin "github.com/baetyl/baetyl-cloud/v2/plugin"
	gomock "github.com/golang/mock/gomock"
	reflect "reflect"
)

// MockJWT is a mock of JWT interface
type MockJWT struct {
	ctrl     *gomock.Controller
	recorder *MockJWTMockRecorder
}

// MockJWTMockRecorder is the mock recorder for MockJWT
type MockJWTMockRecorder struct {
	mock *MockJWT
}

// NewMockJWT creates a new mock instance
func NewMockJWT(ctrl *gomock.Controller) *MockJWT {
	mock := &MockJWT{ctrl: ctrl}
	mock.recorder = &MockJWTMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use
func (m *MockJWT) EXPECT() *MockJWTMockRecorder {
	return m.recorder
}

// CheckAndParseJWT mocks base method
func (m *MockJWT) CheckAndParseJWT(arg0 *common.Context) (map[string]interface{}, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "CheckAndParseJWT", arg0)
	ret0, _ := ret[0].(map[string]interface{})
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// CheckAndParseJWT indicates an expected call of CheckAndParseJWT
func (mr *MockJWTMockRecorder) CheckAndParseJWT(arg0 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "CheckAndParseJWT", reflect.TypeOf((*MockJWT)(nil).CheckAndParseJWT), arg0)
}

// Close mocks base method
func (m *MockJWT) Close() error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Close")
	ret0, _ := ret[0].(error)
	return ret0
}

// Close indicates an expected call of Close
func (mr *MockJWTMockRecorder) Close() *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Close", reflect.TypeOf((*MockJWT)(nil).Close))
}

// GenerateJWT mocks base method
func (m *MockJWT) GenerateJWT(arg0 *common.Context) (*plugin.JWTInfo, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GenerateJWT", arg0)
	ret0, _ := ret[0].(*plugin.JWTInfo)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// GenerateJWT indicates an expected call of GenerateJWT
func (mr *MockJWTMockRecorder) GenerateJWT(arg0 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GenerateJWT", reflect.TypeOf((*MockJWT)(nil).GenerateJWT), arg0)
}

// GetJWT mocks base method
func (m *MockJWT) GetJWT(arg0 *common.Context) (string, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetJWT", arg0)
	ret0, _ := ret[0].(string)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// GetJWT indicates an expected call of GetJWT
func (mr *MockJWTMockRecorder) GetJWT(arg0 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetJWT", reflect.TypeOf((*MockJWT)(nil).GetJWT), arg0)
}

// RefreshJWT mocks base method
func (m *MockJWT) RefreshJWT(arg0 *common.Context) (*plugin.JWTInfo, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "RefreshJWT", arg0)
	ret0, _ := ret[0].(*plugin.JWTInfo)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// RefreshJWT indicates an expected call of RefreshJWT
func (mr *MockJWTMockRecorder) RefreshJWT(arg0 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "RefreshJWT", reflect.TypeOf((*MockJWT)(nil).RefreshJWT), arg0)
}
