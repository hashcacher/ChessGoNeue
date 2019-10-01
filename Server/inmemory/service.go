package inmemory

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
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

// GetUser retrieves a user
func (service *WebService) GetUser(w http.ResponseWriter, r *http.Request) {
	Secret := r.FormValue("secret")
	user, _ := service.usersInteractor.FindBySecret(Secret)
	// Check for empty
	if user.Secret == "" {
		http.Error(w, "user not found", http.StatusNotFound)
		return
	}
	fmt.Fprintf(w, "%v", user)
}

// CreateUser creates a user
func (service *WebService) CreateUser(w http.ResponseWriter, r *http.Request) {
	reqBody, err := ioutil.ReadAll(r.Body)
	if err != nil {
		fmt.Fprintf(w, "user data missing from body")
	}
	// Parse user from json
	var user core.User
	json.Unmarshal(reqBody, &user)

	Secret := r.FormValue("Secret")
	user, _ = service.usersInteractor.FindBySecret(Secret)
	// Check for empty
	if user.Secret == "" {
		http.Error(w, "user not found", http.StatusNotFound)
		return
	}
	fmt.Fprintf(w, "%v", user)
}

// MatchMeRequest is the format of the request sent to this endpoint
type MatchMeRequest struct {
	Secret string `json:"secret"`
}

// MatchMeResponse is the format of the response sent from this endpoint
type MatchMeResponse struct {
	Err      string `json:"err"`
	GameID   int    `json:"gameId"`
	AreWhite bool   `json:"areWhite"`
}

func (service *WebService) MatchMe(w http.ResponseWriter, r *http.Request) {
	decoder := json.NewDecoder(r.Body)
	var request MatchMeRequest
	err := decoder.Decode(&request)
	if err != nil {
		panic(err)
	}

	// Authenticate user
	user, err := service.usersInteractor.FindBySecret(request.Secret)
	// If server error, respond with 500
	if err != nil {
		// Respond
		resp := MatchMeResponse{
			Err: err.Error(),
		}
		json, _ := json.Marshal(resp)
		w.WriteHeader(500)
		w.Write(json)
		return
	}
	// If user is empty, respond 404
	if reflect.DeepEqual(user, core.User{}) {
		// Respond
		resp := MatchMeResponse{
			Err: "user not found",
		}
		json, _ := json.Marshal(resp)
		w.WriteHeader(404)
		w.Write(json)
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
		log.Printf("ERROR: %v", err)
		// Respond
		resp := MatchMeResponse{
			Err: err.Error(),
		}
		json, _ := json.Marshal(resp)
		w.WriteHeader(500)
		w.Write(json)
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

// func (service *Webservice) MatchMe(w http.ResponseWriter, r *http.Request) {
// 	Secret, _ := strconv.Atoi(r.FormValue("Secret"))
// 	// Get and wait for a match
// 	// Return the game
// }
