package interfaces

import (
	"github.com/hashcacher/ChessGoServer/domain"
)

type MemGameRepository struct {
	games map[string]Game
}

func (r *MemGameRepository) Store(game domain.Game) error {
	r.games[game.Id] = game
}

func (r *MemGameRepository) FindById(id int) (domain.Game, error) {
	return r.games[game.Id]
}

func (r *MemGameRepository) Update(game domain.Game) error {
	r.games[game.Id] = game
}
