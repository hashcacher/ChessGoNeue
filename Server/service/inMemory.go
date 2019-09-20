package service

import (
	"fmt"
	"net/http"
	"strconv"

	"github.com/hashcacher/ChessGoNeue/Server/v2/core"
)

// Webservice struct holds data to be injected for use in implementing the web service
type Webservice struct {
	gameInteractor         core.GameInteractor
	userInteractor         core.UserInteractor
	matchRequestInteractor core.MatchRequestInteractor
}

// NewWebservice takes in interactors and creates a new web service that will use those interactors to fetch
// and manipulate data
func NewWebservice(gameInteractor core.GameInteractor, userInteractor core.UserInteractor, matchRequestInteractor core.MatchRequestInteractor) Webservice {
	return Webservice{
		gameInteractor,
		userInteractor,
		matchRequestInteractor,
	}
}

// GetUser retrieves a user
func (service *Webservice) GetUser(w http.ResponseWriter, r *http.Request) {
	userId, _ := strconv.Atoi(r.FormValue("userId"))
	user, _ := service.userInteractor.FindById(userId)
	// Check for empty
	if user.Id == 0 {
		http.Error(w, "user not found", http.StatusNotFound)
		return
	}
	fmt.Fprintf(w, "%v", user)
}

// func (service *Webservice) matchMeHandler(w http.ResponseWriter, r *http.Request) {
// 	userId, _ := strconv.Atoi(r.FormValue("userId"))
// 	// Get and wait for a match
// 	// Return the game
// }
