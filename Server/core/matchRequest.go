package core

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
func (i *MatchRequestsInteractor) MatchMe(clientID string) {

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
