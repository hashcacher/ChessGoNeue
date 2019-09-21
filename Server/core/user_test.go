package core_test

import (
	"reflect"
	"testing"

	"github.com/golang/mock/gomock"
	"github.com/hashcacher/ChessGoNeue/Server/v2/core"
	"github.com/hashcacher/ChessGoNeue/Server/v2/core/mocks"
)

func TestCreateUserOK(t *testing.T) {
	mockCtrl := gomock.NewController(t)
	defer mockCtrl.Finish()

	expect := core.User{Username: "zac", ClientID: "test-id", ID: 0}

	mockUsers := mocks.NewMockUsers(mockCtrl)
	mockUsers.EXPECT().Store(core.User{Username: "zac"}).Return(expect, nil)
	interactor := core.NewUsersInteractor(mockUsers)

	newUser, err := interactor.Create(core.User{Username: "zac"})
	if err != nil {
		t.Errorf("unexpected error: %v", err)
	}

	if !reflect.DeepEqual(newUser, expect) {
		t.Errorf("got: %v, expected: %v", newUser, expect)
	}
}
