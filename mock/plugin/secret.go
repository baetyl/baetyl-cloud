// Code generated by MockGen. DO NOT EDIT.
// Source: github.com/baetyl/baetyl-cloud/v2/plugin (interfaces: Secret)

// Package plugin is a generated GoMock package.
package plugin

import (
	reflect "reflect"

	models "github.com/baetyl/baetyl-cloud/v2/models"
	v1 "github.com/baetyl/baetyl-go/v2/spec/v1"
	gomock "github.com/golang/mock/gomock"
)

// MockSecret is a mock of Secret interface.
type MockSecret struct {
	ctrl     *gomock.Controller
	recorder *MockSecretMockRecorder
}

// MockSecretMockRecorder is the mock recorder for MockSecret.
type MockSecretMockRecorder struct {
	mock *MockSecret
}

// NewMockSecret creates a new mock instance.
func NewMockSecret(ctrl *gomock.Controller) *MockSecret {
	mock := &MockSecret{ctrl: ctrl}
	mock.recorder = &MockSecretMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use.
func (m *MockSecret) EXPECT() *MockSecretMockRecorder {
	return m.recorder
}

// CreateSecret mocks base method.
func (m *MockSecret) CreateSecret(arg0 interface{}, arg1 string, arg2 *v1.Secret) (*v1.Secret, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "CreateSecret", arg0, arg1, arg2)
	ret0, _ := ret[0].(*v1.Secret)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// CreateSecret indicates an expected call of CreateSecret.
func (mr *MockSecretMockRecorder) CreateSecret(arg0, arg1, arg2 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "CreateSecret", reflect.TypeOf((*MockSecret)(nil).CreateSecret), arg0, arg1, arg2)
}

// DeleteSecret mocks base method.
func (m *MockSecret) DeleteSecret(arg0 interface{}, arg1, arg2 string) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "DeleteSecret", arg0, arg1, arg2)
	ret0, _ := ret[0].(error)
	return ret0
}

// DeleteSecret indicates an expected call of DeleteSecret.
func (mr *MockSecretMockRecorder) DeleteSecret(arg0, arg1, arg2 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "DeleteSecret", reflect.TypeOf((*MockSecret)(nil).DeleteSecret), arg0, arg1, arg2)
}

// GetSecret mocks base method.
func (m *MockSecret) GetSecret(arg0 interface{}, arg1, arg2, arg3 string) (*v1.Secret, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetSecret", arg0, arg1, arg2, arg3)
	ret0, _ := ret[0].(*v1.Secret)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// GetSecret indicates an expected call of GetSecret.
func (mr *MockSecretMockRecorder) GetSecret(arg0, arg1, arg2, arg3 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetSecret", reflect.TypeOf((*MockSecret)(nil).GetSecret), arg0, arg1, arg2, arg3)
}

// ListSecret mocks base method.
func (m *MockSecret) ListSecret(arg0 string, arg1 *models.ListOptions) (*models.SecretList, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "ListSecret", arg0, arg1)
	ret0, _ := ret[0].(*models.SecretList)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// ListSecret indicates an expected call of ListSecret.
func (mr *MockSecretMockRecorder) ListSecret(arg0, arg1 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "ListSecret", reflect.TypeOf((*MockSecret)(nil).ListSecret), arg0, arg1)
}

// UpdateSecret mocks base method.
func (m *MockSecret) UpdateSecret(arg0 string, arg1 *v1.Secret) (*v1.Secret, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "UpdateSecret", arg0, arg1)
	ret0, _ := ret[0].(*v1.Secret)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// UpdateSecret indicates an expected call of UpdateSecret.
func (mr *MockSecretMockRecorder) UpdateSecret(arg0, arg1 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "UpdateSecret", reflect.TypeOf((*MockSecret)(nil).UpdateSecret), arg0, arg1)
}
