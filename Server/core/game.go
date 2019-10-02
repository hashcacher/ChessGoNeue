package core

import (
	"errors"
	"log"
	"reflect"
)

// Game stores chessgo game state such as the current board and who is in the game
type Game struct {
	ID        int        `json:"id"`
	BlackUser int        `json:"blackUser"`
	WhiteUser int        `json:"WhiteUser"`
	Board     [8][8]byte `json:"board"`
	WhiteTurn bool       `json:"whiteTurn"`
}

// Games is the use case for Game entitiy
type Games interface {
	MakeMove(*Game, *User, string) error
	GetMove(*Game, *User) (string, error)
	GetBoard(*Game) (board [8][8]byte)
	Store(Game) (id int, err error)
	ListenForStoreByUserID(userID int) (*Game, error)
	ListenForMoveByUserID(userID int) (string, error)
	FindById(id int) (Game, error)
	FindByUserId(id int) ([]*Game, error)
	Update(Game) error
}

// GamesInteractor is a struct that holds data to be injected for use cases
type GamesInteractor struct {
	games         Games
	users         Users
	matchRequests MatchRequests
}

// NewGamesInteractor generates a new GamesInteractor from the given Users store
func NewGamesInteractor(games Games, users Users, matchRequests MatchRequests) GamesInteractor {
	i := GamesInteractor{
		games,
		users,
		matchRequests,
	}
	// Return the interractor
	return i
}

// Create validates an incoming game's data (users, board) and then stores it
func (i *GamesInteractor) Create(game Game) (id int, err error) {
	if game.WhiteUser == game.BlackUser {
		return 0, errors.New("you cannot play a game with yourself")
	}

	whiteUser, err := i.users.FindByID(game.WhiteUser)
	if err != nil {
		return 0, err
	}
	blackUser, err := i.users.FindByID(game.BlackUser)
	if err != nil {
		return 0, err
	}

	// Validate white user
	if reflect.DeepEqual(whiteUser, User{}) {
		return 0, errors.New("could not find white user by that id")
	}
	// Validate black user
	if reflect.DeepEqual(blackUser, User{}) {
		return 0, errors.New("could not find black user by that id")
	}

	// Clear the board
	game.Board = [8][8]byte{}

	// Store game
	id, err = i.games.Store(game)
	if err != nil {
		return 0, err
	}

	return id, nil
}

func (i GamesInteractor) GetBoard(secret string, gameID int) ([8][8]byte, error) {
	user, err := i.users.FindBySecret(secret)
	if err != nil {
		return [8][8]byte{}, errors.New("couldnt find user by that id")
	}

	games := i.getGamesForUser(user.ID)
	if len(games) == 0 {
		return [8][8]byte{}, errors.New("couldnt find game")
	}

	return games[0].Board, nil // TODO take gameID into account
}

func (i GamesInteractor) getGamesForUser(userID int) []*Game {
	games, err := i.games.FindByUserId(userID)
	if err != nil {
		return []*Game{}
	}

	return games
}

func (i GamesInteractor) StartGameCreateDaemon() {
	for {
		// Wait until a store happens
		i.matchRequests.ListenForStore()
		// Get all match requests
		matchRequests, err := i.matchRequests.FindAll()
		if err != nil {
			log.Printf("ERROR: %v\n", err)
			continue
		}
		if len(matchRequests) <= 1 {
			continue
		}
		// Delete the first two requests from the store
		_, err = i.matchRequests.Delete(matchRequests[0].ID)
		if err != nil {
			log.Printf("ERROR: %v\n", err)
		}
		_, err = i.matchRequests.Delete(matchRequests[1].ID)
		if err != nil {
			log.Printf("ERROR: %v\n", err)
		}
		// Use the first two requests to create a game
		game := Game{
			WhiteUser: matchRequests[0].UserID,
			BlackUser: matchRequests[1].UserID,
			Board:     defaultBoard(),
		}
		i.games.Store(game)
		log.Printf("INFO: Created game: %v\n", game)
	}
}

func defaultBoard() [8][8]byte {
	board := [8][8]byte{}
	board[0] = [8]byte{'r', 'n', 'b', 'q', 'k', 'b', 'n', 'r'}
	board[1] = [8]byte{'p', 'p', 'p', 'p', 'p', 'p', 'p', 'p'}

	for i := 2; i < 6; i++ {
		board[i] = [8]byte{' ', ' ', ' ', ' ', ' ', ' ', ' ', ' '}
	}

	board[6] = [8]byte{'P', 'P', 'P', 'P', 'P', 'P', 'P', 'P'}
	board[7] = [8]byte{'R', 'N', 'B', 'Q', 'K', 'B', 'N', 'R'}

	return board
}

// MakeMove validates a user and then performs a move
func (i *GamesInteractor) MakeMove(secret string, gameId int, move string) error {
	user, err := i.users.FindBySecret(secret)
	if err != nil {
		return errors.New("couldnt find user by that id")
	}

	games := i.getGamesForUser(user.ID)
	if len(games) == 0 {
		return errors.New("couldnt find game")
	}

	game := games[0] // TODO take gameID into account
	if game.WhiteTurn && game.WhiteUser != user.ID ||
		!game.WhiteTurn && game.BlackUser != user.ID {
		return errors.New("not your turn")
	}

	return i.games.MakeMove(games[0], &user, move)
}

func (i *GamesInteractor) GetMove(secret string, gameId int) (string, error) {
	user, err := i.users.FindBySecret(secret)
	if err != nil {
		return "", errors.New("couldnt find user by that id")
	}

	games := i.getGamesForUser(user.ID)
	if len(games) == 0 {
		return "", errors.New("couldnt find game")
	}

	game := games[0] // TODO take gameID into account

	return i.games.GetMove(game, &user)
}
