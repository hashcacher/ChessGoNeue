package inmemory

import "github.com/hashcacher/ChessGoNeue/Server/v2/core"

const (
	createdNotifierBufferSize = 100
)

type MatchRequests struct {
	autoIncrement int
	// map from matchRequest ID to matchRequest
	matchRequests map[int]core.MatchRequest
	// map from userID to notification channel
	createdNotifier chan interface{}
}

func (r *MatchRequests) getNextAutoincrementID() int {
	r.autoIncrement++
	return r.autoIncrement
}

func NewMatchRequests(matchRequests map[int]core.MatchRequest) MatchRequests {
	return MatchRequests{
		matchRequests:   matchRequests,
		createdNotifier: make(chan interface{}, createdNotifierBufferSize),
	}
}

func (r *MatchRequests) Store(matchRequest core.MatchRequest) error {
	matchRequest.ID = r.getNextAutoincrementID()
	r.matchRequests[matchRequest.ID] = matchRequest
	return nil
}

func (r *MatchRequests) FindMatchRequestByUserID(userID int) (core.MatchRequest, error) {
	// Find first request and return it
	for _, matchRequest := range r.matchRequests {
		if matchRequest.User == userID {
			return matchRequest, nil
		}
	}
	// If the map of match requests is empty, just return empty
	return core.MatchRequest{}, nil
}

func (r *MatchRequests) Delete(id int) (deleted int, err error) {
	_, ok := r.matchRequests[id]
	deleted = 0
	if ok {
		deleted = 1
	}
	delete(r.matchRequests, id)
	return deleted, nil
}
