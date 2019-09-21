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
	Store(User) (id int, err error)
	FindByClientID(clientID string) (User, error)
	FindByID(id int) (User, error)
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
func (i *UsersInteractor) Create(user User) (int, error) {
	if len(user.Username) == 0 {
		return 0, errors.New("username can't be empty")
	}
	id, err := i.users.Store(user)
	if err != nil {
		return 0, err
	}
	return id, nil
}

// FindByClientID fetches the user from the repository and returns it
func (i *UsersInteractor) FindByClientID(clientID string) (User, error) {
	user, err := i.users.FindByClientID(clientID)
	if err != nil {
		return User{}, err
	}
	return user, nil
}

// FindByID fetches the user from the repository and returns it
func (i *UsersInteractor) FindByID(id int) (User, error) {
	user, err := i.users.FindByID(id)
	if err != nil {
		return User{}, err
	}
	return user, nil
}
