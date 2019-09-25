package main

import (
	"log"
	"net/http"

	inmemory "github.com/hashcacher/ChessGoNeue/Server/v2/inmemory"
)

func main() {

	s := inmemory.NewWebService()

	// TODO add http.servermux with metrics/logging middleware
	http.HandleFunc("/v1/getUser", s.GetUser)
	http.HandleFunc("/v1/matchMe", s.MatchMe)
	// http.HandleFunc("/v1/move", s.moveHandler)

	log.Fatal(http.ListenAndServe(":8080", nil))
}
