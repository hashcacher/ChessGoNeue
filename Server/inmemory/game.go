package inmemory

import (
	"errors"
	"fmt"
	"github.com/hashcacher/ChessGoNeue/Server/v2/core"
	"strings"
)

type Games struct {
	autoIncrement int
	games         map[int]core.Game
	// Map from userID to notification channel
	storedEventsByUserID map[int]chan *core.Game

	// Map from userID to move channel
	moveEventsByUserID map[int]chan string
}

func NewGames(games map[int]core.Game) Games {
	return Games{
		games:                games,
		storedEventsByUserID: make(map[int]chan *core.Game),
		moveEventsByUserID:   make(map[int]chan string),
	}
}

func (g *Games) getNextAutoincrementID() int {
	g.autoIncrement++
	return g.autoIncrement
}

func (g *Games) Store(game core.Game) (int, error) {
	game.ID = g.getNextAutoincrementID()
	g.games[game.ID] = game
	// Notify for white user
	notifyChannel, ok := g.storedEventsByUserID[game.WhiteUser]
	if !ok {
		notifyChannel = make(chan *core.Game, 2)
		g.storedEventsByUserID[game.WhiteUser] = notifyChannel
	}
	notifyChannel <- &game
	// Notify for black user
	notifyChannel, ok = g.storedEventsByUserID[game.BlackUser]
	if !ok {
		notifyChannel = make(chan *core.Game, 2)
		g.storedEventsByUserID[game.BlackUser] = notifyChannel
	}
	notifyChannel <- &game
	// Return
	return game.ID, nil
}

func (g *Games) ListenForMoveByUserID(userID int) (move string, err error) {
	notifyChannel, ok := g.moveEventsByUserID[userID]
	// If the channel doesn't exist, create it
	if !ok {
		notifyChannel = make(chan string, 2)
		g.moveEventsByUserID[userID] = notifyChannel
	}
	return <-notifyChannel, nil
}

func (g *Games) ListenForStoreByUserID(userID int) (game *core.Game, err error) {
	notifyChannel, ok := g.storedEventsByUserID[userID]
	// If the channel doesn't exist, create it
	if !ok {
		notifyChannel = make(chan *core.Game, 2)
		g.storedEventsByUserID[userID] = notifyChannel
	}
	return <-notifyChannel, nil
}

func (g *Games) FindById(id int) (core.Game, error) {
	return g.games[id], nil
}

func (g *Games) FindByUserId(userID int) ([]*core.Game, error) {
	games := make([]*core.Game, 0)
	for _, game := range g.games {
		if game.BlackUser == userID || game.WhiteUser == userID {
			games = append(games, &game)
		}
	}

	return games, nil
}

func (g *Games) Update(game core.Game) error {
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

		fmt.Printf("%d,%d -> %d,%d", fromX, fromY, toX, toY)

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
	g.games[game.ID] = *game

	// Get out opponent's player ID
	var oppID int
	if user.ID == game.BlackUser {
		oppID = game.WhiteUser
	} else {
		oppID = game.BlackUser
	}

	// Send move to other player
	g.moveEventsByUserID[oppID] <- move

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
	return g.games[game.ID].Board
}
