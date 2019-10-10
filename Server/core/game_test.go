package core_test

import (
	"errors"
	"testing"

	"github.com/golang/mock/gomock"
	"github.com/hashcacher/ChessGoNeue/Server/v2/core"
	"github.com/hashcacher/ChessGoNeue/Server/v2/core/mocks"
)

// TestCreateGameOK successful basic case
func TestCreateGameOK(t *testing.T) {
	mockCtrl := gomock.NewController(t)
	defer mockCtrl.Finish()

	mockGameID := 1
	mockGame := core.Game{WhiteUser: 1, BlackUser: 2}
	mockUser1 := core.User{ID: 1, Username: "zac", Secret: "mock-Secret-1"}
	mockUser2 := core.User{ID: 2, Username: "greg", Secret: "mock-Secret-2"}

	// Create mocks
	mockUsers := mocks.NewMockUsers(mockCtrl)
	mockUsers.EXPECT().FindByID(1).Return(mockUser1, nil)
	mockUsers.EXPECT().FindByID(2).Return(mockUser2, nil)
	mockGames := mocks.NewMockGames(mockCtrl)
	mockGames.EXPECT().Store(&mockGame).Return(mockGameID, nil)
	mockMatchRequests := mocks.NewMockMatchRequests(mockCtrl)

	// Create interactor and inject mocks
	interactor := core.NewGamesInteractor(mockGames, mockUsers, mockMatchRequests)

	gotId, err := interactor.Create(&mockGame)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if gotId != mockGameID {
		t.Fatalf("got: %v, expected: %v", gotId, mockGameID)
	}
}

// TestCreateGameOKResetBoard successful but clears board
func TestCreateGameOKResetBoard(t *testing.T) {
	mockCtrl := gomock.NewController(t)
	defer mockCtrl.Finish()

	mockGameID := 1
	board := core.DefaultBoard()
	board[0][0] = 1
	mockGameInCall := core.Game{WhiteUser: 1, BlackUser: 2, Board: board}
	mockGameExpectedToStore := core.Game{WhiteUser: 1, BlackUser: 2, Board: core.DefaultBoard(), WhiteTurn: true}
	mockUser1 := core.User{ID: 1, Username: "zac", Secret: "mock-Secret-1"}
	mockUser2 := core.User{ID: 2, Username: "greg", Secret: "mock-Secret-2"}

	// Create mocks
	mockUsers := mocks.NewMockUsers(mockCtrl)
	mockUsers.EXPECT().FindByID(1).Return(mockUser1, nil)
	mockUsers.EXPECT().FindByID(2).Return(mockUser2, nil)
	mockGames := mocks.NewMockGames(mockCtrl)
	mockGames.EXPECT().Store(&mockGameExpectedToStore).Return(mockGameID, nil)
	mockMatchRequests := mocks.NewMockMatchRequests(mockCtrl)

	// Create interactor and inject mocks
	interactor := core.NewGamesInteractor(mockGames, mockUsers, mockMatchRequests)

	gotId, err := interactor.Create(&mockGameInCall)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if gotId != mockGameID {
		t.Fatalf("got: %v, expected: %v", gotId, mockGameID)
	}

}

// TestGetMoves successful basic case
func TestMakeAndGetMoveOK(t *testing.T) {
	mockCtrl := gomock.NewController(t)
	defer mockCtrl.Finish()

	mockGameID := 1
	mockGame := core.Game{WhiteUser: 1, BlackUser: 2}
	mockUser1 := core.User{ID: 1, Username: "zac", Secret: "mock-Secret-1"}
	mockUser2 := core.User{ID: 2, Username: "greg", Secret: "mock-Secret-2"}
	mockMove := "1,1>2,1"

	// Create mocks
	mockUsers := mocks.NewMockUsers(mockCtrl)
	mockGames := mocks.NewMockGames(mockCtrl)
	mockUsers.EXPECT().FindByID(1).Return(mockUser1, nil)
	mockUsers.EXPECT().FindByID(2).Return(mockUser2, nil)
	mockGames.EXPECT().Store(&mockGame).Return(mockGameID, nil)
	mockMatchRequests := mocks.NewMockMatchRequests(mockCtrl)

	// Create interactor and inject mocks
	interactor := core.NewGamesInteractor(mockGames, mockUsers, mockMatchRequests)

	gotId, err := interactor.Create(&mockGame)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if gotId != mockGameID {
		t.Fatalf("got: %v, expected: %v", gotId, mockGameID)
	}

	mockGame.ID = gotId
	mockGame.WhiteTurn = true

	// Make and get moves
	mockUsers.EXPECT().FindBySecret(mockUser1.Secret).Return(mockUser1, nil)
	mockGames.EXPECT().FindByUserId(mockUser1.ID).Return([]*core.Game{&mockGame}, nil)
	mockGames.EXPECT().MakeMove(&mockGame, &mockUser1, mockMove).Return(nil)
	err = interactor.MakeMove(mockUser1.Secret, mockGameID, mockMove)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	mockUsers.EXPECT().FindBySecret(mockUser2.Secret).Return(mockUser2, nil)
	mockGames.EXPECT().FindByUserId(mockUser2.ID).Return([]*core.Game{&mockGame}, nil)
	mockGames.EXPECT().GetMove(&mockGame, &mockUser2).Return(mockMove, nil)
	move, err := interactor.GetMove(mockUser2.Secret, gotId)
	if move != mockMove {
		t.Fatalf("unexpected move %s", move)
	}
	if err != nil {
		t.Fatalf("unexpected error %v", err)
	}
}

