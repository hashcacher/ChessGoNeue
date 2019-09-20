package implementations

import (
	"errors"

	"github.com/hashcacher/ChessGoNeue/Server/v2/core"
)

// ---------------
// Game Repository
// ---------------
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

// ---------------
// Match Request Repository
// ---------------
type IMMatchRequestRepository struct {
	matchRequests map[int]core.MatchRequest
}

func (r IMMatchRequestRepository) Store(matchRequest core.MatchRequest) error {
	r.matchRequests[matchRequest.Id] = matchRequest
	return nil
}

func (r IMMatchRequestRepository) FindAllMatchRequestsByUserId(userId int) []core.MatchRequest {
	// TODO
	return []core.MatchRequest{}
}

func (r IMMatchRequestRepository) Delete(id int) (deleted int, err error) {
	_, ok := r.matchRequests[id]
	deleted = 0
	if ok {
		deleted = 1
	}
	delete(r.matchRequests, id)
	return deleted, nil
}

// ---------------
// User Repository
// ---------------
type IMUserRepository struct {
	users map[int]core.User
}

func NewIMUserRepIMUserRepository(users map[int]core.User) IMUserRepository {
	return IMUserRepository{
		users: users,
	}
}

func (r *IMUserRepository) Store(user core.User) error {
	r.users[user.Id] = user
	return nil
}

func (r *IMUserRepository) FindById(id int) (core.User, error) {
	return r.users[id], nil
}

func (r *IMUserRepository) Update(user core.User) error {
	_, ok := r.users[user.Id]
	if !ok {
		return errors.New("user does not exist")
	}
	r.users[user.Id] = user
	return nil
}
