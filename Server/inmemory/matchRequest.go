package inmemory

import "github.com/hashcacher/ChessGoNeue/Server/v2/core"

type MatchRequests struct {
	// map from matchRequest ID to matchRequest
	matchRequests map[int]core.MatchRequest
	// map from userID to notification channel
	notifiers map[int]chan int
}

func NewMatchRequests(matchRequests map[int]core.MatchRequest) MatchRequests {
	return MatchRequests{
		matchRequests,
		notifiers: make(map[int]chan int),
	}
}

func (r MatchRequests) Store(matchRequest core.MatchRequest) error {
	r.matchRequests[matchRequest.ID] = matchRequest
	return nil
}

func (r MatchRequests) FindMatchForUser(userID int) core.MatchRequest {
	// Find first request and return it
	for _, matchRequest := range r.matchRequests {
		return matchRequest
	}
	// If the map of match requests is empty, just return empty
	return core.MatchRequest{}
}

func (r MatchRequests) ListenForGameCreatedNotify(userID int) (gameID int) {
	notifyChannel := make(chan int, 2)
	r.notifiers[userID] = notifyChannel
	return <-notifyChannel
}

func (r MatchRequests) NotifyGameCreated(userID, gameID int) {
	notifyChannel := r.notifiers[userID]
	notifyChannel <- gameID
	delete(r.notifiers, userID)
}

func (r MatchRequests) Delete(id int) (deleted int, err error) {
	_, ok := r.matchRequests[id]
	deleted = 0
	if ok {
		deleted = 1
	}
	delete(r.matchRequests, id)
	return deleted, nil
}
