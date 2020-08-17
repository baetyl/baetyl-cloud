// Code generated by MockGen. DO NOT EDIT.
// Source: github.com/baetyl/baetyl-cloud/v2/service (interfaces: RegisterService)

// Package service is a generated GoMock package.
package service

import (
	models "github.com/baetyl/baetyl-cloud/v2/models"
	gomock "github.com/golang/mock/gomock"
	reflect "reflect"
)

// MockRegisterService is a mock of RegisterService interface
type MockRegisterService struct {
	ctrl     *gomock.Controller
	recorder *MockRegisterServiceMockRecorder
}

// MockRegisterServiceMockRecorder is the mock recorder for MockRegisterService
type MockRegisterServiceMockRecorder struct {
	mock *MockRegisterService
}

// NewMockRegisterService creates a new mock instance
func NewMockRegisterService(ctrl *gomock.Controller) *MockRegisterService {
	mock := &MockRegisterService{ctrl: ctrl}
	mock.recorder = &MockRegisterServiceMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use
func (m *MockRegisterService) EXPECT() *MockRegisterServiceMockRecorder {
	return m.recorder
}

// CreateBatch mocks base method
func (m *MockRegisterService) CreateBatch(arg0 *models.Batch) (*models.Batch, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "CreateBatch", arg0)
	ret0, _ := ret[0].(*models.Batch)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// CreateBatch indicates an expected call of CreateBatch
func (mr *MockRegisterServiceMockRecorder) CreateBatch(arg0 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "CreateBatch", reflect.TypeOf((*MockRegisterService)(nil).CreateBatch), arg0)
}

// CreateRecord mocks base method
func (m *MockRegisterService) CreateRecord(arg0 *models.Record) (*models.Record, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "CreateRecord", arg0)
	ret0, _ := ret[0].(*models.Record)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// CreateRecord indicates an expected call of CreateRecord
func (mr *MockRegisterServiceMockRecorder) CreateRecord(arg0 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "CreateRecord", reflect.TypeOf((*MockRegisterService)(nil).CreateRecord), arg0)
}

// DeleteBatch mocks base method
func (m *MockRegisterService) DeleteBatch(arg0, arg1 string) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "DeleteBatch", arg0, arg1)
	ret0, _ := ret[0].(error)
	return ret0
}

// DeleteBatch indicates an expected call of DeleteBatch
func (mr *MockRegisterServiceMockRecorder) DeleteBatch(arg0, arg1 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "DeleteBatch", reflect.TypeOf((*MockRegisterService)(nil).DeleteBatch), arg0, arg1)
}

// DeleteRecord mocks base method
func (m *MockRegisterService) DeleteRecord(arg0, arg1, arg2 string) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "DeleteRecord", arg0, arg1, arg2)
	ret0, _ := ret[0].(error)
	return ret0
}

// DeleteRecord indicates an expected call of DeleteRecord
func (mr *MockRegisterServiceMockRecorder) DeleteRecord(arg0, arg1, arg2 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "DeleteRecord", reflect.TypeOf((*MockRegisterService)(nil).DeleteRecord), arg0, arg1, arg2)
}

// DownloadRecords mocks base method
func (m *MockRegisterService) DownloadRecords(arg0, arg1 string) ([]byte, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "DownloadRecords", arg0, arg1)
	ret0, _ := ret[0].([]byte)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// DownloadRecords indicates an expected call of DownloadRecords
func (mr *MockRegisterServiceMockRecorder) DownloadRecords(arg0, arg1 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "DownloadRecords", reflect.TypeOf((*MockRegisterService)(nil).DownloadRecords), arg0, arg1)
}

