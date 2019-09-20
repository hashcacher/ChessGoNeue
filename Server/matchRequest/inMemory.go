package matchRequest

import "github.com/hashcacher/ChessGoNeue/Server/v2/core"

type IMMatchRequestRepository struct {
	matchRequests map[int]core.MatchRequest
}

func (r IMMatchRequestRepository) Store(matchRequest core.MatchRequest) error {
	r.matchRequests[matchRequest.Id] = matchRequest
	return nil
}

func (r IMMatchRequestRepository) FindAllMatchRequestsByUserId(userId int) []core.MatchRequest {
	// TODO
	return []core.MatchRequest{}
}

func (r IMMatchRequestRepository) Delete(id int) (deleted int, err error) {
	_, ok := r.matchRequests[id]
	deleted = 0
	if ok {
		deleted = 1
	}
	delete(r.matchRequests, id)
	return deleted, nil
}
