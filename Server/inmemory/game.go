package inmemory

import (
	"errors"
	"fmt"
	"github.com/hashcacher/ChessGoNeue/Server/v2/core"
	"strings"
	"sync"
	"time"
)

type Games struct {
	autoIncrement int
	games         map[int]*core.Game
	// Map from userID to notification channel
	storedEventsByUserID map[int]chan *core.Game

	// Map from userID to move channel
	moveEventsByUserID map[int]chan string

	// We're dealing with real threads here
	gamesLock  sync.RWMutex
	moveLock   sync.RWMutex
	storedLock sync.RWMutex
}

func NewGames(games map[int]*core.Game) Games {
	return Games{
		games:                games,
		storedEventsByUserID: make(map[int]chan *core.Game),
		moveEventsByUserID:   make(map[int]chan string),
		storedLock:           sync.RWMutex{},
		moveLock:             sync.RWMutex{},
		gamesLock:            sync.RWMutex{},
	}
}

func (g *Games) getNextAutoincrementID() int {
	g.autoIncrement++
	return g.autoIncrement
}

func (g *Games) Store(game *core.Game) (int, error) {
	game.ID = g.getNextAutoincrementID()

	g.gamesLock.Lock()
	g.games[game.ID] = game
	g.gamesLock.Unlock()

	// Notify for white user
	g.storedLock.RLock()
	notifyChannel, ok := g.storedEventsByUserID[game.WhiteUser]
	g.storedLock.RUnlock()
	if !ok {
		notifyChannel = make(chan *core.Game, 2)
		g.storedLock.Lock()
		g.storedEventsByUserID[game.WhiteUser] = notifyChannel
		g.storedLock.Unlock()
	}
	notifyChannel <- game

	// Notify for black user
	g.storedLock.RLock()
	notifyChannel, ok = g.storedEventsByUserID[game.BlackUser]
	g.storedLock.RUnlock()
	if !ok {
		notifyChannel = make(chan *core.Game, 2)
		g.storedLock.Lock()
		g.storedEventsByUserID[game.BlackUser] = notifyChannel
		g.storedLock.Unlock()
	}
	notifyChannel <- game

	// Create move channels
	whiteMoveChannel := make(chan string, 2)
	blackMoveChannel := make(chan string, 2)
	g.moveLock.Lock()
	g.moveEventsByUserID[game.WhiteUser] = whiteMoveChannel
	g.moveEventsByUserID[game.BlackUser] = blackMoveChannel
	g.moveLock.Unlock()

	return game.ID, nil
}

func (g *Games) ListenForMoveByUserID(userID int) (move string, err error) {
	g.moveLock.RLock()
	notifyChannel, _ := g.moveEventsByUserID[userID]
	g.moveLock.RUnlock()

	return <-notifyChannel, nil
}

func (g *Games) ListenForStoreByUserID(userID int) (game *core.Game, err error) {
	g.storedLock.RLock()
	notifyChannel, ok := g.storedEventsByUserID[userID]
	g.storedLock.RUnlock()

	// If the channel doesn't exist, create it
	if !ok {
		notifyChannel = make(chan *core.Game, 2)
		g.storedLock.Lock()
		g.storedEventsByUserID[userID] = notifyChannel
		g.storedLock.Unlock()
	}
	return <-notifyChannel, nil
}

func (g *Games) FindById(id int) (core.Game, error) {
	g.gamesLock.RLock()
	defer g.gamesLock.RUnlock()
	return *g.games[id], nil
}

func (g *Games) FindByUserId(userID int) ([]*core.Game, error) {
	g.gamesLock.RLock()
	defer g.gamesLock.RUnlock()

	games := []*core.Game{}
	for _, game := range g.games {
		if game.BlackUser == userID || game.WhiteUser == userID {
			games = append(games, game)
		}
	}

	return games, nil
}

func (g *Games) Update(game *core.Game) error {
	g.gamesLock.Lock()
	defer g.gamesLock.Unlock()

	g.games[game.ID] = game
	return nil
}

func (g *Games) GetMove(game *core.Game, user *core.User) (string, error) {
	return g.ListenForMoveByUserID(user.ID)
}

