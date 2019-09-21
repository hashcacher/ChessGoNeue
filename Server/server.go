package main

import (
	"log"
	"net/http"
)

type Server struct {
	games   map[string]Game
	looking map[string]chan bool
}

func NewServer() *Server {
	s := &Server{
		make(map[string]Game),
		make(map[string]chan bool),
	}

	return s
}

func main() {
	s := NewServer()

	s.games = make(map[string]Game)
	s.looking = make(map[string]chan bool)

	// TODO add http.servermux with metrics/logging middleware
	http.HandleFunc("/v1/matchMe", s.matchMeHandler)
	http.HandleFunc("/v1/move", s.moveHandler)
	http.HandleFunc("/", s.webHandler)
	http.HandleFunc("/ding", s.dingHandler)

	log.Fatal(http.ListenAndServe(":8080", nil))
}
