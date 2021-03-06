// Code generated by MockGen. DO NOT EDIT.
// Source: ./script.go

// Package mock_repository is a generated GoMock package.
package mock_repository

import (
	reflect "reflect"

	gomock "github.com/golang/mock/gomock"
	dto "github.com/scriptscat/cloudcat/internal/service/sync/domain/dto"
	entity "github.com/scriptscat/cloudcat/internal/service/sync/domain/entity"
)

// MockScript is a mock of Script interface.
type MockScript struct {
	ctrl     *gomock.Controller
	recorder *MockScriptMockRecorder
}

// MockScriptMockRecorder is the mock recorder for MockScript.
type MockScriptMockRecorder struct {
	mock *MockScript
}

// NewMockScript creates a new mock instance.
func NewMockScript(ctrl *gomock.Controller) *MockScript {
	mock := &MockScript{ctrl: ctrl}
	mock.recorder = &MockScriptMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use.
func (m *MockScript) EXPECT() *MockScriptMockRecorder {
	return m.recorder
}

// ActionList mocks base method.
func (m *MockScript) ActionList(user, device, version int64) ([][]*dto.SyncScript, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "ActionList", user, device, version)
	ret0, _ := ret[0].([][]*dto.SyncScript)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// ActionList indicates an expected call of ActionList.
func (mr *MockScriptMockRecorder) ActionList(user, device, version interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "ActionList", reflect.TypeOf((*MockScript)(nil).ActionList), user, device, version)
}

// FindByUUID mocks base method.
func (m *MockScript) FindByUUID(user, device int64, uuid string) (*entity.SyncScript, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "FindByUUID", user, device, uuid)
	ret0, _ := ret[0].(*entity.SyncScript)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// FindByUUID indicates an expected call of FindByUUID.
func (mr *MockScriptMockRecorder) FindByUUID(user, device, uuid interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "FindByUUID", reflect.TypeOf((*MockScript)(nil).FindByUUID), user, device, uuid)
}

// LatestVersion mocks base method.
func (m *MockScript) LatestVersion(user, device int64) (int64, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "LatestVersion", user, device)
	ret0, _ := ret[0].(int64)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// LatestVersion indicates an expected call of LatestVersion.
func (mr *MockScriptMockRecorder) LatestVersion(user, device interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "LatestVersion", reflect.TypeOf((*MockScript)(nil).LatestVersion), user, device)
}

// ListScript mocks base method.
func (m *MockScript) ListScript(user, device int64) ([]*entity.SyncScript, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "ListScript", user, device)
	ret0, _ := ret[0].([]*entity.SyncScript)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// ListScript indicates an expected call of ListScript.
func (mr *MockScriptMockRecorder) ListScript(user, device interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "ListScript", reflect.TypeOf((*MockScript)(nil).ListScript), user, device)
}

// PushVersion mocks base method.
func (m *MockScript) PushVersion(user, device int64, data []*dto.SyncScript) (int64, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "PushVersion", user, device, data)
	ret0, _ := ret[0].(int64)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// PushVersion indicates an expected call of PushVersion.
func (mr *MockScriptMockRecorder) PushVersion(user, device, data interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "PushVersion", reflect.TypeOf((*MockScript)(nil).PushVersion), user, device, data)
}

// Save mocks base method.
func (m *MockScript) Save(entity *entity.SyncScript) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Save", entity)
	ret0, _ := ret[0].(error)
	return ret0
}

// Save indicates an expected call of Save.
func (mr *MockScriptMockRecorder) Save(entity interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Save", reflect.TypeOf((*MockScript)(nil).Save), entity)
}

// SetStatus mocks base method.
func (m *MockScript) SetStatus(id int64, status int8) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "SetStatus", id, status)
	ret0, _ := ret[0].(error)
	return ret0
}

// SetStatus indicates an expected call of SetStatus.
func (mr *MockScriptMockRecorder) SetStatus(id, status interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "SetStatus", reflect.TypeOf((*MockScript)(nil).SetStatus), id, status)
}
