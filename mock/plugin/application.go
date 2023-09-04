// Code generated by MockGen. DO NOT EDIT.
// Source: github.com/baetyl/baetyl-cloud/v2/plugin (interfaces: Application)

// Package plugin is a generated GoMock package.
package plugin

import (
	reflect "reflect"

	models "github.com/baetyl/baetyl-cloud/v2/models"
	v1 "github.com/baetyl/baetyl-go/v2/spec/v1"
	gomock "github.com/golang/mock/gomock"
)

// MockApplication is a mock of Application interface.
type MockApplication struct {
	ctrl     *gomock.Controller
	recorder *MockApplicationMockRecorder
}

// MockApplicationMockRecorder is the mock recorder for MockApplication.
type MockApplicationMockRecorder struct {
	mock *MockApplication
}

// NewMockApplication creates a new mock instance.
func NewMockApplication(ctrl *gomock.Controller) *MockApplication {
	mock := &MockApplication{ctrl: ctrl}
	mock.recorder = &MockApplicationMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use.
func (m *MockApplication) EXPECT() *MockApplicationMockRecorder {
	return m.recorder
}

// CreateApplication mocks base method.
func (m *MockApplication) CreateApplication(arg0 interface{}, arg1 string, arg2 *v1.Application) (*v1.Application, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "CreateApplication", arg0, arg1, arg2)
	ret0, _ := ret[0].(*v1.Application)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// CreateApplication indicates an expected call of CreateApplication.
func (mr *MockApplicationMockRecorder) CreateApplication(arg0, arg1, arg2 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "CreateApplication", reflect.TypeOf((*MockApplication)(nil).CreateApplication), arg0, arg1, arg2)
}

// DeleteApplication mocks base method.
func (m *MockApplication) DeleteApplication(arg0 interface{}, arg1, arg2 string) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "DeleteApplication", arg0, arg1, arg2)
	ret0, _ := ret[0].(error)
	return ret0
}

// DeleteApplication indicates an expected call of DeleteApplication.
func (mr *MockApplicationMockRecorder) DeleteApplication(arg0, arg1, arg2 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "DeleteApplication", reflect.TypeOf((*MockApplication)(nil).DeleteApplication), arg0, arg1, arg2)
}

// GetApplication mocks base method.
func (m *MockApplication) GetApplication(arg0 interface{}, arg1, arg2, arg3 string) (*v1.Application, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetApplication", arg0, arg1, arg2, arg3)
	ret0, _ := ret[0].(*v1.Application)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// GetApplication indicates an expected call of GetApplication.
func (mr *MockApplicationMockRecorder) GetApplication(arg0, arg1, arg2, arg3 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetApplication", reflect.TypeOf((*MockApplication)(nil).GetApplication), arg0, arg1, arg2, arg3)
}

// ListApplication mocks base method.
func (m *MockApplication) ListApplication(arg0 interface{}, arg1 string, arg2 *models.ListOptions) (*models.ApplicationList, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "ListApplication", arg0, arg1, arg2)
	ret0, _ := ret[0].(*models.ApplicationList)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// ListApplication indicates an expected call of ListApplication.
func (mr *MockApplicationMockRecorder) ListApplication(arg0, arg1, arg2 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "ListApplication", reflect.TypeOf((*MockApplication)(nil).ListApplication), arg0, arg1, arg2)
}

// ListApplicationsByNames mocks base method.
func (m *MockApplication) ListApplicationsByNames(arg0 interface{}, arg1 string, arg2 []string) ([]models.AppItem, int, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "ListApplicationsByNames", arg0, arg1, arg2)
	ret0, _ := ret[0].([]models.AppItem)
	ret1, _ := ret[1].(int)
	ret2, _ := ret[2].(error)
	return ret0, ret1, ret2
}

// ListApplicationsByNames indicates an expected call of ListApplicationsByNames.
func (mr *MockApplicationMockRecorder) ListApplicationsByNames(arg0, arg1, arg2 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "ListApplicationsByNames", reflect.TypeOf((*MockApplication)(nil).ListApplicationsByNames), arg0, arg1, arg2)
}

// UpdateApplication mocks base method.
func (m *MockApplication) UpdateApplication(arg0 interface{}, arg1 string, arg2 *v1.Application) (*v1.Application, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "UpdateApplication", arg0, arg1, arg2)
	ret0, _ := ret[0].(*v1.Application)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// UpdateApplication indicates an expected call of UpdateApplication.
func (mr *MockApplicationMockRecorder) UpdateApplication(arg0, arg1, arg2 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "UpdateApplication", reflect.TypeOf((*MockApplication)(nil).UpdateApplication), arg0, arg1, arg2)
}
