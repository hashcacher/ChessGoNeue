package core

// User stores user profile data
type User struct {
	Id       int
	Username string
}

// UserRepository is the use case for User entitiy
type Users interface {
	Store(User) error
	FindById(id int) (User, error)
	Update(User) error
}

// UserInteractor is used to interact with user repositories and other related repositories
type UsersInteractor struct {
	users Users
}

// NewUsersInteractor generates a new UsersInteractor from the given Users store
func NewUsersInteractor(users Users) UsersInteractor {
	return UsersInteractor{
		users,
	}
}

// FindById fetches the user from the repository and returns it
func (i *UsersInteractor) FindById(id int) (User, error) {
	user, err := i.users.FindById(id)
	if err != nil {
		return User{}, err
	}
	return user, nil
}
