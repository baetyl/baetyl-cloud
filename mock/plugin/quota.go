// Code generated by MockGen. DO NOT EDIT.
// Source: github.com/baetyl/baetyl-cloud/v2/plugin (interfaces: Quota)

// Package plugin is a generated GoMock package.
package plugin

import (
	gomock "github.com/golang/mock/gomock"
	reflect "reflect"
)

// MockQuota is a mock of Quota interface
type MockQuota struct {
	ctrl     *gomock.Controller
	recorder *MockQuotaMockRecorder
}

// MockQuotaMockRecorder is the mock recorder for MockQuota
type MockQuotaMockRecorder struct {
	mock *MockQuota
}

// NewMockQuota creates a new mock instance
func NewMockQuota(ctrl *gomock.Controller) *MockQuota {
	mock := &MockQuota{ctrl: ctrl}
	mock.recorder = &MockQuotaMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use
func (m *MockQuota) EXPECT() *MockQuotaMockRecorder {
	return m.recorder
}

// AcquireQuota mocks base method
func (m *MockQuota) AcquireQuota(arg0, arg1 string, arg2 int) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "AcquireQuota", arg0, arg1, arg2)
	ret0, _ := ret[0].(error)
	return ret0
}

// AcquireQuota indicates an expected call of AcquireQuota
func (mr *MockQuotaMockRecorder) AcquireQuota(arg0, arg1, arg2 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "AcquireQuota", reflect.TypeOf((*MockQuota)(nil).AcquireQuota), arg0, arg1, arg2)
}

// Close mocks base method
func (m *MockQuota) Close() error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Close")
	ret0, _ := ret[0].(error)
	return ret0
}

// Close indicates an expected call of Close
func (mr *MockQuotaMockRecorder) Close() *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Close", reflect.TypeOf((*MockQuota)(nil).Close))
}

// CreateQuota mocks base method
func (m *MockQuota) CreateQuota(arg0 string, arg1 map[string]int) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "CreateQuota", arg0, arg1)
	ret0, _ := ret[0].(error)
	return ret0
}

// CreateQuota indicates an expected call of CreateQuota
func (mr *MockQuotaMockRecorder) CreateQuota(arg0, arg1 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "CreateQuota", reflect.TypeOf((*MockQuota)(nil).CreateQuota), arg0, arg1)
}

// DeleteQuota mocks base method
func (m *MockQuota) DeleteQuota(arg0, arg1 string) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "DeleteQuota", arg0, arg1)
	ret0, _ := ret[0].(error)
	return ret0
}

// DeleteQuota indicates an expected call of DeleteQuota
func (mr *MockQuotaMockRecorder) DeleteQuota(arg0, arg1 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "DeleteQuota", reflect.TypeOf((*MockQuota)(nil).DeleteQuota), arg0, arg1)
}

// DeleteQuotaByNamespace mocks base method
func (m *MockQuota) DeleteQuotaByNamespace(arg0 string) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "DeleteQuotaByNamespace", arg0)
	ret0, _ := ret[0].(error)
	return ret0
}

// DeleteQuotaByNamespace indicates an expected call of DeleteQuotaByNamespace
func (mr *MockQuotaMockRecorder) DeleteQuotaByNamespace(arg0 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "DeleteQuotaByNamespace", reflect.TypeOf((*MockQuota)(nil).DeleteQuotaByNamespace), arg0)
}

// GetDefaultQuotas mocks base method
func (m *MockQuota) GetDefaultQuotas(arg0 string) (map[string]int, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetDefaultQuotas", arg0)
	ret0, _ := ret[0].(map[string]int)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// GetDefaultQuotas indicates an expected call of GetDefaultQuotas
func (mr *MockQuotaMockRecorder) GetDefaultQuotas(arg0 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetDefaultQuotas", reflect.TypeOf((*MockQuota)(nil).GetDefaultQuotas), arg0)
}

// GetQuota mocks base method
func (m *MockQuota) GetQuota(arg0 string) (map[string]int, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetQuota", arg0)
	ret0, _ := ret[0].(map[string]int)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// GetQuota indicates an expected call of GetQuota
func (mr *MockQuotaMockRecorder) GetQuota(arg0 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetQuota", reflect.TypeOf((*MockQuota)(nil).GetQuota), arg0)
}

// ReleaseQuota mocks base method
func (m *MockQuota) ReleaseQuota(arg0, arg1 string, arg2 int) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "ReleaseQuota", arg0, arg1, arg2)
	ret0, _ := ret[0].(error)
	return ret0
}

// ReleaseQuota indicates an expected call of ReleaseQuota
func (mr *MockQuotaMockRecorder) ReleaseQuota(arg0, arg1, arg2 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "ReleaseQuota", reflect.TypeOf((*MockQuota)(nil).ReleaseQuota), arg0, arg1, arg2)
}

// UpdateQuota mocks base method
func (m *MockQuota) UpdateQuota(arg0, arg1 string, arg2 int) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "UpdateQuota", arg0, arg1, arg2)
	ret0, _ := ret[0].(error)
	return ret0
}

// UpdateQuota indicates an expected call of UpdateQuota
func (mr *MockQuotaMockRecorder) UpdateQuota(arg0, arg1, arg2 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "UpdateQuota", reflect.TypeOf((*MockQuota)(nil).UpdateQuota), arg0, arg1, arg2)
}
