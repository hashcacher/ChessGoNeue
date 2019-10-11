package core

import (
	"reflect"
)

// MatchRequest holds data to determine how to match users together for a game
type MatchRequest struct {
	ID     int `json:"id"`
	UserID int `json:"user"`
	Elo    int `json:"elo"`
}

// MatchRequests is the use case for Match entitiy
type MatchRequests interface {
	// Store this matchrequest into the database
	Store(MatchRequest) error
	ListenForStore() (MatchRequest, error)
	// Find by a specific user
	FindByUserID(userID int) (MatchRequest, error)
	// Find all
	FindAll() ([]MatchRequest, error)
	// Delete a MatchRequest by ID
	Delete(id int) (deleted int, err error)
	// Delete a MatchRequest by userID
	DeleteByUserID(id int) (deleted int, err error)
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
func (i *MatchRequestsInteractor) MatchMe(userID int) (game *Game, err error) {
	// Make sure the user isn't already queued for a game by seeing if they have a match
	// request created
	matchRequest, err := i.matchRequests.FindByUserID(userID)
	isMatchRequestEmpty := reflect.DeepEqual(matchRequest, MatchRequest{})
	if err != nil {
		return &Game{}, err
	}

	// Create request
	if isMatchRequestEmpty {
		newMatchRequest := MatchRequest{UserID: userID}
		i.matchRequests.Store(newMatchRequest)
	}

	// Listen (blocking) for notify
	game, err = i.games.ListenForStoreByUserID(userID)

	// Return the gameID
	return game, nil
}

// DeleteMatchMe will dequeue a user from matchmaking
func (i *MatchRequestsInteractor) DeleteMatchMe(userID int) (err error) {
	_, err = i.matchRequests.DeleteByUserID(userID)
	return err
}
