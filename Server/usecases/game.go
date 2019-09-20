package main

import (
	"github.com/hashcacher/ChessGoServer/domain"
)

// GameInterractor is a struct that holds data to be injected for use cases
type GameInterractor struct {
	UserRepository domain.UserRepository
	GameRepository domain.GameRepository
}

// ExecuteMove validates a user and then performs a move
func (i *GameInterractor) ExecuteMove(m string, userId, gameId int) {
	// (UserRepository) Validate user is in match and it is their turn
	// (GameRepository) Perform update
	// (-)              Notify other user about the update
}
