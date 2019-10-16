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
	gamesMap := make(map[int]*core.Game)

	// Create in memory data stores
	games := NewGames(gamesMap)
	users := NewUsers()
	matchRequests := NewMatchRequests()
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
	Secret   string `json:"secret"`
	Name     string `json:"username"`
	Duration int    `json:"duration"`
}

// MatchMeRequest is the format of the request sent to this endpoint
type GetBoardRequest struct {
	Secret string `json:"secret"`
	GameID int    `json:"gameID"`
}

// MatchMeResponse is the format of the response sent from this endpoint
type MatchMeResponse struct {
	Err      string `json:"err"`
	GameID   int    `json:"gameID"`
	AreWhite bool   `json:"areWhite"`
	OppName  string `json:"oppName"`
	Duration string `json:"duration"`
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
	core.Game
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
	core.Debug(fmt.Sprintf("Matchme request %+v", request))

	// Authenticate user
	user, err := service.usersInteractor.FindBySecret(request.Secret, request.Name)
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

	// If user disconnects, dequeue
	closeNotify := w.(http.CloseNotifier).CloseNotify()
	go func(secret string) {
		<-closeNotify
		service.matchRequestsInteractor.DeleteMatchMe(user.ID)
		core.Debug(fmt.Sprintf("closeConnection: user %d", user.ID))
	}(request.Secret)

	core.Debug(fmt.Sprintf("Matchme request for user %d", user.ID))

	// Wait for match
	game, err := service.matchRequestsInteractor.MatchMe(user.ID, request.Duration)
	if err != nil {
		log.Printf("MATCHME ERROR: %v", err)
		w.WriteHeader(500)
		w.Write(errorResponse(err))
		return
	}

	// What's our opponemts name
	oppID := game.BlackUser
	if oppID == user.ID {
		oppID = game.WhiteUser
	}
	opp, _ := service.usersInteractor.FindByID(oppID)
	oppName := opp.Name
	core.Debug(fmt.Sprintf("Opponent name: %s", oppName))

	// Respond
	resp := MatchMeResponse{
		Err:      "",
		GameID:   game.ID,
		AreWhite: game.WhiteUser == user.ID,
		OppName:  oppName,
	}
	json, _ := json.Marshal(resp)
	w.Write(json)

	core.Debug("MatchMe succeeded")
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

	core.Debug(fmt.Sprintf("getboard request for user %s gameid %d", request.Secret, request.GameID))

	board, err := service.gamesInteractor.GetBoard(request.Secret, request.GameID)
	if err != nil {
		w.WriteHeader(500)
		w.Write(errorResponse(errors.New("error getting board " + err.Error())))
		return
	}

	for _, row := range board {
		w.Write(append(row[:], '\n'))
	}

	core.Debug("Getboard succeeded")
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

	core.Debug(fmt.Sprintf("MakeMove request for user %s gameid %d move %s", request.Secret, request.GameID, request.Move))

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

	core.Debug("MakeMove succeeded")
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

	core.Debug(fmt.Sprintf("GetMove request for user %s gameid %d", request.Secret, request.GameID))

	move, game, err := service.gamesInteractor.GetMove(request.Secret, request.GameID)
	if err != nil {
		w.WriteHeader(400)
		w.Write(errorResponse(err))
		return
	}

	resp := GetMoveResponse{
		Move: move,
		Game: *game,
	}
	json, _ := json.Marshal(resp)
	w.Write([]byte(json))

	core.Debug(fmt.Sprintf("GetMove succeeded: %+v", resp))
}
