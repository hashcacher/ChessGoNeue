package interfaces

import (
	"net/http"
	"strconv"

	usecases "github.com/hashcacher/ChessGoServer/usecases"
)

// Webservice struct holds data to be injected for use in implementing the web service
type Webservice struct {
	GameInterractor         usecases.GameInterractor
	MatchRequestInterractor usecases.MatchRequestInterractor
}

func (service *Webservice) matchMeHandler(w http.ResponseWriter, r *http.Request) {
	userId, _ := strconv.Atoi(r.FormValue("userId"))
	// Get and wait for a match
	// Return the game
}
