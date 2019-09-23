package inmemory

import (
	"fmt"
	"net/http"

	"github.com/hashcacher/ChessGoNeue/Server/v2/core"
)

// Webservice struct holds data to be injected for use in implementing the web service
type Webservice struct {
	gamesInteractor         core.GamesInteractor
	usersInteractor         core.UsersInteractor
	matchRequestsInteractor core.MatchRequestsInteractor
}

// NewWebservice takes in interactors and creates a new web service that will use those interactors to fetch
// and manipulate data
func NewWebservice(gameInteractor core.GamesInteractor, userInteractor core.UsersInteractor, matchRequestInteractor core.MatchRequestsInteractor) Webservice {
	return Webservice{
		gameInteractor,
		userInteractor,
		matchRequestInteractor,
	}
}

// GetUser retrieves a user
func (service *Webservice) GetUser(w http.ResponseWriter, r *http.Request) {
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
