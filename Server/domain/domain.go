package domain

// The domain consists of our application's entities. These are the inner most layer of our application
// and should not depend on anything. They define the enterprise wide business logic of our application.
// for example here we can implement the logic for making a move in chessgo, and how that affects a
// game entity, but should not include how to handle storing that data.
//
// Entity: an object with methods that can be used by many different applications
//  and depends on nothing.

// User stores user profile data
type User struct {
	Id       int
	Username string
}

// Game stores chessgo game state such as the current board and who is in the game
type Game struct {
	Id        int
	BlackUser int
	WhiteUser int
	Board     string
}

// MatchRequest holds data to determine how to match users together for a game
type MatchRequest struct {
	Id   int
	User int
	Elo  int
}

// UserRepository is the use case for User entitiy
type UserRepository interface {
	Store(User) error
	FindById(id int) (User, error)
	Update(User) error
}

// GameRepository is the use case for Game entitiy
type GameRepository interface {
	Store(Game) error
	FindById(id int) (Game, error)
	Update(Game) error
}

// MatchRequestRepository is the use case for Match entitiy
type MatchRequestRepository interface {
	Store(MatchRequest) error
	FindAllMatchRequests(User)
	Delete(id int) (deleted int, err error)
}
