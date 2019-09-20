package game

import "github.com/hashcacher/ChessGoNeue/Server/v2/core"

type IMGameRepository struct {
	games map[int]core.Game
}

func (r *IMGameRepository) Store(game core.Game) error {
	r.games[game.Id] = game
	return nil
}

func (r *IMGameRepository) FindById(id int) (core.Game, error) {
	return r.games[id], nil
}

func (r *IMGameRepository) Update(game core.Game) error {
	r.games[game.Id] = game
	return nil
}
