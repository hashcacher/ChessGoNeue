package core

// Game stores chessgo game state such as the current board and who is in the game
type Game struct {
	ID        int        `json:"id"`
	BlackUser int        `json:"blackUser"`
	WhiteUser int        `json:"WhiteUser"`
	Board     [8][8]byte `json:"board"`
}

// Games is the use case for Game entitiy
type Games interface {
	Store(Game) error
	FindById(id int) (Game, error)
	Update(Game) error
}

// GamesInteractor is a struct that holds data to be injected for use cases
type GamesInteractor struct {
	games Games
	users Users
}

// NewGamesInteractor generates a new GamesInteractor from the given Users store
func NewGamesInteractor(games Games, users Users) GamesInteractor {
	return GamesInteractor{
		games,
		users,
	}
}

// ExecuteMove validates a user and then performs a move
func (i *GamesInteractor) ExecuteMove(m string, userID, gameId int) {
	// (UserRepository) Validate user is in match and it is their turn
	// (GameRepository) Perform update
	// (-)              Notify other user about the update
}
