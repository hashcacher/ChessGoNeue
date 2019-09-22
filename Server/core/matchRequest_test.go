package core_test

import (
	"errors"
	"testing"

	"github.com/golang/mock/gomock"
	"github.com/hashcacher/ChessGoNeue/Server/v2/core"
	"github.com/hashcacher/ChessGoNeue/Server/v2/core/mocks"
)

// TODO: Error on store

// TestMatchMeSuccessWithNoInitialMatch without an initial match, it will create a
// match request then call wait for a match to be created
func TestMatchMeSuccessWithNoInitialMatch(t *testing.T) {
	mockCtrl := gomock.NewController(t)
	defer mockCtrl.Finish()

	// Create mock objects
	mockClientID := "mock-clientid-1"
	mockUser := core.User{
		ID:       1,
		ClientID: mockClientID,
	}
	expectStoreMatchRequest := core.MatchRequest{User: mockUser.ID}
	// mockMatchRequest := core.MatchRequest{ID: 2, User: mockUser.ID}
	mockGame := core.Game{ID: 3}

	// Create mocks
	mockUsers := mocks.NewMockUsers(mockCtrl)
	mockUsers.EXPECT().FindByClientID(mockClientID).Return(mockUser, nil)
	mockGames := mocks.NewMockGames(mockCtrl)
	mockMatchRequests := mocks.NewMockMatchRequests(mockCtrl)
	mockMatchRequests.EXPECT().FindMatchForUser(mockUser.ID).Return(core.MatchRequest{})
	mockMatchRequests.EXPECT().Store(expectStoreMatchRequest)
	mockMatchRequests.EXPECT().ListenForGameCreatedNotify(mockUser.ID).Return(mockGame.ID)

	// Create interactor and inject mocks
	interactor := core.NewMatchRequestsInteractor(mockMatchRequests, mockUsers, mockGames)

	gotGameId, err := interactor.MatchMe(mockClientID)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if gotGameId != mockGame.ID {
		t.Fatalf("got: %v, expected: %v", gotGameId, mockGame.ID)
	}
}

// TestMatchMeSuccessWithInitialMatch with an initial match, it will delete the
// match request, create a game and return the game id
func TestMatchMeSuccessWithInitialMatch(t *testing.T) {
	mockCtrl := gomock.NewController(t)
	defer mockCtrl.Finish()

	// Mock users (mockUser is the user making the call, mockOtherUser is the one we are matched with)
	mockClientID := "mock-clientid-1"
	mockClientID2 := "mock-clientid-2"
	mockUser := core.User{
		ID:       1,
		ClientID: mockClientID,
	}
	mockOtherUser := core.User{
		ID:       2,
		ClientID: mockClientID2,
	}
	// Match request (this assumes otherUser had already called matchMe and created a request)
	mockMatchRequest := core.MatchRequest{ID: 3, User: mockOtherUser.ID}
	// Mock values for the game object
	mockGameID := 4
	mockGame := core.Game{ID: 4, WhiteUser: mockUser.ID, BlackUser: mockOtherUser.ID}

	// Create mocks
	mockUsers := mocks.NewMockUsers(mockCtrl)
	mockUsers.EXPECT().FindByClientID(mockClientID).Return(mockUser, nil)
	mockGames := mocks.NewMockGames(mockCtrl)
	mockGames.EXPECT().Store(core.Game{WhiteUser: mockUser.ID, BlackUser: mockOtherUser.ID}).Return(mockGameID, nil)
	mockMatchRequests := mocks.NewMockMatchRequests(mockCtrl)
	mockMatchRequests.EXPECT().FindMatchForUser(mockUser.ID).Return(mockMatchRequest)
	mockMatchRequests.EXPECT().Delete(mockMatchRequest.ID).Return(1, nil)

	// Create interactor and inject mocks
	interactor := core.NewMatchRequestsInteractor(mockMatchRequests, mockUsers, mockGames)

	gotGameId, err := interactor.MatchMe(mockClientID)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if gotGameId != mockGame.ID {
		t.Fatalf("got: %v, expected: %v", gotGameId, mockGame.ID)
	}
}

// TestMatchMeBackendError test on error with the back end
func TestMatchMeBackendError(t *testing.T) {
	mockCtrl := gomock.NewController(t)
	defer mockCtrl.Finish()

	expectErr := errors.New("error fetching")
	mockClientID := "mock-clientid-1"

	// Create mocks
	mockUsers := mocks.NewMockUsers(mockCtrl)
	mockUsers.EXPECT().FindByClientID(mockClientID).Return(core.User{}, expectErr)
	mockGames := mocks.NewMockGames(mockCtrl)
	mockMatchRequests := mocks.NewMockMatchRequests(mockCtrl)

	// Create interactor and inject mocks
	interactor := core.NewMatchRequestsInteractor(mockMatchRequests, mockUsers, mockGames)

	_, err := interactor.MatchMe(mockClientID)
	if err != nil {
		if err.Error() != expectErr.Error() {
			t.Fatalf("got error: %v, expected error: %v", err, expectErr)
		}
	} else {
		t.Fatalf("Expected error, but call was succesful")
	}
}

// TestMatchMeUserDNE error if user by clientID comes back empty (doesn't exist)
func TestMatchMeUserDNE(t *testing.T) {
	mockCtrl := gomock.NewController(t)
	defer mockCtrl.Finish()

	expectErr := errors.New("could not find user with that client id")
	mockClientID := "mock-clientid-1"

	// Create mocks
	mockUsers := mocks.NewMockUsers(mockCtrl)
	mockUsers.EXPECT().FindByClientID(mockClientID).Return(core.User{}, nil)
	mockGames := mocks.NewMockGames(mockCtrl)
	mockMatchRequests := mocks.NewMockMatchRequests(mockCtrl)

	// Create interactor and inject mocks
	interactor := core.NewMatchRequestsInteractor(mockMatchRequests, mockUsers, mockGames)

	_, err := interactor.MatchMe(mockClientID)
	if err != nil {
		if err.Error() != expectErr.Error() {
			t.Fatalf("got error: %v, expected error: %v", err, expectErr)
		}
	} else {
		t.Fatalf("Expected error, but call was succesful")
	}
}
