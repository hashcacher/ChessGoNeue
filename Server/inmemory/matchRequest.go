package inmemory

import (
	"github.com/hashcacher/ChessGoNeue/Server/v2/core"
	"sync"
)

const (
	allStoredEventsBufferSize = 100
)

type MatchRequests struct {
	autoIncrement int
	// map from matchRequest ID to matchRequest
	matchRequests map[int]core.MatchRequest
	// Map from userID to notification channel
	allStoredEvents chan core.MatchRequest

	lock sync.RWMutex
}

func (r *MatchRequests) getNextAutoincrementID() int {
	r.autoIncrement++
	return r.autoIncrement
}

func NewMatchRequests(matchRequests map[int]core.MatchRequest) MatchRequests {
	return MatchRequests{
		matchRequests:   matchRequests,
		allStoredEvents: make(chan core.MatchRequest, allStoredEventsBufferSize),
		lock:            sync.RWMutex{},
	}
}

func (r *MatchRequests) Store(matchRequest core.MatchRequest) error {
	matchRequest.ID = r.getNextAutoincrementID()

	r.lock.Lock()
	defer r.lock.Unlock()
	r.matchRequests[matchRequest.ID] = matchRequest
	r.allStoredEvents <- matchRequest

	return nil
}

func (r *MatchRequests) AttemptMatch(matchRequest core.MatchRequest) error {
	matchRequest.ID = r.getNextAutoincrementID()

	r.lock.Lock()
	defer r.lock.Unlock()
	r.matchRequests[matchRequest.ID] = matchRequest

	return nil
}

func (r *MatchRequests) FindByUserID(userID int) (core.MatchRequest, error) {
	r.lock.RLock()
	defer r.lock.RUnlock()

	// Find first request and return it
	for _, matchRequest := range r.matchRequests {
		if matchRequest.UserID == userID {
			return matchRequest, nil
		}
	}
	// If the map of match requests is empty, just return empty
	return core.MatchRequest{}, nil
}

func (r *MatchRequests) FindAll() ([]core.MatchRequest, error) {
	r.lock.RLock()
	defer r.lock.RUnlock()

	matchRequests := []core.MatchRequest{}
	for _, matchRequest := range r.matchRequests {
		matchRequests = append(matchRequests, matchRequest)
	}
	return matchRequests, nil
}

func (r *MatchRequests) Delete(id int) (deleted int, err error) {
	r.lock.RLock()
	defer r.lock.RUnlock()

	_, ok := r.matchRequests[id]
	deleted = 0
	if ok {
		deleted = 1
	}
	delete(r.matchRequests, id)
	return deleted, nil
}

func (r *MatchRequests) ListenForStore() (core.MatchRequest, error) {
	return <-r.allStoredEvents, nil
}
