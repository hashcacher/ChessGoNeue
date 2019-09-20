package inmemory

import "github.com/hashcacher/ChessGoNeue/Server/v2/core"

type MatchRequests struct {
	matchRequests map[int]core.MatchRequest
}

func NewMatchRequests(matchRequests map[int]core.MatchRequest) MatchRequests {
	return MatchRequests{
		matchRequests,
	}
}

func (r MatchRequests) Store(matchRequest core.MatchRequest) error {
	r.matchRequests[matchRequest.Id] = matchRequest
	return nil
}

func (r MatchRequests) FindAllMatchRequestsByUserId(userId int) []core.MatchRequest {
	// TODO
	return []core.MatchRequest{}
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
