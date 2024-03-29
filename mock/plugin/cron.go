// Code generated by MockGen. DO NOT EDIT.
// Source: github.com/baetyl/baetyl-cloud/v2/plugin (interfaces: Cron)

// Package plugin is a generated GoMock package.
package plugin

import (
	models "github.com/baetyl/baetyl-cloud/v2/models"
	gomock "github.com/golang/mock/gomock"
	reflect "reflect"
)

// MockCron is a mock of Cron interface
type MockCron struct {
	ctrl     *gomock.Controller
	recorder *MockCronMockRecorder
}

// MockCronMockRecorder is the mock recorder for MockCron
type MockCronMockRecorder struct {
	mock *MockCron
}

// NewMockCron creates a new mock instance
func NewMockCron(ctrl *gomock.Controller) *MockCron {
	mock := &MockCron{ctrl: ctrl}
	mock.recorder = &MockCronMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use
func (m *MockCron) EXPECT() *MockCronMockRecorder {
	return m.recorder
}

// Close mocks base method
func (m *MockCron) Close() error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Close")
	ret0, _ := ret[0].(error)
	return ret0
}

// Close indicates an expected call of Close
func (mr *MockCronMockRecorder) Close() *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Close", reflect.TypeOf((*MockCron)(nil).Close))
}

// CreateCron mocks base method
func (m *MockCron) CreateCron(arg0 *models.Cron) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "CreateCron", arg0)
	ret0, _ := ret[0].(error)
	return ret0
}

// CreateCron indicates an expected call of CreateCron
func (mr *MockCronMockRecorder) CreateCron(arg0 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "CreateCron", reflect.TypeOf((*MockCron)(nil).CreateCron), arg0)
}

// DeleteCron mocks base method
func (m *MockCron) DeleteCron(arg0, arg1 string) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "DeleteCron", arg0, arg1)
	ret0, _ := ret[0].(error)
	return ret0
}

// DeleteCron indicates an expected call of DeleteCron
func (mr *MockCronMockRecorder) DeleteCron(arg0, arg1 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "DeleteCron", reflect.TypeOf((*MockCron)(nil).DeleteCron), arg0, arg1)
}

// DeleteExpiredApps mocks base method
func (m *MockCron) DeleteExpiredApps(arg0 []uint64) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "DeleteExpiredApps", arg0)
	ret0, _ := ret[0].(error)
	return ret0
}

// DeleteExpiredApps indicates an expected call of DeleteExpiredApps
func (mr *MockCronMockRecorder) DeleteExpiredApps(arg0 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "DeleteExpiredApps", reflect.TypeOf((*MockCron)(nil).DeleteExpiredApps), arg0)
}

// GetCron mocks base method
func (m *MockCron) GetCron(arg0, arg1 string) (*models.Cron, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetCron", arg0, arg1)
	ret0, _ := ret[0].(*models.Cron)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// GetCron indicates an expected call of GetCron
func (mr *MockCronMockRecorder) GetCron(arg0, arg1 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetCron", reflect.TypeOf((*MockCron)(nil).GetCron), arg0, arg1)
}

// ListExpiredApps mocks base method
func (m *MockCron) ListExpiredApps() ([]models.Cron, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "ListExpiredApps")
	ret0, _ := ret[0].([]models.Cron)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// ListExpiredApps indicates an expected call of ListExpiredApps
func (mr *MockCronMockRecorder) ListExpiredApps() *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "ListExpiredApps", reflect.TypeOf((*MockCron)(nil).ListExpiredApps))
}

// UpdateCron mocks base method
func (m *MockCron) UpdateCron(arg0 *models.Cron) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "UpdateCron", arg0)
	ret0, _ := ret[0].(error)
	return ret0
}

// UpdateCron indicates an expected call of UpdateCron
func (mr *MockCronMockRecorder) UpdateCron(arg0 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "UpdateCron", reflect.TypeOf((*MockCron)(nil).UpdateCron), arg0)
}
