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
	// map from game duration to matchRequest
	matchRequests map[int][]core.MatchRequest
	// Map from userID to notification channel
	allStoredEvents chan core.MatchRequest

	lock sync.RWMutex
}

func (r *MatchRequests) getNextAutoincrementID() int {
	r.autoIncrement++
	return r.autoIncrement
}

func NewMatchRequests() MatchRequests {
	return MatchRequests{
		matchRequests:   map[int][]core.MatchRequest{},
		allStoredEvents: make(chan core.MatchRequest, allStoredEventsBufferSize),
		lock:            sync.RWMutex{},
	}
}

func (r *MatchRequests) Store(matchRequest core.MatchRequest) error {
	matchRequest.ID = r.getNextAutoincrementID()

	r.lock.Lock()
	defer r.lock.Unlock()
	queue := r.matchRequests[matchRequest.Duration]
	r.matchRequests[matchRequest.Duration] = append(queue, matchRequest)
	r.allStoredEvents <- matchRequest

	return nil
}

func (r *MatchRequests) FindByUserID(userID int) (core.MatchRequest, error) {
	r.lock.RLock()
	defer r.lock.RUnlock()

	// Find first request and return it
	for _, queue := range r.matchRequests {
		for _, matchRequest := range queue {
			if matchRequest.UserID == userID {
				return matchRequest, nil
			}
		}
	}
	// If the map of match requests is empty, just return empty
	return core.MatchRequest{}, nil
}

func (r *MatchRequests) find(ID int) (int, int) {
	for duration, queue := range r.matchRequests {
		for i, matchRequest := range queue {
			if matchRequest.ID == ID {
				return duration, i
			}
		}
	}

	return -1, -1
}

func (r *MatchRequests) FindAll() (map[int][]core.MatchRequest, error) {
	return r.matchRequests, nil
}

func (r *MatchRequests) Delete(id int) (deleted int, err error) {
	r.lock.RLock()
	defer r.lock.RUnlock()

	duration, idx := r.find(id)
	if idx != -1 {
		r.matchRequests[duration] = core.RemoveMatchRequest(r.matchRequests[duration], idx)
		return 1, nil
	}

	return 0, nil
}

func (r *MatchRequests) DeleteByUserID(userID int) (deleted int, err error) {
	matchMe, _ := r.FindByUserID(userID)
	return r.Delete(matchMe.ID)
}

func (r *MatchRequests) ListenForStore() (core.MatchRequest, error) {
	var request core.MatchRequest
	for {
		request = <-r.allStoredEvents
		r.lock.RLock()
		if duration, _ := r.find(request.ID); duration != -1 {
			// event is still fresh
			r.lock.RUnlock()
			break
		}
		// toss this event because the request has expired and wait for another
		r.lock.RUnlock()
	}
	return request, nil
}
