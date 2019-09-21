package inmemory

import "github.com/hashcacher/ChessGoNeue/Server/v2/core"

type Games struct {
	games map[int]core.Game
}

func NewGames(games map[int]core.Game) Games {
	return Games{
		games,
	}
}

func (r *Games) Store(game core.Game) error {
	r.games[game.ID] = game
	return nil
}

func (r *Games) FindById(id int) (core.Game, error) {
	return r.games[id], nil
}

func (r *Games) Update(game core.Game) error {
	r.games[game.ID] = game
	return nil
}
