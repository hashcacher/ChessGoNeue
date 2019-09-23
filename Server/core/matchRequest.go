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
	// Find a match request for a specific user
	FindMatchRequestByUserID(userID int) (MatchRequest, error)
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
func (i *MatchRequestsInteractor) MatchMe(Secret string) (gameID int, err error) {
	// Make sure user exists and get their info
	user, err := i.users.FindBySecret(Secret)
	isUserEmpty := reflect.DeepEqual(user, User{})
	if err != nil {
		return 0, err
	}
	if isUserEmpty {
		return 0, errors.New("could not find user with that client id")
	}

	// Make sure the user isn't already queued for a game by seeing if they have a match
	// request created
	matchRequest, err := i.matchRequests.FindMatchRequestByUserID(user.ID)
	isMatchRequestEmpty := reflect.DeepEqual(matchRequest, MatchRequest{})
	if err != nil {
		return 0, err
	}
	if !isMatchRequestEmpty {
		return 0, errors.New("you can only queue for one game at a time")
	}

	// Create request
	newMatchRequest := MatchRequest{User: user.ID}
	i.matchRequests.Store(newMatchRequest)

	// Listen (blocking) for notify
	gameID = i.games.ListenForGameCreatedNotification(user.ID)

	// Return the gameID
	return gameID, nil
}
