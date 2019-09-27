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
}

// Games is the use case for Game entitiy
type Games interface {
	Store(Game) (id int, err error)
	ListenForStoreByUserID(userID int) (Game, error)
	FindById(id int) (Game, error)
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
	// Start daemon that listens for match requests and creates games accordingly
	go i.startGameCreateDaemon()
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

func (i GamesInteractor) startGameCreateDaemon() {
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
		}
		i.games.Store(game)
		log.Printf("INFO: Created game: %v\n", game)
	}
}

// // ExecuteMove validates a user and then performs a move
// func (i *GamesInteractor) ExecuteMove(m string, userID, gameId int) {
// 	// (UserRepository) Validate user is in match and it is their turn
// 	// (GameRepository) Perform update
// 	// (-)              Notify other user about the update
// }
