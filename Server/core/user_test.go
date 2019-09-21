package core_test

import (
	"errors"
	"reflect"
	"testing"

	"github.com/golang/mock/gomock"
	"github.com/hashcacher/ChessGoNeue/Server/v2/core"
	"github.com/hashcacher/ChessGoNeue/Server/v2/core/mocks"
)

// Successful create user call
func TestCreateUserOK(t *testing.T) {
	mockCtrl := gomock.NewController(t)
	defer mockCtrl.Finish()

	expectId := 1

	mockUsers := mocks.NewMockUsers(mockCtrl)
	mockUsers.EXPECT().Store(core.User{Username: "zac"}).Return(1, nil)
	interactor := core.NewUsersInteractor(mockUsers)

	gotId, err := interactor.Create(core.User{Username: "zac"})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if gotId != expectId {
		t.Fatalf("got: %v, expected: %v", gotId, expectId)
	}
}

// Expect error if username is blank
func TestCreateUserError(t *testing.T) {
	mockCtrl := gomock.NewController(t)
	defer mockCtrl.Finish()

	expectErr := errors.New("username can't be empty")

	mockUsers := mocks.NewMockUsers(mockCtrl)
	interactor := core.NewUsersInteractor(mockUsers)

	_, err := interactor.Create(core.User{Username: ""})
	if err != nil {
		if err.Error() != expectErr.Error() {
			t.Fatalf("got error: %v, expected error: %v", err, expectErr)
		}
	} else {
		t.Fatalf("Expected error, but call was succesful")
	}
}

// Find user by clientID is succesful if user exists
func TestFindByClientIDOK(t *testing.T) {
	mockCtrl := gomock.NewController(t)
	defer mockCtrl.Finish()

	mockClientID := "mock-client-id"
	mockUser := core.User{ID: 1, ClientID: mockClientID, Username: "zac"}

	mockUsers := mocks.NewMockUsers(mockCtrl)
	mockUsers.EXPECT().FindByClientID(mockClientID).Return(mockUser, nil)
	interactor := core.NewUsersInteractor(mockUsers)

	got, err := interactor.FindByClientID(mockClientID)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if !reflect.DeepEqual(got, mockUser) {
		t.Fatalf("got: %v, expected: %v", got, mockUser)
	}
}

// Find user by clientID error if db throws error
func TestFindByClientIDError(t *testing.T) {
	mockCtrl := gomock.NewController(t)
	defer mockCtrl.Finish()

	mockClientID := "mock-client-id"
	expectErr := errors.New("mock database error")

	mockUsers := mocks.NewMockUsers(mockCtrl)
	mockUsers.EXPECT().FindByClientID(mockClientID).Return(core.User{}, expectErr)
	interactor := core.NewUsersInteractor(mockUsers)

	_, err := interactor.FindByClientID(mockClientID)
	if err != nil {
		if err.Error() != expectErr.Error() {
			t.Fatalf("got error: %v, expected error: %v", err, expectErr)
		}
	} else {
		t.Fatalf("Expected error, but call was succesful")
	}
}

// Find user by ID is succesful if user exists
func TestFindByIDOK(t *testing.T) {
	mockCtrl := gomock.NewController(t)
	defer mockCtrl.Finish()

	mockID := 1
	mockUser := core.User{ID: mockID, ClientID: "mock-client-id", Username: "zac"}

	mockUsers := mocks.NewMockUsers(mockCtrl)
	mockUsers.EXPECT().FindByID(mockID).Return(mockUser, nil)
	interactor := core.NewUsersInteractor(mockUsers)

	got, err := interactor.FindByID(mockID)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if !reflect.DeepEqual(got, mockUser) {
		t.Fatalf("got: %v, expected: %v", got, mockUser)
	}
}

// Find user by ID error if db throws error
func TestFindByIDError(t *testing.T) {
	mockCtrl := gomock.NewController(t)
	defer mockCtrl.Finish()

	mockID := 1
	expectErr := errors.New("mock database error")

	mockUsers := mocks.NewMockUsers(mockCtrl)
	mockUsers.EXPECT().FindByID(mockID).Return(core.User{}, expectErr)
	interactor := core.NewUsersInteractor(mockUsers)

	_, err := interactor.FindByID(mockID)
	if err != nil {
		if err.Error() != expectErr.Error() {
			t.Fatalf("got error: %v, expected error: %v", err, expectErr)
		}
	} else {
		t.Fatalf("Expected error, but call was succesful")
	}
}