// MakeMove assumes a previous function checked it's this user's turn
func (g *Games) MakeMove(game *core.Game, user *core.User, move string) error {
	squares := strings.Split(move, ">")

	if len(squares) == 2 {
		// Chess move
		fromX, fromY, toX, toY := unpackChessMove(squares[0], squares[1])

		if game.Board[fromX][fromY] == ' ' {
			return errors.New("Nothing at " + squares[0])
		}

		core.Debug(fmt.Sprintf("  Game %d user %d: %d,%d -> %d,%d\n", game.ID, user.ID, fromX, fromY, toX, toY))

		game.Board[toX][toY] = game.Board[fromX][fromY]
		game.Board[fromX][fromY] = ' '
	} else if len(squares) == 1 {
		// Go Move
		toX, toY := unpackGoMove(squares[0])

		if game.Board[toX][toY] != ' ' {
			return errors.New("Cant put a go stone on a square thats taken: " + squares[0])
		}

		var stone byte
		if user.ID == game.BlackUser {
			stone = 's'
		} else {
			stone = 'S'
		}
		game.Board[toX][toY] = stone
	} else {
		return errors.New("Invalid Move Format")
	}

	// Next player's turn
	game.WhiteTurn = !game.WhiteTurn

	// Actually override the game in the map
	g.gamesLock.Lock()
	g.games[game.ID] = game
	g.gamesLock.Unlock()

	core.Debug(fmt.Sprintf("    Updated game %d board", game.ID))

	// Get out opponent's player ID
	var oppID int
	if user.ID == game.BlackUser {
		oppID = game.WhiteUser
	} else {
		oppID = game.BlackUser
	}

	// Send move to other player
	core.Debug(fmt.Sprintf("    Sending game %d move events to user %d", game.ID, oppID))
	g.moveLock.RLock()
	moveChan := g.moveEventsByUserID[oppID]
	g.moveLock.RUnlock()

	moveChan <- move

	core.Debug(fmt.Sprintf("    Sent game %d move events", game.ID))

	return nil
}

func unpackChessMove(fromMove string, toMove string) (int, int, int, int) {
	fromCoords := strings.Split(fromMove, ",")
	toCoords := strings.Split(toMove, ",")

	toX := int(toCoords[0][0] - '0')
	toY := int(toCoords[1][0] - '0')
	fromX := int(fromCoords[0][0] - '0')
	fromY := int(fromCoords[1][0] - '0')

	return fromX, fromY, toX, toY
}

func unpackGoMove(move string) (int, int) {
	coords := strings.Split(move, ",")
	from := charToInt(coords[0][0])
	to := charToInt(coords[1][0])

	return from, to
}

func charToInt(char byte) int {
	return int(char - '0')
}

func (g *Games) GetBoard(game *core.Game) [8][8]byte {
	g.gamesLock.RLock()
	defer g.gamesLock.RUnlock()
	return g.games[game.ID].Board
}

func (g *Games) ListenForTimeout(game *core.Game, userID int) error {
	zero, _ := time.ParseDuration("0s")

	for {
		// This only needs to run for one turn
		if game.WhiteTurn && userID != game.WhiteUser ||
			!game.WhiteTurn && userID != game.BlackUser {
			break
		}

		var left time.Duration
		var turnStarted time.Time
		var color string
		if userID == game.WhiteUser {
			left = game.WhiteLeft
			turnStarted = game.WhiteTurnStarted
			color = "white"
		} else {
			left = game.BlackLeft
			turnStarted = game.BlackTurnStarted
			color = "black"
		}

		left = left - time.Now().Sub(turnStarted)

		core.Debug(fmt.Sprintf("Game %d %s has %+v left", game.ID, color, left))
		if left <= zero {
			msg := "timeout " + color
			core.Debug(fmt.Sprintf("Game %d %s", game.ID, msg))
			notifyChannel, _ := g.moveEventsByUserID[game.WhiteUser]
			notifyChannel <- msg
			notifyChannel, _ = g.moveEventsByUserID[game.BlackUser]
			notifyChannel <- msg
		}

		core.Debug(fmt.Sprintf("Game %d sleeping %+v", game.ID, left))
		time.Sleep(left)
	}

	return nil
}
