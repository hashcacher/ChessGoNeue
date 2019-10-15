package inmemory

import (
	"errors"
	"log"

	"github.com/hashcacher/ChessGoNeue/Server/v2/core"
	"sync"
)

type Users struct {
	autoIncrement int
	secretUserMap map[string]core.User
	idUserMap     map[int]core.User
	lock          sync.RWMutex
}

func NewUsers() Users {
	return Users{
		secretUserMap: make(map[string]core.User),
		idUserMap:     make(map[int]core.User),
		lock:          sync.RWMutex{},
	}
}

func (r *Users) getNextAutoincrementID() int {
	r.autoIncrement++
	return r.autoIncrement
}

func (r *Users) Store(user core.User) (int, error) {
	user.ID = r.getNextAutoincrementID()

	r.lock.Lock()
	defer r.lock.Unlock()

	r.secretUserMap[user.Secret] = user
	r.idUserMap[user.ID] = user
	return user.ID, nil
}

func (r *Users) FindByID(id int) (core.User, error) {
	r.lock.RLock()
	defer r.lock.RUnlock()
	return r.idUserMap[id], nil
}

func (r *Users) FindBySecret(secret string, name string) (core.User, error) {
	if secret == "" {
		return core.User{}, errors.New("secret cannot be empty string")
	}

	r.lock.RLock()
	user, ok := r.secretUserMap[secret]
	r.lock.RUnlock()
	// Create a new user if one doesn't exist with that secret
	if !ok {
		user = core.User{Secret: secret, Name: name}
		id, err := r.Store(user)
		if err != nil {
			return core.User{}, nil
		}
		user.ID = id
		log.Printf("%v", user)
		return user, nil
	}
	// Return the user
	return user, nil
}

func (r *Users) Update(user core.User) error {
	r.lock.Lock()
	defer r.lock.Unlock()

	_, ok := r.idUserMap[user.ID]
	if !ok {
		return errors.New("user does not exist")
	}
	r.idUserMap[user.ID] = user
	r.secretUserMap[user.Secret] = user
	return nil
}
