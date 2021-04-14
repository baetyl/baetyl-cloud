// Code generated by MockGen. DO NOT EDIT.
// Source: github.com/baetyl/baetyl-cloud/v2/service (interfaces: SystemAppService)

// Package service is a generated GoMock package.
package service

import (
	v1 "github.com/baetyl/baetyl-go/v2/spec/v1"
	gomock "github.com/golang/mock/gomock"
	reflect "reflect"
)

// MockSystemAppService is a mock of SystemAppService interface
type MockSystemAppService struct {
	ctrl     *gomock.Controller
	recorder *MockSystemAppServiceMockRecorder
}

// MockSystemAppServiceMockRecorder is the mock recorder for MockSystemAppService
type MockSystemAppServiceMockRecorder struct {
	mock *MockSystemAppService
}

// NewMockSystemAppService creates a new mock instance
func NewMockSystemAppService(ctrl *gomock.Controller) *MockSystemAppService {
	mock := &MockSystemAppService{ctrl: ctrl}
	mock.recorder = &MockSystemAppServiceMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use
func (m *MockSystemAppService) EXPECT() *MockSystemAppServiceMockRecorder {
	return m.recorder
}

// GenApps mocks base method
func (m *MockSystemAppService) GenApps(arg0 string, arg1 *v1.Node) ([]*v1.Application, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GenApps", arg0, arg1)
	ret0, _ := ret[0].([]*v1.Application)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// GenApps indicates an expected call of GenApps
func (mr *MockSystemAppServiceMockRecorder) GenApps(arg0, arg1 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GenApps", reflect.TypeOf((*MockSystemAppService)(nil).GenApps), arg0, arg1)
}

// GenOptionalApps mocks base method
func (m *MockSystemAppService) GenOptionalApps(arg0, arg1 string, arg2 []string) ([]*v1.Application, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GenOptionalApps", arg0, arg1, arg2)
	ret0, _ := ret[0].([]*v1.Application)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// GenOptionalApps indicates an expected call of GenOptionalApps
func (mr *MockSystemAppServiceMockRecorder) GenOptionalApps(arg0, arg1, arg2 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GenOptionalApps", reflect.TypeOf((*MockSystemAppService)(nil).GenOptionalApps), arg0, arg1, arg2)
}

// GetOptionalApps mocks base method
func (m *MockSystemAppService) GetOptionalApps() []string {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetOptionalApps")
	ret0, _ := ret[0].([]string)
	return ret0
}

// GetOptionalApps indicates an expected call of GetOptionalApps
func (mr *MockSystemAppServiceMockRecorder) GetOptionalApps() *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetOptionalApps", reflect.TypeOf((*MockSystemAppService)(nil).GetOptionalApps))
}
