package inmemory

import (
	"errors"

	"github.com/hashcacher/ChessGoNeue/Server/v2/core"
)

type Users struct {
	users map[int]core.User
}

func NewUsers(users map[int]core.User) Users {
	return Users{
		users: users,
	}
}

func (r *Users) Store(user core.User) error {
	r.users[user.Id] = user
	return nil
}

func (r *Users) FindById(id int) (core.User, error) {
	return r.users[id], nil
}

func (r *Users) Update(user core.User) error {
	_, ok := r.users[user.Id]
	if !ok {
		return errors.New("user does not exist")
	}
	r.users[user.Id] = user
	return nil
}
