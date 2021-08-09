// Code generated by MockGen. DO NOT EDIT.
// Source: github.com/baetyl/baetyl-cloud/v2/service (interfaces: CronService)

// Package service is a generated GoMock package.
package service

import (
	models "github.com/baetyl/baetyl-cloud/v2/models"
	gomock "github.com/golang/mock/gomock"
	reflect "reflect"
)

// MockCronService is a mock of CronService interface
type MockCronService struct {
	ctrl     *gomock.Controller
	recorder *MockCronServiceMockRecorder
}

// MockCronServiceMockRecorder is the mock recorder for MockCronService
type MockCronServiceMockRecorder struct {
	mock *MockCronService
}

// NewMockCronService creates a new mock instance
func NewMockCronService(ctrl *gomock.Controller) *MockCronService {
	mock := &MockCronService{ctrl: ctrl}
	mock.recorder = &MockCronServiceMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use
func (m *MockCronService) EXPECT() *MockCronServiceMockRecorder {
	return m.recorder
}

// CreateCron mocks base method
func (m *MockCronService) CreateCron(arg0 *models.Cron) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "CreateCron", arg0)
	ret0, _ := ret[0].(error)
	return ret0
}

// CreateCron indicates an expected call of CreateCron
func (mr *MockCronServiceMockRecorder) CreateCron(arg0 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "CreateCron", reflect.TypeOf((*MockCronService)(nil).CreateCron), arg0)
}

// DeleteCron mocks base method
func (m *MockCronService) DeleteCron(arg0, arg1 string) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "DeleteCron", arg0, arg1)
	ret0, _ := ret[0].(error)
	return ret0
}

// DeleteCron indicates an expected call of DeleteCron
func (mr *MockCronServiceMockRecorder) DeleteCron(arg0, arg1 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "DeleteCron", reflect.TypeOf((*MockCronService)(nil).DeleteCron), arg0, arg1)
}

// DeleteExpiredApps mocks base method
func (m *MockCronService) DeleteExpiredApps(arg0 []uint64) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "DeleteExpiredApps", arg0)
	ret0, _ := ret[0].(error)
	return ret0
}

// DeleteExpiredApps indicates an expected call of DeleteExpiredApps
func (mr *MockCronServiceMockRecorder) DeleteExpiredApps(arg0 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "DeleteExpiredApps", reflect.TypeOf((*MockCronService)(nil).DeleteExpiredApps), arg0)
}

// GetCron mocks base method
func (m *MockCronService) GetCron(arg0, arg1 string) (*models.Cron, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetCron", arg0, arg1)
	ret0, _ := ret[0].(*models.Cron)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// GetCron indicates an expected call of GetCron
func (mr *MockCronServiceMockRecorder) GetCron(arg0, arg1 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetCron", reflect.TypeOf((*MockCronService)(nil).GetCron), arg0, arg1)
}

// ListExpiredApps mocks base method
func (m *MockCronService) ListExpiredApps() ([]models.Cron, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "ListExpiredApps")
	ret0, _ := ret[0].([]models.Cron)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// ListExpiredApps indicates an expected call of ListExpiredApps
func (mr *MockCronServiceMockRecorder) ListExpiredApps() *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "ListExpiredApps", reflect.TypeOf((*MockCronService)(nil).ListExpiredApps))
}

// UpdateCron mocks base method
func (m *MockCronService) UpdateCron(arg0 *models.Cron) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "UpdateCron", arg0)
	ret0, _ := ret[0].(error)
	return ret0
}

// UpdateCron indicates an expected call of UpdateCron
func (mr *MockCronServiceMockRecorder) UpdateCron(arg0 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "UpdateCron", reflect.TypeOf((*MockCronService)(nil).UpdateCron), arg0)
}
