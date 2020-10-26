// Code generated by MockGen. DO NOT EDIT.
// Source: github.com/baetyl/baetyl-cloud/v2/service (interfaces: NodeService)

// Package service is a generated GoMock package.
package service

import (
	models "github.com/baetyl/baetyl-cloud/v2/models"
	v1 "github.com/baetyl/baetyl-go/v2/spec/v1"
	gomock "github.com/golang/mock/gomock"
	reflect "reflect"
)

// MockNodeService is a mock of NodeService interface
type MockNodeService struct {
	ctrl     *gomock.Controller
	recorder *MockNodeServiceMockRecorder
}

// MockNodeServiceMockRecorder is the mock recorder for MockNodeService
type MockNodeServiceMockRecorder struct {
	mock *MockNodeService
}

// NewMockNodeService creates a new mock instance
func NewMockNodeService(ctrl *gomock.Controller) *MockNodeService {
	mock := &MockNodeService{ctrl: ctrl}
	mock.recorder = &MockNodeServiceMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use
func (m *MockNodeService) EXPECT() *MockNodeServiceMockRecorder {
	return m.recorder
}

// CountNumber mocks base method
func (m *MockNodeService) CountNumber(arg0 string) (map[string]int, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "CountNumber", arg0)
	ret0, _ := ret[0].(map[string]int)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// CountNumber indicates an expected call of CountNumber
func (mr *MockNodeServiceMockRecorder) CountNumber(arg0 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "CountNumber", reflect.TypeOf((*MockNodeService)(nil).CountNumber), arg0)
}

// Create mocks base method
func (m *MockNodeService) Create(arg0 string, arg1 *v1.Node) (*v1.Node, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Create", arg0, arg1)
	ret0, _ := ret[0].(*v1.Node)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// Create indicates an expected call of Create
func (mr *MockNodeServiceMockRecorder) Create(arg0, arg1 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Create", reflect.TypeOf((*MockNodeService)(nil).Create), arg0, arg1)
}

// Delete mocks base method
func (m *MockNodeService) Delete(arg0, arg1 string) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Delete", arg0, arg1)
	ret0, _ := ret[0].(error)
	return ret0
}

// Delete indicates an expected call of Delete
func (mr *MockNodeServiceMockRecorder) Delete(arg0, arg1 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Delete", reflect.TypeOf((*MockNodeService)(nil).Delete), arg0, arg1)
}

// DeleteNodeAppVersion mocks base method
func (m *MockNodeService) DeleteNodeAppVersion(arg0 string, arg1 *v1.Application) ([]string, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "DeleteNodeAppVersion", arg0, arg1)
	ret0, _ := ret[0].([]string)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// DeleteNodeAppVersion indicates an expected call of DeleteNodeAppVersion
func (mr *MockNodeServiceMockRecorder) DeleteNodeAppVersion(arg0, arg1 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "DeleteNodeAppVersion", reflect.TypeOf((*MockNodeService)(nil).DeleteNodeAppVersion), arg0, arg1)
}

// Get mocks base method
func (m *MockNodeService) Get(arg0, arg1 string) (*v1.Node, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Get", arg0, arg1)
	ret0, _ := ret[0].(*v1.Node)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// Get indicates an expected call of Get
func (mr *MockNodeServiceMockRecorder) Get(arg0, arg1 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Get", reflect.TypeOf((*MockNodeService)(nil).Get), arg0, arg1)
}

// GetDesire mocks base method
func (m *MockNodeService) GetDesire(arg0, arg1 string) (*v1.Desire, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetDesire", arg0, arg1)
	ret0, _ := ret[0].(*v1.Desire)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// GetDesire indicates an expected call of GetDesire
func (mr *MockNodeServiceMockRecorder) GetDesire(arg0, arg1 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetDesire", reflect.TypeOf((*MockNodeService)(nil).GetDesire), arg0, arg1)
}

// List mocks base method
func (m *MockNodeService) List(arg0 string, arg1 *models.ListOptions) (*models.NodeList, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "List", arg0, arg1)
	ret0, _ := ret[0].(*models.NodeList)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// List indicates an expected call of List
func (mr *MockNodeServiceMockRecorder) List(arg0, arg1 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "List", reflect.TypeOf((*MockNodeService)(nil).List), arg0, arg1)
}

// Update mocks base method
func (m *MockNodeService) Update(arg0 string, arg1 *v1.Node) (*v1.Node, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Update", arg0, arg1)
	ret0, _ := ret[0].(*v1.Node)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// Update indicates an expected call of Update
func (mr *MockNodeServiceMockRecorder) Update(arg0, arg1 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Update", reflect.TypeOf((*MockNodeService)(nil).Update), arg0, arg1)
}

// UpdateDesire mocks base method
func (m *MockNodeService) UpdateDesire(arg0, arg1 string, arg2 v1.Desire) (*models.Shadow, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "UpdateDesire", arg0, arg1, arg2)
	ret0, _ := ret[0].(*models.Shadow)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// UpdateDesire indicates an expected call of UpdateDesire
func (mr *MockNodeServiceMockRecorder) UpdateDesire(arg0, arg1, arg2 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "UpdateDesire", reflect.TypeOf((*MockNodeService)(nil).UpdateDesire), arg0, arg1, arg2)
}

// UpdateNodeAppVersion mocks base method
func (m *MockNodeService) UpdateNodeAppVersion(arg0 string, arg1 *v1.Application) ([]string, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "UpdateNodeAppVersion", arg0, arg1)
	ret0, _ := ret[0].([]string)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// UpdateNodeAppVersion indicates an expected call of UpdateNodeAppVersion
func (mr *MockNodeServiceMockRecorder) UpdateNodeAppVersion(arg0, arg1 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "UpdateNodeAppVersion", reflect.TypeOf((*MockNodeService)(nil).UpdateNodeAppVersion), arg0, arg1)
}

// UpdateReport mocks base method
func (m *MockNodeService) UpdateReport(arg0, arg1 string, arg2 v1.Report) (*models.Shadow, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "UpdateReport", arg0, arg1, arg2)
	ret0, _ := ret[0].(*models.Shadow)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// UpdateReport indicates an expected call of UpdateReport
func (mr *MockNodeServiceMockRecorder) UpdateReport(arg0, arg1, arg2 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "UpdateReport", reflect.TypeOf((*MockNodeService)(nil).UpdateReport), arg0, arg1, arg2)
}
