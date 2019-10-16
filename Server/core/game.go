package core

import (
	"errors"
	"fmt"
	"log"
	"math/rand"
	"reflect"
	"strconv"
	"time"
)

// Game stores chessgo game state such as the current board and who is in the game
type Game struct {
	ID               int           `json:"-"`
	BlackUser        int           `json:"-"`
	WhiteUser        int           `json:"-"`
	Board            [8][8]byte    `json:"board"`
	WhiteTurn        bool          `json:"whiteTurn"`
	BlackLeft        time.Duration `json:"blackLeft"`
	WhiteLeft        time.Duration `json:"whiteLeft"`
	BlackTurnStarted time.Time     `json:"blackTurnStarted"`
	WhiteTurnStarted time.Time     `json:"whiteTurnStarted"`
	Duration         time.Duration `json:"duration"`
}

// Games is the use case for Game entitiy
type Games interface {
	MakeMove(*Game, *User, string) error
	GetMove(*Game, *User) (string, error)
	GetBoard(*Game) (board [8][8]byte)
	Store(*Game) (id int, err error)
	ListenForStoreByUserID(userID int) (*Game, error)
	ListenForMoveByUserID(userID int) (string, error)
	FindById(id int) (Game, error)
	FindByUserId(id int) ([]*Game, error)
	Update(*Game) error
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
func (i *GamesInteractor) Create(game *Game) (id int, err error) {
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

	// White goes first
	game.WhiteTurn = true

	// Clear the board
	game.Board = DefaultBoard()

	// Time
	game.BlackLeft = game.Duration
	game.WhiteLeft = game.Duration

	lagCompensation, _ := time.ParseDuration("2s")
	initialTime := time.Now().Add(lagCompensation)
	game.BlackTurnStarted = initialTime
	game.WhiteTurnStarted = initialTime

	// Store game
	id, err = i.games.Store(game)
	if err != nil {
		return 0, err
	}

	return id, nil
}

func (i GamesInteractor) GetBoard(secret string, gameID int) ([8][8]byte, error) {
	game, err := i.getGameForUserSecret(secret, gameID)

	if err != nil {
		return [8][8]byte{}, err
	}
	return game.Board, nil

}

func (i GamesInteractor) getGameForUserSecret(secret string, gameID int) (*Game, error) {
	user, err := i.users.FindBySecret(secret, "")
	if err != nil {
		return nil, errors.New("couldnt find user by that id")
	}

	return i.getGameForUser(user.ID, gameID), nil
}

func (i GamesInteractor) getGameForUser(userID int, gameID int) *Game {
	games, err := i.games.FindByUserId(userID)
	if err != nil {
		return nil
	}

	Debug(fmt.Sprintf("  Looking for gameID %d", gameID))
	var game *Game
	found := false
	for _, game = range games {
		if game.ID == gameID {
			found = true
			break
		}
	}

	if found == true {
		return game
	}

	return nil
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

		for duration, queue := range matchRequests {
			if len(queue) <= 1 {
				continue
			}
			Debug(fmt.Sprintf("Matching a %d min game with queue %+v\n", duration, queue))

			durationTime, _ := time.ParseDuration(strconv.Itoa(duration) + "m")

			Debug(fmt.Sprintf("Matching a %d min game with queue %+v\n", duration, queue))
			// Use the first two requests to create a game
			game := Game{
				WhiteUser: queue[0].UserID,
				BlackUser: queue[1].UserID,
				Duration:  durationTime,
			}

			// Randomize color
			if rand.Float32() < .5 {
				temp := game.WhiteUser
				game.WhiteUser = game.BlackUser
				game.BlackUser = temp
			}

			i.Create(&game)
			log.Printf("INFO: Created game: %+v\n", game)

			// Delete the first two requests from the store
			_, err = i.matchRequests.Delete(queue[0].ID)
			if err != nil {
				log.Printf("ERROR: %v\n", err)
			}
			_, err = i.matchRequests.Delete(queue[1].ID)
			if err != nil {
				log.Printf("ERROR: %v\n", err)
			}
		}
	}
}

func DefaultBoard() [8][8]byte {
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
func (i *GamesInteractor) MakeMove(secret string, gameID int, move string) error {
	user, err := i.users.FindBySecret(secret, "")
	if err != nil {
		return errors.New("--couldnt find user by that id")
	}

	game := i.getGameForUser(user.ID, gameID)
	if game == nil {
		return errors.New("--couldnt find game")
	}

	if game.WhiteTurn && game.WhiteUser != user.ID ||
		!game.WhiteTurn && game.BlackUser != user.ID {
		Debug(fmt.Sprintf("--Wrong turn for user %s for game: %+v\n", secret, game))

		return errors.New("--not your turn")
	}

	err = i.games.MakeMove(game, &user, move)
	if err != nil {
		return err
	}

	lagCompensation, _ := time.ParseDuration("250ms")
	if user.ID == game.WhiteUser {
		game.BlackTurnStarted = time.Now().Add(lagCompensation)
		game.WhiteLeft = game.Duration - time.Now().Sub(game.WhiteTurnStarted)
	} else {
		game.WhiteTurnStarted = time.Now().Add(lagCompensation)
		game.BlackLeft = game.Duration - time.Now().Sub(game.BlackTurnStarted)
	}

	return nil
}

func (i *GamesInteractor) GetMove(secret string, gameID int) (string, *Game, error) {
	user, err := i.users.FindBySecret(secret, "")
	if err != nil {
		return "", nil, errors.New("--couldnt find user by that id")
	}

	game := i.getGameForUser(user.ID, gameID)

	if game == nil {
		msg := fmt.Sprintf("--couldnt find game for secret: %s gameID: %d", secret, gameID)
		Debug(msg)
		return "", nil, errors.New(msg)
	}

	move, err := i.games.GetMove(game, &user)
	return move, game, err
}
