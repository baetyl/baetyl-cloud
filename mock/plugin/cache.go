// Code generated by MockGen. DO NOT EDIT.
// Source: ./plugin/cache.go

// Package plugin is a generated GoMock package.
package plugin

import (
	models "github.com/baetyl/baetyl-cloud/models"
	gomock "github.com/golang/mock/gomock"
	reflect "reflect"
)

// MockCacheStorage is a mock of CacheStorage interface.
type MockCacheStorage struct {
	ctrl     *gomock.Controller
	recorder *MockCacheStorageMockRecorder
}

// MockCacheStorageMockRecorder is the mock recorder for MockCacheStorage.
type MockCacheStorageMockRecorder struct {
	mock *MockCacheStorage
}

// NewMockCacheStorage creates a new mock instance.
func NewMockCacheStorage(ctrl *gomock.Controller) *MockCacheStorage {
	mock := &MockCacheStorage{ctrl: ctrl}
	mock.recorder = &MockCacheStorageMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use.
func (m *MockCacheStorage) EXPECT() *MockCacheStorageMockRecorder {
	return m.recorder
}

// GetCache mocks base method.
func (m *MockCacheStorage) GetCache(key string) (string, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetCache", key)
	ret0, _ := ret[0].(string)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// GetCache indicates an expected call of GetCache.
func (mr *MockCacheStorageMockRecorder) GetCache(key interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetCache", reflect.TypeOf((*MockCacheStorage)(nil).GetCache), key)
}

// SetCache mocks base method.
func (m *MockCacheStorage) SetCache(key, value string) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "SetCache", key, value)
	ret0, _ := ret[0].(error)
	return ret0
}

// SetCache indicates an expected call of SetCache.
func (mr *MockCacheStorageMockRecorder) SetCache(key, value interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "SetCache", reflect.TypeOf((*MockCacheStorage)(nil).SetCache), key, value)
}

// DeleteCache mocks base method.
func (m *MockCacheStorage) DeleteCache(key string) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "DeleteCache", key)
	ret0, _ := ret[0].(error)
	return ret0
}

// DeleteCache indicates an expected call of DeleteCache.
func (mr *MockCacheStorageMockRecorder) DeleteCache(key interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "DeleteCache", reflect.TypeOf((*MockCacheStorage)(nil).DeleteCache), key)
}

// ListCache mocks base method.
func (m *MockCacheStorage) ListCache(page *models.Filter) (*models.AmisListView, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "ListCache", page)
	ret0, _ := ret[0].(*models.AmisListView)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// ListCache indicates an expected call of ListCache.
func (mr *MockCacheStorageMockRecorder) ListCache(page interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "ListCache", reflect.TypeOf((*MockCacheStorage)(nil).ListCache), page)
}

// Close mocks base method.
func (m *MockCacheStorage) Close() error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Close")
	ret0, _ := ret[0].(error)
	return ret0
}

// Close indicates an expected call of Close.
func (mr *MockCacheStorageMockRecorder) Close() *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Close", reflect.TypeOf((*MockCacheStorage)(nil).Close))
}