// CreateGame Error if users are the same
// 	Expect no calls to mocks
func TestCreateGameErrorUsersSame(t *testing.T) {
	mockCtrl := gomock.NewController(t)
	defer mockCtrl.Finish()

	expectErr := errors.New("you cannot play a game with yourself")
	mockGame := core.Game{WhiteUser: 1, BlackUser: 1}

	// Create mocks
	mockUsers := mocks.NewMockUsers(mockCtrl)
	mockGames := mocks.NewMockGames(mockCtrl)
	mockMatchRequests := mocks.NewMockMatchRequests(mockCtrl)

	// Create interactor and inject mocks
	interactor := core.NewGamesInteractor(mockGames, mockUsers, mockMatchRequests)

	_, err := interactor.Create(&mockGame)
	if err != nil {
		if err.Error() != expectErr.Error() {
			t.Fatalf("got error: %v, expected error: %v", err, expectErr)
		}
	} else {
		t.Fatalf("Expected error, but call was succesful")
	}
}

// CreateGame error white doesn't exist
func TestCreateGameErrorWhiteNotFound(t *testing.T) {
	mockCtrl := gomock.NewController(t)
	defer mockCtrl.Finish()

	expectErr := errors.New("could not find white user by that id")
	mockGame := core.Game{WhiteUser: 1, BlackUser: 2}
	mockUser1 := core.User{}
	mockUser2 := core.User{ID: 2, Username: "greg", Secret: "mock-Secret-2"}

	// Create mocks
	mockUsers := mocks.NewMockUsers(mockCtrl)
	mockUsers.EXPECT().FindByID(1).Return(mockUser1, nil)
	mockUsers.EXPECT().FindByID(2).Return(mockUser2, nil)
	mockGames := mocks.NewMockGames(mockCtrl)
	mockMatchRequests := mocks.NewMockMatchRequests(mockCtrl)

	// Create interactor and inject mocks
	interactor := core.NewGamesInteractor(mockGames, mockUsers, mockMatchRequests)

	_, err := interactor.Create(&mockGame)
	if err != nil {
		if err.Error() != expectErr.Error() {
			t.Fatalf("got error: %v, expected error: %v", err, expectErr)
		}
	} else {
		t.Fatalf("Expected error, but call was succesful")
	}
}

// CreateGame error black doesn't exist
func TestCreateGameErrorBlackNotFound(t *testing.T) {
	mockCtrl := gomock.NewController(t)
	defer mockCtrl.Finish()

	expectErr := errors.New("could not find black user by that id")
	mockGame := core.Game{WhiteUser: 1, BlackUser: 2}
	mockUser1 := core.User{ID: 1, Username: "zac", Secret: "mock-Secret-1"}
	mockUser2 := core.User{}

	// Create mocks
	mockUsers := mocks.NewMockUsers(mockCtrl)
	mockUsers.EXPECT().FindByID(1).Return(mockUser1, nil)
	mockUsers.EXPECT().FindByID(2).Return(mockUser2, nil)
	mockGames := mocks.NewMockGames(mockCtrl)
	mockMatchRequests := mocks.NewMockMatchRequests(mockCtrl)

	// Create interactor and inject mocks
	interactor := core.NewGamesInteractor(mockGames, mockUsers, mockMatchRequests)

	_, err := interactor.Create(&mockGame)
	if err != nil {
		if err.Error() != expectErr.Error() {
			t.Fatalf("got error: %v, expected error: %v", err, expectErr)
		}
	} else {
		t.Fatalf("Expected error, but call was succesful")
	}
}

// CreateGame error if the store is unsuccesful for some reason
func TestCreateGameErrorStoreUnsuccesful(t *testing.T) {
	mockCtrl := gomock.NewController(t)
	defer mockCtrl.Finish()

	expectErr := errors.New("store error")
	mockGame := core.Game{WhiteUser: 1, BlackUser: 2}
	mockUser1 := core.User{ID: 1, Username: "zac", Secret: "mock-Secret-1"}
	mockUser2 := core.User{ID: 2, Username: "greg", Secret: "mock-Secret-2"}

	// Create mocks
	mockUsers := mocks.NewMockUsers(mockCtrl)
	mockUsers.EXPECT().FindByID(1).Return(mockUser1, nil)
	mockUsers.EXPECT().FindByID(2).Return(mockUser2, nil)
	mockGames := mocks.NewMockGames(mockCtrl)
	mockGames.EXPECT().Store(&mockGame).Return(0, expectErr)
	mockMatchRequests := mocks.NewMockMatchRequests(mockCtrl)

	// Create interactor and inject mocks
	interactor := core.NewGamesInteractor(mockGames, mockUsers, mockMatchRequests)

	_, err := interactor.Create(&mockGame)
	if err != nil {
		if err.Error() != expectErr.Error() {
			t.Fatalf("got error: %v, expected error: %v", err, expectErr)
		}
	} else {
		t.Fatalf("Expected error, but call was succesful")
	}
}
