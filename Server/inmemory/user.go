package inmemory

import (
	"errors"

	"github.com/hashcacher/ChessGoNeue/Server/v2/core"
)

type Users struct {
	autoIncrement int
	secretUserMap map[string]core.User
	idUserMap     map[int]core.User
}

func NewUsers() Users {
	return Users{
		secretUserMap: make(map[string]core.User),
		idUserMap:     make(map[int]core.User),
	}
}

func (r *Users) getNextAutoincrementID() int {
	r.autoIncrement++
	return r.autoIncrement
}

func (r *Users) Store(user core.User) (int, error) {
	user.ID = r.getNextAutoincrementID()
	r.secretUserMap[user.Secret] = user
	r.idUserMap[user.ID] = user
	return user.ID, nil
}

func (r *Users) FindByID(id int) (core.User, error) {
	return r.idUserMap[id], nil
}

func (r *Users) FindBySecret(secret string) (core.User, error) {
	return r.secretUserMap[secret], nil
}

func (r *Users) Update(user core.User) error {
	_, ok := r.idUserMap[user.ID]
	if !ok {
		return errors.New("user does not exist")
	}
	r.idUserMap[user.ID] = user
	r.secretUserMap[user.Secret] = user
	return nil
}
