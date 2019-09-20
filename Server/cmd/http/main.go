package main

import (
	"log"
	"net/http"

	"github.com/hashcacher/ChessGoNeue/Server/v2/core"
	"github.com/hashcacher/ChessGoNeue/Server/v2/implementations"
)

func main() {

	gameRepository := implementations.IMGameRepository{}
	users := make(map[int]core.User)
	users[1] = core.User{
		Id:       1,
		Username: "zac",
	}
	userRepository := implementations.NewIMUserRepIMUserRepository(users)
	matchRequestRepository := implementations.IMMatchRequestRepository{}

	gameInteractor := core.GameInteractor{&userRepository, &gameRepository}
	userInteractor := core.UserInteractor{&userRepository}
	matchRequestInteractor := core.MatchRequestInteractor{matchRequestRepository, &userRepository, &gameRepository}

	s := implementations.NewWebservice(gameInteractor, userInteractor, matchRequestInteractor)

	// TODO add http.servermux with metrics/logging middleware
	http.HandleFunc("/v1/getUser", s.GetUser)
	// http.HandleFunc("/v1/matchMe", s.matchMeHandler)
	// http.HandleFunc("/v1/move", s.moveHandler)

	log.Fatal(http.ListenAndServe(":8080", nil))
}
