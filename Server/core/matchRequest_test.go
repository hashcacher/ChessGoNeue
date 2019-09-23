package core_test

import (
	"errors"
	"testing"

	"github.com/golang/mock/gomock"
	"github.com/hashcacher/ChessGoNeue/Server/v2/core"
	"github.com/hashcacher/ChessGoNeue/Server/v2/core/mocks"
)

// TODO: Error on store

// TestMatchMeErrorFindingUser error if error finding user
func TestMatchMeErrorFindingUser(t *testing.T) {
	mockCtrl := gomock.NewController(t)
	defer mockCtrl.Finish()

	expectErr := errors.New("db error")
	mockSecret := "mock-Secret-1"

	// Create mocks
	mockUsers := mocks.NewMockUsers(mockCtrl)
	mockUsers.EXPECT().FindBySecret(mockSecret).Return(core.User{}, expectErr)
	mockGames := mocks.NewMockGames(mockCtrl)
	mockMatchRequests := mocks.NewMockMatchRequests(mockCtrl)

	// Create interactor and inject mocks
	interactor := core.NewMatchRequestsInteractor(mockMatchRequests, mockUsers, mockGames)

	_, err := interactor.MatchMe(mockSecret)
	if err != nil {
		if err.Error() != expectErr.Error() {
			t.Fatalf("got error: %v, expected error: %v", err, expectErr)
		}
	} else {
		t.Fatalf("Expected error, but call was succesful")
	}
}

// TestMatchMeUserDNE error if user by Secret comes back empty (doesn't exist)
func TestMatchMeUserDNE(t *testing.T) {
	mockCtrl := gomock.NewController(t)
	defer mockCtrl.Finish()

	expectErr := errors.New("could not find user with that client id")
	mockSecret := "mock-Secret-1"

	// Create mocks
	mockUsers := mocks.NewMockUsers(mockCtrl)
	mockUsers.EXPECT().FindBySecret(mockSecret).Return(core.User{}, nil)
	mockGames := mocks.NewMockGames(mockCtrl)
	mockMatchRequests := mocks.NewMockMatchRequests(mockCtrl)

	// Create interactor and inject mocks
	interactor := core.NewMatchRequestsInteractor(mockMatchRequests, mockUsers, mockGames)

	_, err := interactor.MatchMe(mockSecret)
	if err != nil {
		if err.Error() != expectErr.Error() {
			t.Fatalf("got error: %v, expected error: %v", err, expectErr)
		}
	} else {
		t.Fatalf("Expected error, but call was succesful")
	}
}

// TestFindMatchRequestByUserIDError error if we encounter an error trying to check if a MatchRequest
// already exists for a user
func TestFindMatchRequestByUserIDError(t *testing.T) {
	mockCtrl := gomock.NewController(t)
	defer mockCtrl.Finish()

	expectErr := errors.New("db error")
	mockSecret := "mock-secret-1"
	mockUser := core.User{ID: 1, Secret: mockSecret}

	// Create mocks
	mockUsers := mocks.NewMockUsers(mockCtrl)
	mockUsers.EXPECT().FindBySecret(mockSecret).Return(mockUser, nil)
	mockGames := mocks.NewMockGames(mockCtrl)
	mockMatchRequests := mocks.NewMockMatchRequests(mockCtrl)
	mockMatchRequests.EXPECT().FindMatchRequestByUserID(mockUser.ID).Return(core.MatchRequest{}, expectErr)

	// Create interactor and inject mocks
	interactor := core.NewMatchRequestsInteractor(mockMatchRequests, mockUsers, mockGames)

	_, err := interactor.MatchMe(mockSecret)
	if err != nil {
		if err.Error() != expectErr.Error() {
			t.Fatalf("got error: %v, expected error: %v", err, expectErr)
		}
	} else {
		t.Fatalf("Expected error, but call was succesful")
	}
}

// TestFindMatchRequestByUserIDExists error if we find a request already exists
func TestFindMatchRequestByUserIDExists(t *testing.T) {
	mockCtrl := gomock.NewController(t)
	defer mockCtrl.Finish()

	expectErr := errors.New("you can only queue for one game at a time")
	mockSecret := "mock-secret-1"
	mockUser := core.User{ID: 1, Secret: mockSecret}
	mockMatchRequest := core.MatchRequest{ID: 2, User: mockUser.ID}

	// Create mocks
	mockUsers := mocks.NewMockUsers(mockCtrl)
	mockUsers.EXPECT().FindBySecret(mockSecret).Return(mockUser, nil)
	mockGames := mocks.NewMockGames(mockCtrl)
	mockMatchRequests := mocks.NewMockMatchRequests(mockCtrl)
	mockMatchRequests.EXPECT().FindMatchRequestByUserID(mockUser.ID).Return(mockMatchRequest, nil)

	// Create interactor and inject mocks
	interactor := core.NewMatchRequestsInteractor(mockMatchRequests, mockUsers, mockGames)

	_, err := interactor.MatchMe(mockSecret)
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
	expectStoreMatchRequest := core.MatchRequest{User: mockUser.ID}
	mockGame := core.Game{ID: 3}

	// Create mocks
	mockUsers := mocks.NewMockUsers(mockCtrl)
	mockUsers.EXPECT().FindBySecret(mockSecret).Return(mockUser, nil)
	mockGames := mocks.NewMockGames(mockCtrl)
	mockGames.EXPECT().ListenForGameCreatedNotification(mockUser.ID).Return(mockGame.ID)
	mockMatchRequests := mocks.NewMockMatchRequests(mockCtrl)
	mockMatchRequests.EXPECT().Store(expectStoreMatchRequest)
	mockMatchRequests.EXPECT().FindMatchRequestByUserID(mockUser.ID).Return(core.MatchRequest{}, nil)

	// Create interactor and inject mocks
	interactor := core.NewMatchRequestsInteractor(mockMatchRequests, mockUsers, mockGames)

	gotGameId, err := interactor.MatchMe(mockSecret)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if gotGameId != mockGame.ID {
		t.Fatalf("got: %v, expected: %v", gotGameId, mockGame.ID)
	}
}
