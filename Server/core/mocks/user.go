// Code generated by MockGen. DO NOT EDIT.
// Source: ./core/user.go

// Package mocks is a generated GoMock package.
package mocks

import (
	gomock "github.com/golang/mock/gomock"
	core "github.com/hashcacher/ChessGoNeue/Server/v2/core"
	reflect "reflect"
)

// MockUsers is a mock of Users interface
type MockUsers struct {
	ctrl     *gomock.Controller
	recorder *MockUsersMockRecorder
}

// MockUsersMockRecorder is the mock recorder for MockUsers
type MockUsersMockRecorder struct {
	mock *MockUsers
}

// NewMockUsers creates a new mock instance
func NewMockUsers(ctrl *gomock.Controller) *MockUsers {
	mock := &MockUsers{ctrl: ctrl}
	mock.recorder = &MockUsersMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use
func (m *MockUsers) EXPECT() *MockUsersMockRecorder {
	return m.recorder
}

// Store mocks base method
func (m *MockUsers) Store(arg0 core.User) (int, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Store", arg0)
	ret0, _ := ret[0].(int)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// Store indicates an expected call of Store
func (mr *MockUsersMockRecorder) Store(arg0 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Store", reflect.TypeOf((*MockUsers)(nil).Store), arg0)
}

// FindBySecret mocks base method
func (m *MockUsers) FindBySecret(secret string) (core.User, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "FindBySecret", secret)
	ret0, _ := ret[0].(core.User)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// FindBySecret indicates an expected call of FindBySecret
func (mr *MockUsersMockRecorder) FindBySecret(secret interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "FindBySecret", reflect.TypeOf((*MockUsers)(nil).FindBySecret), secret)
}

// FindByID mocks base method
func (m *MockUsers) FindByID(id int) (core.User, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "FindByID", id)
	ret0, _ := ret[0].(core.User)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// FindByID indicates an expected call of FindByID
func (mr *MockUsersMockRecorder) FindByID(id interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "FindByID", reflect.TypeOf((*MockUsers)(nil).FindByID), id)
}

// Update mocks base method
func (m *MockUsers) Update(arg0 core.User) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Update", arg0)
	ret0, _ := ret[0].(error)
	return ret0
}

// Update indicates an expected call of Update
func (mr *MockUsersMockRecorder) Update(arg0 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Update", reflect.TypeOf((*MockUsers)(nil).Update), arg0)
}
