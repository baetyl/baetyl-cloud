// Code generated by MockGen. DO NOT EDIT.
// Source: github.com/baetyl/baetyl-cloud/service (interfaces: PKIService)

// Package plugin is a generated GoMock package.
package plugin

import (
	models "github.com/baetyl/baetyl-cloud/v2/models"
	gomock "github.com/golang/mock/gomock"
	reflect "reflect"
)

// MockPKIService is a mock of PKIService interface
type MockPKIService struct {
	ctrl     *gomock.Controller
	recorder *MockPKIServiceMockRecorder
}

// MockPKIServiceMockRecorder is the mock recorder for MockPKIService
type MockPKIServiceMockRecorder struct {
	mock *MockPKIService
}

// NewMockPKIService creates a new mock instance
func NewMockPKIService(ctrl *gomock.Controller) *MockPKIService {
	mock := &MockPKIService{ctrl: ctrl}
	mock.recorder = &MockPKIServiceMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use
func (m *MockPKIService) EXPECT() *MockPKIServiceMockRecorder {
	return m.recorder
}

// DeleteClientCertificate mocks base method
func (m *MockPKIService) DeleteClientCertificate(arg0 string) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "DeleteClientCertificate", arg0)
	ret0, _ := ret[0].(error)
	return ret0
}

// DeleteClientCertificate indicates an expected call of DeleteClientCertificate
func (mr *MockPKIServiceMockRecorder) DeleteClientCertificate(arg0 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "DeleteClientCertificate", reflect.TypeOf((*MockPKIService)(nil).DeleteClientCertificate), arg0)
}

// DeleteServerCertificate mocks base method
func (m *MockPKIService) DeleteServerCertificate(arg0 string) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "DeleteServerCertificate", arg0)
	ret0, _ := ret[0].(error)
	return ret0
}

// DeleteServerCertificate indicates an expected call of DeleteServerCertificate
func (mr *MockPKIServiceMockRecorder) DeleteServerCertificate(arg0 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "DeleteServerCertificate", reflect.TypeOf((*MockPKIService)(nil).DeleteServerCertificate), arg0)
}

// GetCA mocks base method
func (m *MockPKIService) GetCA() ([]byte, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetCA")
	ret0, _ := ret[0].([]byte)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// GetCA indicates an expected call of GetCA
func (mr *MockPKIServiceMockRecorder) GetCA() *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetCA", reflect.TypeOf((*MockPKIService)(nil).GetCA))
}

// SignClientCertificate mocks base method
func (m *MockPKIService) SignClientCertificate(arg0 string, arg1 models.AltNames) (*models.PEMCredential, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "SignClientCertificate", arg0, arg1)
	ret0, _ := ret[0].(*models.PEMCredential)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// SignClientCertificate indicates an expected call of SignClientCertificate
func (mr *MockPKIServiceMockRecorder) SignClientCertificate(arg0, arg1 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "SignClientCertificate", reflect.TypeOf((*MockPKIService)(nil).SignClientCertificate), arg0, arg1)
}

// SignServerCertificate mocks base method
func (m *MockPKIService) SignServerCertificate(arg0 string, arg1 models.AltNames) (*models.PEMCredential, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "SignServerCertificate", arg0, arg1)
	ret0, _ := ret[0].(*models.PEMCredential)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// SignServerCertificate indicates an expected call of SignServerCertificate
func (mr *MockPKIServiceMockRecorder) SignServerCertificate(arg0, arg1 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "SignServerCertificate", reflect.TypeOf((*MockPKIService)(nil).SignServerCertificate), arg0, arg1)
}
