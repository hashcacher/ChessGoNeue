package inmemory

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"

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
	// TODO: notificaiton channels

	// Create in memory data stores
	games := NewGames(gamesMap)
	users := NewUsers()
	matchRequests := NewMatchRequests(matchRequestsMap)
	//
	gamesInterractor := core.NewGamesInteractor(&games, &users)
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
	user, _ := service.usersInteractor.FindBySecret(Secret)
	// Check for empty
	if user.Secret == "" {
		http.Error(w, "user not found", http.StatusNotFound)
		return
	}
	fmt.Fprintf(w, "%v", user)
}

// func (service *Webservice) MatchMe(w http.ResponseWriter, r *http.Request) {
// 	Secret, _ := strconv.Atoi(r.FormValue("Secret"))
// 	// Get and wait for a match
// 	// Return the game
// }
