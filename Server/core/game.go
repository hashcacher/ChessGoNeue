package core

// Game stores chessgo game state such as the current board and who is in the game
type Game struct {
	Id        int
	BlackUser int
	WhiteUser int
	Board     string
}

// GameRepository is the use case for Game entitiy
type GameRepository interface {
	Store(Game) error
	FindById(id int) (Game, error)
	Update(Game) error
}

// GameInteractor is a struct that holds data to be injected for use cases
type GameInteractor struct {
	UserRepository UserRepository
	GameRepository GameRepository
}

// ExecuteMove validates a user and then performs a move
func (i *GameInteractor) ExecuteMove(m string, userId, gameId int) {
	// (UserRepository) Validate user is in match and it is their turn
	// (GameRepository) Perform update
	// (-)              Notify other user about the update
}
