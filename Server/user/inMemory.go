package user

import (
	"errors"

	"github.com/hashcacher/ChessGoNeue/Server/v2/core"
)

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
