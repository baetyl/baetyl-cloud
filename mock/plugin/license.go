// Code generated by MockGen. DO NOT EDIT.
// Source: github.com/baetyl/baetyl-cloud/v2/plugin (interfaces: License)

// Package plugin is a generated GoMock package.
package plugin

import (
	gomock "github.com/golang/mock/gomock"
	reflect "reflect"
)

// MockLicense is a mock of License interface
type MockLicense struct {
	ctrl     *gomock.Controller
	recorder *MockLicenseMockRecorder
}

// MockLicenseMockRecorder is the mock recorder for MockLicense
type MockLicenseMockRecorder struct {
	mock *MockLicense
}

// NewMockLicense creates a new mock instance
func NewMockLicense(ctrl *gomock.Controller) *MockLicense {
	mock := &MockLicense{ctrl: ctrl}
	mock.recorder = &MockLicenseMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use
func (m *MockLicense) EXPECT() *MockLicenseMockRecorder {
	return m.recorder
}

// AcquireQuota mocks base method
func (m *MockLicense) AcquireQuota(arg0, arg1 string, arg2 int) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "AcquireQuota", arg0, arg1, arg2)
	ret0, _ := ret[0].(error)
	return ret0
}

// AcquireQuota indicates an expected call of AcquireQuota
func (mr *MockLicenseMockRecorder) AcquireQuota(arg0, arg1, arg2 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "AcquireQuota", reflect.TypeOf((*MockLicense)(nil).AcquireQuota), arg0, arg1, arg2)
}

// CheckLicense mocks base method
func (m *MockLicense) CheckLicense() error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "CheckLicense")
	ret0, _ := ret[0].(error)
	return ret0
}

// CheckLicense indicates an expected call of CheckLicense
func (mr *MockLicenseMockRecorder) CheckLicense() *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "CheckLicense", reflect.TypeOf((*MockLicense)(nil).CheckLicense))
}

// Close mocks base method
func (m *MockLicense) Close() error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Close")
	ret0, _ := ret[0].(error)
	return ret0
}

// Close indicates an expected call of Close
func (mr *MockLicenseMockRecorder) Close() *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Close", reflect.TypeOf((*MockLicense)(nil).Close))
}

// CreateQuota mocks base method
func (m *MockLicense) CreateQuota(arg0 string, arg1 map[string]int) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "CreateQuota", arg0, arg1)
	ret0, _ := ret[0].(error)
	return ret0
}

// CreateQuota indicates an expected call of CreateQuota
func (mr *MockLicenseMockRecorder) CreateQuota(arg0, arg1 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "CreateQuota", reflect.TypeOf((*MockLicense)(nil).CreateQuota), arg0, arg1)
}

// DeleteQuota mocks base method
func (m *MockLicense) DeleteQuota(arg0, arg1 string) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "DeleteQuota", arg0, arg1)
	ret0, _ := ret[0].(error)
	return ret0
}

// DeleteQuota indicates an expected call of DeleteQuota
func (mr *MockLicenseMockRecorder) DeleteQuota(arg0, arg1 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "DeleteQuota", reflect.TypeOf((*MockLicense)(nil).DeleteQuota), arg0, arg1)
}

// DeleteQuotaByNamespace mocks base method
func (m *MockLicense) DeleteQuotaByNamespace(arg0 string) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "DeleteQuotaByNamespace", arg0)
	ret0, _ := ret[0].(error)
	return ret0
}

// DeleteQuotaByNamespace indicates an expected call of DeleteQuotaByNamespace
func (mr *MockLicenseMockRecorder) DeleteQuotaByNamespace(arg0 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "DeleteQuotaByNamespace", reflect.TypeOf((*MockLicense)(nil).DeleteQuotaByNamespace), arg0)
}

// GetDefaultQuotas mocks base method
func (m *MockLicense) GetDefaultQuotas(arg0 string) (map[string]int, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetDefaultQuotas", arg0)
	ret0, _ := ret[0].(map[string]int)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// GetDefaultQuotas indicates an expected call of GetDefaultQuotas
func (mr *MockLicenseMockRecorder) GetDefaultQuotas(arg0 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetDefaultQuotas", reflect.TypeOf((*MockLicense)(nil).GetDefaultQuotas), arg0)
}

// GetQuota mocks base method
func (m *MockLicense) GetQuota(arg0 string) (map[string]int, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetQuota", arg0)
	ret0, _ := ret[0].(map[string]int)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// GetQuota indicates an expected call of GetQuota
func (mr *MockLicenseMockRecorder) GetQuota(arg0 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetQuota", reflect.TypeOf((*MockLicense)(nil).GetQuota), arg0)
}

// ProtectCode mocks base method
func (m *MockLicense) ProtectCode() error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "ProtectCode")
	ret0, _ := ret[0].(error)
	return ret0
}

// ProtectCode indicates an expected call of ProtectCode
func (mr *MockLicenseMockRecorder) ProtectCode() *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "ProtectCode", reflect.TypeOf((*MockLicense)(nil).ProtectCode))
}

// ReleaseQuota mocks base method
func (m *MockLicense) ReleaseQuota(arg0, arg1 string, arg2 int) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "ReleaseQuota", arg0, arg1, arg2)
	ret0, _ := ret[0].(error)
	return ret0
}

// ReleaseQuota indicates an expected call of ReleaseQuota
func (mr *MockLicenseMockRecorder) ReleaseQuota(arg0, arg1, arg2 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "ReleaseQuota", reflect.TypeOf((*MockLicense)(nil).ReleaseQuota), arg0, arg1, arg2)
}

// UpdateQuota mocks base method
func (m *MockLicense) UpdateQuota(arg0, arg1 string, arg2 int) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "UpdateQuota", arg0, arg1, arg2)
	ret0, _ := ret[0].(error)
	return ret0
}

// UpdateQuota indicates an expected call of UpdateQuota
func (mr *MockLicenseMockRecorder) UpdateQuota(arg0, arg1, arg2 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "UpdateQuota", reflect.TypeOf((*MockLicense)(nil).UpdateQuota), arg0, arg1, arg2)
}
