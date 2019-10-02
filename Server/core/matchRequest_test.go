package core_test

import (
	"errors"
	"testing"

	"github.com/golang/mock/gomock"
	"github.com/hashcacher/ChessGoNeue/Server/v2/core"
	"github.com/hashcacher/ChessGoNeue/Server/v2/core/mocks"
)

// TODO: Error on store

// TestFindByUserIDError error if we encounter an error trying to check if a MatchRequest
// already exists for a user
func TestFindByUserIDError(t *testing.T) {
	mockCtrl := gomock.NewController(t)
	defer mockCtrl.Finish()

	expectErr := errors.New("db error")
	mockSecret := "mock-secret-1"
	mockUser := core.User{ID: 1, Secret: mockSecret}

	// Create mocks
	mockUsers := mocks.NewMockUsers(mockCtrl)
	mockGames := mocks.NewMockGames(mockCtrl)
	mockMatchRequests := mocks.NewMockMatchRequests(mockCtrl)
	mockMatchRequests.EXPECT().FindByUserID(mockUser.ID).Return(core.MatchRequest{}, expectErr)

	// Create interactor and inject mocks
	interactor := core.NewMatchRequestsInteractor(mockMatchRequests, mockUsers, mockGames)

	_, err := interactor.MatchMe(mockUser.ID)
	if err != nil {
		if err.Error() != expectErr.Error() {
			t.Fatalf("got error: %v, expected error: %v", err, expectErr)
		}
	} else {
		t.Fatalf("Expected error, but call was succesful")
	}
}

// TestFindByUserIDExists error if we find a request already exists
func TestFindByUserIDExists(t *testing.T) {
	mockCtrl := gomock.NewController(t)
	defer mockCtrl.Finish()

	expectErr := errors.New("you can only queue for one game at a time")
	mockSecret := "mock-secret-1"
	mockUser := core.User{ID: 1, Secret: mockSecret}
	mockMatchRequest := core.MatchRequest{ID: 2, UserID: mockUser.ID}

	// Create mocks
	mockUsers := mocks.NewMockUsers(mockCtrl)
	mockGames := mocks.NewMockGames(mockCtrl)
	mockMatchRequests := mocks.NewMockMatchRequests(mockCtrl)
	mockMatchRequests.EXPECT().FindByUserID(mockUser.ID).Return(mockMatchRequest, nil)

	// Create interactor and inject mocks
	interactor := core.NewMatchRequestsInteractor(mockMatchRequests, mockUsers, mockGames)

	_, err := interactor.MatchMe(mockUser.ID)
	if err != nil {
		if err.Error() != expectErr.Error() {
			t.Fatalf("got error: %v, expected error: %v", err, expectErr)
		}
	} else {
		t.Fatalf("Expected error, but call was succesful")
	}
}

// TestMatchMeSuccess will succesfully create a match and call for a listener
func TestMatchMeSuccess(t *testing.T) {
	mockCtrl := gomock.NewController(t)
	defer mockCtrl.Finish()

	// Create mock objects
	mockSecret := "mock-Secret-1"
	mockUser := core.User{
		ID:     1,
		Secret: mockSecret,
	}
	expectStoreMatchRequest := core.MatchRequest{UserID: mockUser.ID}
	mockGame := core.Game{ID: 3}

	// Create mocks
	mockUsers := mocks.NewMockUsers(mockCtrl)
	mockGames := mocks.NewMockGames(mockCtrl)
	mockGames.EXPECT().ListenForStoreByUserID(mockUser.ID).Return(mockGame, nil)
	mockMatchRequests := mocks.NewMockMatchRequests(mockCtrl)
	mockMatchRequests.EXPECT().Store(expectStoreMatchRequest)
	mockMatchRequests.EXPECT().FindByUserID(mockUser.ID).Return(core.MatchRequest{}, nil)

	// Create interactor and inject mocks
	interactor := core.NewMatchRequestsInteractor(mockMatchRequests, mockUsers, mockGames)

	gotGame, err := interactor.MatchMe(mockUser.ID)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if gotGame.ID != mockGame.ID {
		t.Fatalf("got: %v, expected: %v", gotGame.ID, mockGame.ID)
	}
}
