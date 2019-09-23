package inmemory

import (
	"errors"

	"github.com/hashcacher/ChessGoNeue/Server/v2/core"
)

type Users struct {
	SecretUserMap map[string]core.User
	idUserMap       map[int]core.User
}

func NewUsers() Users {
	return Users{
		SecretUserMap: make(map[string]core.User),
		idUserMap:       make(map[int]core.User),
	}
}

func (r *Users) Store(user core.User) error {
	r.SecretUserMap[user.Secret] = user
	r.idUserMap[user.ID] = user
	return nil
}

func (r *Users) FindBySecret(id int) (core.User, error) {
	return r.idUserMap[id], nil
}

func (r *Users) Update(user core.User) error {
	_, ok := r.idUserMap[user.ID]
	if !ok {
		return errors.New("user does not exist")
	}
	r.idUserMap[user.ID] = user
	r.SecretUserMap[user.Secret] = user
	return nil
}
