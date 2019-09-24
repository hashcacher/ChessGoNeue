package inmemory

import "github.com/hashcacher/ChessGoNeue/Server/v2/core"

type Games struct {
	autoIncrement int
	games         map[int]core.Game
	// Map from userID to notification channel
	createdNotifiers map[int]chan int
}

func NewGames(games map[int]core.Game) Games {
	return Games{
		games:            games,
		createdNotifiers: make(map[int]chan int),
	}
}

func (g *Games) getNextAutoincrementID() int {
	g.autoIncrement++
	return g.autoIncrement
}

func (g *Games) Store(game core.Game) (int, error) {
	game.ID = g.getNextAutoincrementID()
	g.games[game.ID] = game
	g.NotifyGameCreated(game.WhiteUser, game.ID)
	g.NotifyGameCreated(game.BlackUser, game.ID)
	return game.ID, nil
}

func (g *Games) FindById(id int) (core.Game, error) {
	return g.games[id], nil
}

func (g *Games) Update(game core.Game) error {
	g.games[game.ID] = game
	return nil
}

func (g *Games) ListenForGameCreatedNotification(userID int) (gameID int) {
	notifyChannel := make(chan int, 2)
	g.createdNotifiers[userID] = notifyChannel
	return <-notifyChannel
}

func (g *Games) NotifyGameCreated(userID, gameID int) error {
	notifyChannel := g.createdNotifiers[userID]
	notifyChannel <- gameID
	delete(g.createdNotifiers, userID)
	return nil
}
