package core

import "errors"

// User stores user profile data
type User struct {
	ID       int    `json:"id"`
	ClientID string `json:"clientID"`
	Username string `json:"username"`
}

// Users is the use case for User entitiy
type Users interface {
	Store(User) (User, error)
	FindByClientID(clientID string) (User, error)
	Update(User) error
}

// UsersInteractor is used to interact with user repositories and other related repositories
type UsersInteractor struct {
	users Users
}

// NewUsersInteractor generates a new UsersInteractor from the given Users store
func NewUsersInteractor(users Users) UsersInteractor {
	return UsersInteractor{
		users,
	}
}

// Create creates a new user
func (i *UsersInteractor) Create(user User) (User, error) {
	if len(user.Username) == 0 {
		return User{}, errors.New("must specify username")
	}
	storedUser, err := i.users.Store(user)
	if err != nil {
		return User{}, err
	}
	return storedUser, nil
}

// FindByClientID fetches the user from the repository and returns it
func (i *UsersInteractor) FindByClientID(id string) (User, error) {
	user, err := i.users.FindByClientID(id)
	if err != nil {
		return User{}, err
	}
	return user, nil
}
