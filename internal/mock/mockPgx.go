// Code generated by MockGen. DO NOT EDIT.
// Source: github.com/Azcarot/PasswordStorage/internal/storage (interfaces: PgxStorage)

// Package mock_storage is a generated GoMock package.
package mock_storage

import (
	context "context"
	reflect "reflect"

	gomock "github.com/golang/mock/gomock"
)

// MockPgxStorage is a mock of PgxStorage interface.
type MockPgxStorage struct {
	ctrl     *gomock.Controller
	recorder *MockPgxStorageMockRecorder
}

// MockPgxStorageMockRecorder is the mock recorder for MockPgxStorage.
type MockPgxStorageMockRecorder struct {
	mock *MockPgxStorage
}

// NewMockPgxStorage creates a new mock instance.
func NewMockPgxStorage(ctrl *gomock.Controller) *MockPgxStorage {
	mock := &MockPgxStorage{ctrl: ctrl}
	mock.recorder = &MockPgxStorageMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use.
func (m *MockPgxStorage) EXPECT() *MockPgxStorageMockRecorder {
	return m.recorder
}

// AddData mocks base method.
func (m *MockPgxStorage) AddData(arg0 interface{}) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "AddData", arg0)
	ret0, _ := ret[0].(error)
	return ret0
}

// AddData indicates an expected call of AddData.
func (mr *MockPgxStorageMockRecorder) AddData(arg0 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "AddData", reflect.TypeOf((*MockPgxStorage)(nil).AddData), arg0)
}

// CreateNewRecord mocks base method.
func (m *MockPgxStorage) CreateNewRecord(arg0 context.Context) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "CreateNewRecord", arg0)
	ret0, _ := ret[0].(error)
	return ret0
}

// CreateNewRecord indicates an expected call of CreateNewRecord.
func (mr *MockPgxStorageMockRecorder) CreateNewRecord(arg0 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "CreateNewRecord", reflect.TypeOf((*MockPgxStorage)(nil).CreateNewRecord), arg0)
}

// DeleteRecord mocks base method.
func (m *MockPgxStorage) DeleteRecord(arg0 context.Context) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "DeleteRecord", arg0)
	ret0, _ := ret[0].(error)
	return ret0
}

// DeleteRecord indicates an expected call of DeleteRecord.
func (mr *MockPgxStorageMockRecorder) DeleteRecord(arg0 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "DeleteRecord", reflect.TypeOf((*MockPgxStorage)(nil).DeleteRecord), arg0)
}

// GetAllRecords mocks base method.
func (m *MockPgxStorage) GetAllRecords(arg0 context.Context) (interface{}, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetAllRecords", arg0)
	ret0, _ := ret[0].(interface{})
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// GetAllRecords indicates an expected call of GetAllRecords.
func (mr *MockPgxStorageMockRecorder) GetAllRecords(arg0 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetAllRecords", reflect.TypeOf((*MockPgxStorage)(nil).GetAllRecords), arg0)
}

// GetData mocks base method.
func (m *MockPgxStorage) GetData() interface{} {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetData")
	ret0, _ := ret[0].(interface{})
	return ret0
}

// GetData indicates an expected call of GetData.
func (mr *MockPgxStorageMockRecorder) GetData() *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetData", reflect.TypeOf((*MockPgxStorage)(nil).GetData))
}

// GetRecord mocks base method.
func (m *MockPgxStorage) GetRecord(arg0 context.Context) (interface{}, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetRecord", arg0)
	ret0, _ := ret[0].(interface{})
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// GetRecord indicates an expected call of GetRecord.
func (mr *MockPgxStorageMockRecorder) GetRecord(arg0 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetRecord", reflect.TypeOf((*MockPgxStorage)(nil).GetRecord), arg0)
}

// HashDatabaseData mocks base method.
func (m *MockPgxStorage) HashDatabaseData(arg0 context.Context) (string, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "HashDatabaseData", arg0)
	ret0, _ := ret[0].(string)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// HashDatabaseData indicates an expected call of HashDatabaseData.
func (mr *MockPgxStorageMockRecorder) HashDatabaseData(arg0 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "HashDatabaseData", reflect.TypeOf((*MockPgxStorage)(nil).HashDatabaseData), arg0)
}

// SearchRecord mocks base method.
func (m *MockPgxStorage) SearchRecord(arg0 context.Context) (interface{}, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "SearchRecord", arg0)
	ret0, _ := ret[0].(interface{})
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// SearchRecord indicates an expected call of SearchRecord.
func (mr *MockPgxStorageMockRecorder) SearchRecord(arg0 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "SearchRecord", reflect.TypeOf((*MockPgxStorage)(nil).SearchRecord), arg0)
}

// UpdateRecord mocks base method.
func (m *MockPgxStorage) UpdateRecord(arg0 context.Context) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "UpdateRecord", arg0)
	ret0, _ := ret[0].(error)
	return ret0
}

// UpdateRecord indicates an expected call of UpdateRecord.
func (mr *MockPgxStorageMockRecorder) UpdateRecord(arg0 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "UpdateRecord", reflect.TypeOf((*MockPgxStorage)(nil).UpdateRecord), arg0)
}