package core

import (
	"errors"
	"reflect"
)

// MatchRequest holds data to determine how to match users together for a game
type MatchRequest struct {
	ID   int `json:"id"`
	User int `json:"user"`
	Elo  int `json:"elo"`
}

// MatchRequests is the use case for Match entitiy
type MatchRequests interface {
	Store(MatchRequest) error
	FindAllMatchRequestsByUserId(userID int) []MatchRequest
	// Go through MatchRequests and see if there is a valid match
	FindMatchForUser(userID int) MatchRequest
	// Block and listen for a notification saying a game was created for you
	ListenForGameCreatedNotify(userID int) (gameID int)
	// Notify someone who is listening that a game was created for them
	NotifyGameCreated(userID int) (gameID int)
	Delete(id int) (deleted int, err error)
}

// MatchRequestsInteractor is a struct that holds data to be injected for use cases
type MatchRequestsInteractor struct {
	matchRequests MatchRequests
	users         Users
	games         Games
}

// NewMatchRequestsInteractor generates a new MatchRequestsInteractor from the given Users store
func NewMatchRequestsInteractor(matchRequests MatchRequests, users Users, games Games) MatchRequestsInteractor {
	return MatchRequestsInteractor{
		matchRequests,
		users,
		games,
	}
}

// MatchMe will take in a user, create a match request, and wait for a notification
// saying a match was succesful
func (i *MatchRequestsInteractor) MatchMe(clientID string) (gameID int, err error) {
	// Make sure user exists and get their info
	user, err := i.users.FindByClientID(clientID)
	isUserEmpty := reflect.DeepEqual(user, User{})
	if err != nil {
		return 0, err
	}
	if isUserEmpty {
		return 0, errors.New("could not find user with that client id")
	}

	// Does a valid match request from another user already exist?
	matchRequest := i.matchRequests.FindMatchForUser(user.ID)
	isMREmpty := reflect.DeepEqual(matchRequest, MatchRequest{})
	// If the match request is empty, create one and waot
	if isMREmpty {
		// Create request
		matchRequest := MatchRequest{User: user.ID}
		i.matchRequests.Store(matchRequest)
		// Listen (blocking) for notify
		gameID := i.matchRequests.ListenForGameCreatedNotify(user.ID)
		// Return gameID we were notified about
		return gameID, nil
	}

	// Delete the match request
	i.matchRequests.Delete(matchRequest.ID)
	// Create a game
	// TODO: randomize black and white
	game := Game{
		WhiteUser: user.ID,
		BlackUser: matchRequest.User,
	}
	// Store the new game
	id, err := i.games.Store(game)
	if err != nil {
		return 0, err
	}
	// Return the id
	return id, nil
}
