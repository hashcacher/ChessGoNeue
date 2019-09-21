package main

import (
	"log"
	"net/http"

	"github.com/hashcacher/ChessGoNeue/Server/v2/core"
	inmemory "github.com/hashcacher/ChessGoNeue/Server/v2/inmemory"
)

func main() {

	// Create in memory data stores
	games := inmemory.NewGames(map[int]core.Game{})
	users := inmemory.NewUsers()
	matchRequests := inmemory.NewMatchRequests(map[int]core.MatchRequest{})

	// Create interactors based on the data stores above
	gameInteractor := core.NewGamesInteractor(&games, &users)
	userInteractor := core.NewUsersInteractor(&users)
	matchRequestInteractor := core.NewMatchRequestsInteractor(matchRequests, &users, &games)

	s := inmemory.NewWebservice(gameInteractor, userInteractor, matchRequestInteractor)

	// TODO add http.servermux with metrics/logging middleware
	http.HandleFunc("/v1/getUser", s.GetUser)
	// http.HandleFunc("/v1/matchMe", s.matchMeHandler)
	// http.HandleFunc("/v1/move", s.moveHandler)

	log.Fatal(http.ListenAndServe(":8080", nil))
}
