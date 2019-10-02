package inmemory

import (
	"encoding/json"
	"errors"
	"fmt"
	"log"
	"net/http"
	"reflect"

	"github.com/hashcacher/ChessGoNeue/Server/v2/core"
)

// WebService struct holds data to be injected for use in implementing the web service
type WebService struct {
	gamesInteractor         core.GamesInteractor
	usersInteractor         core.UsersInteractor
	matchRequestsInteractor core.MatchRequestsInteractor
}

// NewWebService takes in interactors and creates a new web service that will use those interactors to fetch
// and manipulate data
func NewWebService() WebService {
	// Create some context for inmemory data stores
	gamesMap := make(map[int]core.Game)
	matchRequestsMap := map[int]core.MatchRequest{}

	// Create in memory data stores
	games := NewGames(gamesMap)
	users := NewUsers()
	matchRequests := NewMatchRequests(matchRequestsMap)
	gamesInterractor := core.NewGamesInteractor(&games, &users, &matchRequests)
	go gamesInterractor.StartGameCreateDaemon()
	usersInteractor := core.NewUsersInteractor(&users)
	matchRequestsInteractor := core.NewMatchRequestsInteractor(&matchRequests, &users, &games)

	return WebService{
		gamesInterractor,
		usersInteractor,
		matchRequestsInteractor,
	}
}

// MatchMeRequest is the format of the request sent to this endpoint
type MatchMeRequest struct {
	Secret string `json:"secret"`
}

// MatchMeRequest is the format of the request sent to this endpoint
type GetBoardRequest struct {
	Secret string `json:"secret"`
	GameID int    `json:"gameID"`
}

// MatchMeResponse is the format of the response sent from this endpoint
type MatchMeResponse struct {
	Err      string `json:"err"`
	GameID   int    `json:"gameId"`
	AreWhite bool   `json:"areWhite"`
}

// MoveRequest is the format of the request sent for GetMove/MakeMove
type MoveRequest struct {
	Secret string `json:"secret"`
	GameID int    `json:"gameID"`
	Move   string `json:"move"`
}

// MoveResponse is the format of the response sent from this endpoint
type MakeMoveResponse struct {
	Success bool `json:"success"`
}

// MoveResponse is the format of the response sent from this endpoint
type GetMoveResponse struct {
	Move string `json:"move"`
}

func errorResponse(err error) []byte {
	return []byte(fmt.Sprintf(`{ "err": "%s" }`, err.Error()))
}

func (service *WebService) MatchMe(w http.ResponseWriter, r *http.Request) {
	decoder := json.NewDecoder(r.Body)
	var request MatchMeRequest
	err := decoder.Decode(&request)
	if err != nil {
		w.WriteHeader(400)
		w.Write(errorResponse(err))
		return
	}

	// Authenticate user
	user, err := service.usersInteractor.FindBySecret(request.Secret)
	if err != nil {
		w.WriteHeader(500)
		w.Write(errorResponse(err))
		return
	}

	// If user is empty, respond 404
	if reflect.DeepEqual(user, core.User{}) {
		w.Write(errorResponse(errors.New("user not found")))
		w.WriteHeader(404)
		return
	}

	/*
		closeNotify := w.(http.CloseNotifier).CloseNotify()
		go func(secret string) {
			<-closeNotify
			// Cleanup
			service.matchRequestsInteractor.Delete(secret)
			// Need to close the connection and cleanup the rest of the objects
		}(matchMeReq.secret)
	*/

	// Wait for match
	game, err := service.matchRequestsInteractor.MatchMe(user.ID)
	if err != nil {
		log.Printf("MATCHME ERROR: %v", err)
		w.WriteHeader(500)
		w.Write(errorResponse(err))
		return
	}

	// Respond
	resp := MatchMeResponse{
		Err:      "",
		GameID:   game.ID,
		AreWhite: game.WhiteUser == user.ID,
	}
	json, _ := json.Marshal(resp)
	w.WriteHeader(500)
	w.Write(json)

}

func (service *WebService) GetBoard(w http.ResponseWriter, r *http.Request) {
	decoder := json.NewDecoder(r.Body)
	var request GetBoardRequest
	err := decoder.Decode(&request)
	if err != nil {
		w.WriteHeader(400)
		w.Write(errorResponse(errors.New("invalid arguments")))
		return
	}

	board, err := service.gamesInteractor.GetBoard(request.Secret, request.GameID)
	if err != nil {
		w.WriteHeader(500)
		w.Write(errorResponse(errors.New("error getting board " + err.Error())))
		return
	}

	for _, row := range board {
		w.Write(append(row[:], '\n'))
	}
}

func (service *WebService) MakeMove(w http.ResponseWriter, r *http.Request) {
	decoder := json.NewDecoder(r.Body)
	var request MoveRequest
	err := decoder.Decode(&request)
	if err != nil {
		w.WriteHeader(400)
		w.Write(errorResponse(errors.New("invalid arguments")))
		return
	}

	err = service.gamesInteractor.MakeMove(request.Secret, request.GameID, request.Move)
	if err != nil {
		w.WriteHeader(400)
		w.Write(errorResponse(err))
		return
	}

	resp := MakeMoveResponse{
		Success: true,
	}
	json, _ := json.Marshal(resp)
	w.Write([]byte(json))
}

func (service *WebService) GetMove(w http.ResponseWriter, r *http.Request) {
	decoder := json.NewDecoder(r.Body)
	var request MoveRequest
	err := decoder.Decode(&request)
	if err != nil {
		w.WriteHeader(400)
		w.Write([]byte("invalid arguments"))
		return
	}

	move, err := service.gamesInteractor.GetMove(request.Secret, request.GameID)
	if err != nil {
		w.WriteHeader(400)
		w.Write(errorResponse(err))
		return
	}

	resp := GetMoveResponse{move}
	json, _ := json.Marshal(resp)
	w.Write([]byte(json))
}
