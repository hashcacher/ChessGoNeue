package core

import (
	"errors"
	"reflect"
)

// Game stores chessgo game state such as the current board and who is in the game
type Game struct {
	ID        int        `json:"id"`
	BlackUser int        `json:"blackUser"`
	WhiteUser int        `json:"WhiteUser"`
	Board     [8][8]byte `json:"board"`
}

// Games is the use case for Game entitiy
type Games interface {
	Store(Game) (id int, err error)
	FindById(id int) (Game, error)
	Update(Game) error
	// Block and listen for a notification saying a game was created for the specified user
	ListenForGameCreatedNotification(userID int) (gameID int)
	// Notify the specified user that a game was created for them
	NotifyGameCreated(userID, gameID int) error
}

// GamesInteractor is a struct that holds data to be injected for use cases
type GamesInteractor struct {
	games Games
	users Users
}

// NewGamesInteractor generates a new GamesInteractor from the given Users store
func NewGamesInteractor(games Games, users Users) GamesInteractor {
	return GamesInteractor{
		games,
		users,
	}
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

	// Notify users
	i.games.NotifyGameCreated(game.WhiteUser, id)
	i.games.NotifyGameCreated(game.BlackUser, id)

	return id, nil
}

// // ExecuteMove validates a user and then performs a move
// func (i *GamesInteractor) ExecuteMove(m string, userID, gameId int) {
// 	// (UserRepository) Validate user is in match and it is their turn
// 	// (GameRepository) Perform update
// 	// (-)              Notify other user about the update
// }
