package main

import (
	domain "github.com/hashcacher/ChessGoServer/domain"
)

// MatchRequestInterractor is a struct that holds data to be injected for use cases
type MatchRequestInterractor struct {
	MatchRequestRepository domain.MatchRequestRepository
	UserRepository         domain.UserRepository
	GameRepository         domain.GameRepository
}

// MatchMe will take in a user, create a match request, and wait for a notification
// saying a match was succesful
func (i *MatchRequestInterractor) MatchMe(userId int) {
	// (UserRepo) Validate user exists

	// Does a valid match request from another user already exist?
	//	yes:
	//		(MatchRequestRepository) delete the match request
	//		(GameRepository)         create a new game
	//    (-)                      notify the other user that a match was found and a game was created
	//	no:
	//		(MatchRequestRepository) create a match request
	//    (-)                      wait for a notification that a match was found and a game was created
	// Return the new game
}