// GenRecordRandom mocks base method
func (m *MockRegisterService) GenRecordRandom(arg0, arg1 string, arg2 int) ([]string, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GenRecordRandom", arg0, arg1, arg2)
	ret0, _ := ret[0].([]string)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// GenRecordRandom indicates an expected call of GenRecordRandom
func (mr *MockRegisterServiceMockRecorder) GenRecordRandom(arg0, arg1, arg2 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GenRecordRandom", reflect.TypeOf((*MockRegisterService)(nil).GenRecordRandom), arg0, arg1, arg2)
}

// GetBatch mocks base method
func (m *MockRegisterService) GetBatch(arg0, arg1 string) (*models.Batch, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetBatch", arg0, arg1)
	ret0, _ := ret[0].(*models.Batch)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// GetBatch indicates an expected call of GetBatch
func (mr *MockRegisterServiceMockRecorder) GetBatch(arg0, arg1 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetBatch", reflect.TypeOf((*MockRegisterService)(nil).GetBatch), arg0, arg1)
}

// GetRecord mocks base method
func (m *MockRegisterService) GetRecord(arg0, arg1, arg2 string) (*models.Record, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetRecord", arg0, arg1, arg2)
	ret0, _ := ret[0].(*models.Record)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// GetRecord indicates an expected call of GetRecord
func (mr *MockRegisterServiceMockRecorder) GetRecord(arg0, arg1, arg2 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetRecord", reflect.TypeOf((*MockRegisterService)(nil).GetRecord), arg0, arg1, arg2)
}

// GetRecordByFingerprint mocks base method
func (m *MockRegisterService) GetRecordByFingerprint(arg0, arg1, arg2 string) (*models.Record, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetRecordByFingerprint", arg0, arg1, arg2)
	ret0, _ := ret[0].(*models.Record)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// GetRecordByFingerprint indicates an expected call of GetRecordByFingerprint
func (mr *MockRegisterServiceMockRecorder) GetRecordByFingerprint(arg0, arg1, arg2 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetRecordByFingerprint", reflect.TypeOf((*MockRegisterService)(nil).GetRecordByFingerprint), arg0, arg1, arg2)
}

// ListBatch mocks base method
func (m *MockRegisterService) ListBatch(arg0 string, arg1 *models.Filter) (*models.ListView, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "ListBatch", arg0, arg1)
	ret0, _ := ret[0].(*models.ListView)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// ListBatch indicates an expected call of ListBatch
func (mr *MockRegisterServiceMockRecorder) ListBatch(arg0, arg1 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "ListBatch", reflect.TypeOf((*MockRegisterService)(nil).ListBatch), arg0, arg1)
}

// ListRecord mocks base method
func (m *MockRegisterService) ListRecord(arg0, arg1 string, arg2 *models.Filter) (*models.ListView, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "ListRecord", arg0, arg1, arg2)
	ret0, _ := ret[0].(*models.ListView)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// ListRecord indicates an expected call of ListRecord
func (mr *MockRegisterServiceMockRecorder) ListRecord(arg0, arg1, arg2 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "ListRecord", reflect.TypeOf((*MockRegisterService)(nil).ListRecord), arg0, arg1, arg2)
}

// UpdateBatch mocks base method
func (m *MockRegisterService) UpdateBatch(arg0 *models.Batch) (*models.Batch, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "UpdateBatch", arg0)
	ret0, _ := ret[0].(*models.Batch)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// UpdateBatch indicates an expected call of UpdateBatch
func (mr *MockRegisterServiceMockRecorder) UpdateBatch(arg0 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "UpdateBatch", reflect.TypeOf((*MockRegisterService)(nil).UpdateBatch), arg0)
}

// UpdateRecord mocks base method
func (m *MockRegisterService) UpdateRecord(arg0 *models.Record) (*models.Record, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "UpdateRecord", arg0)
	ret0, _ := ret[0].(*models.Record)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// UpdateRecord indicates an expected call of UpdateRecord
func (mr *MockRegisterServiceMockRecorder) UpdateRecord(arg0 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "UpdateRecord", reflect.TypeOf((*MockRegisterService)(nil).UpdateRecord), arg0)
}

// GenRecordFromUpload mocks base method
func (m *MockRegisterService) GenRecordFromUpload(arg0, arg1 string, arg2 []string) ([]string, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GenRecordFromUpload", arg0, arg1, arg2)
	ret0, _ := ret[0].([]string)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// GenRecordFromUpload indicates an expected call of GenRecordFromUpload
func (mr *MockRegisterServiceMockRecorder) GenRecordFromUpload(arg0, arg1, arg2 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GenRecordFromUpload", reflect.TypeOf((*MockRegisterService)(nil).GenRecordFromUpload), arg0, arg1, arg2)
}