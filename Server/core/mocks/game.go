// Code generated by MockGen. DO NOT EDIT.
// Source: ./core/game.go

// Package mocks is a generated GoMock package.
package mocks

import (
	gomock "github.com/golang/mock/gomock"
	core "github.com/hashcacher/ChessGoNeue/Server/v2/core"
	reflect "reflect"
)

// MockGames is a mock of Games interface
type MockGames struct {
	ctrl     *gomock.Controller
	recorder *MockGamesMockRecorder
}

// MockGamesMockRecorder is the mock recorder for MockGames
type MockGamesMockRecorder struct {
	mock *MockGames
}

// NewMockGames creates a new mock instance
func NewMockGames(ctrl *gomock.Controller) *MockGames {
	mock := &MockGames{ctrl: ctrl}
	mock.recorder = &MockGamesMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use
func (m *MockGames) EXPECT() *MockGamesMockRecorder {
	return m.recorder
}

// MakeMove mocks base method
func (m *MockGames) MakeMove(arg0 *core.Game, arg1 *core.User, arg2 string) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "MakeMove", arg0, arg1, arg2)
	ret0, _ := ret[0].(error)
	return ret0
}

// MakeMove indicates an expected call of MakeMove
func (mr *MockGamesMockRecorder) MakeMove(arg0, arg1, arg2 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "MakeMove", reflect.TypeOf((*MockGames)(nil).MakeMove), arg0, arg1, arg2)
}

// GetMove mocks base method
func (m *MockGames) GetMove(arg0 *core.Game, arg1 *core.User) (string, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetMove", arg0, arg1)
	ret0, _ := ret[0].(string)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// GetMove indicates an expected call of GetMove
func (mr *MockGamesMockRecorder) GetMove(arg0, arg1 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetMove", reflect.TypeOf((*MockGames)(nil).GetMove), arg0, arg1)
}

// GetBoard mocks base method
func (m *MockGames) GetBoard(arg0 *core.Game) [8][8]byte {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetBoard", arg0)
	ret0, _ := ret[0].([8][8]byte)
	return ret0
}

// GetBoard indicates an expected call of GetBoard
func (mr *MockGamesMockRecorder) GetBoard(arg0 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetBoard", reflect.TypeOf((*MockGames)(nil).GetBoard), arg0)
}

// Store mocks base method
func (m *MockGames) Store(arg0 *core.Game) (int, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Store", arg0)
	ret0, _ := ret[0].(int)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// Store indicates an expected call of Store
func (mr *MockGamesMockRecorder) Store(arg0 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Store", reflect.TypeOf((*MockGames)(nil).Store), arg0)
}

// ListenForStoreByUserID mocks base method
func (m *MockGames) ListenForStoreByUserID(userID int) (*core.Game, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "ListenForStoreByUserID", userID)
	ret0, _ := ret[0].(*core.Game)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// ListenForStoreByUserID indicates an expected call of ListenForStoreByUserID
func (mr *MockGamesMockRecorder) ListenForStoreByUserID(userID interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "ListenForStoreByUserID", reflect.TypeOf((*MockGames)(nil).ListenForStoreByUserID), userID)
}

// ListenForMoveByUserID mocks base method
func (m *MockGames) ListenForMoveByUserID(userID int) (string, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "ListenForMoveByUserID", userID)
	ret0, _ := ret[0].(string)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// ListenForMoveByUserID indicates an expected call of ListenForMoveByUserID
func (mr *MockGamesMockRecorder) ListenForMoveByUserID(userID interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "ListenForMoveByUserID", reflect.TypeOf((*MockGames)(nil).ListenForMoveByUserID), userID)
}

// FindById mocks base method
func (m *MockGames) FindById(id int) (core.Game, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "FindById", id)
	ret0, _ := ret[0].(core.Game)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// FindById indicates an expected call of FindById
func (mr *MockGamesMockRecorder) FindById(id interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "FindById", reflect.TypeOf((*MockGames)(nil).FindById), id)
}

// FindByUserId mocks base method
func (m *MockGames) FindByUserId(id int) ([]*core.Game, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "FindByUserId", id)
	ret0, _ := ret[0].([]*core.Game)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// FindByUserId indicates an expected call of FindByUserId
func (mr *MockGamesMockRecorder) FindByUserId(id interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "FindByUserId", reflect.TypeOf((*MockGames)(nil).FindByUserId), id)
}

// Update mocks base method
func (m *MockGames) Update(arg0 core.Game) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Update", arg0)
	ret0, _ := ret[0].(error)
	return ret0
}

// Update indicates an expected call of Update
func (mr *MockGamesMockRecorder) Update(arg0 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Update", reflect.TypeOf((*MockGames)(nil).Update), arg0)
}
