package inmemory

import (
	"github.com/hashcacher/ChessGoNeue/Server/v2/core"
)

type Games struct {
	autoIncrement int
	games         map[int]core.Game
	// Map from userID to notification channel
	storedEventsByUserID map[int]chan core.Game
}

func NewGames(games map[int]core.Game) Games {
	return Games{
		games:                games,
		storedEventsByUserID: make(map[int]chan core.Game),
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
		notifyChannel = make(chan core.Game, 2)
		g.storedEventsByUserID[game.WhiteUser] = notifyChannel
	}
	notifyChannel <- game
	// Notify for black user
	notifyChannel, ok = g.storedEventsByUserID[game.BlackUser]
	if !ok {
		notifyChannel = make(chan core.Game, 2)
		g.storedEventsByUserID[game.BlackUser] = notifyChannel
	}
	notifyChannel <- game
	// Return
	return game.ID, nil
}

func (g *Games) ListenForStoreByUserID(userID int) (game core.Game, err error) {
	notifyChannel, ok := g.storedEventsByUserID[userID]
	// If the channel doesn't exist, create it
	if !ok {
		notifyChannel = make(chan core.Game, 2)
		g.storedEventsByUserID[userID] = notifyChannel
	}
	return <-notifyChannel, nil
}

func (g *Games) FindById(id int) (core.Game, error) {
	return g.games[id], nil
}

func (g *Games) Update(game core.Game) error {
	g.games[game.ID] = game
	return nil
}
