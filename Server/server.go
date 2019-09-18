package main

import (
	"fmt"
	"log"
	"net/http"
)

type Server struct {
	games   map[string]Game
	looking map[string]bool
}

func main() {
	games = make(map[string]Game)
	looking = make(map[string]bool)

	http.HandleFunc("/v1/matchMe", matchHandler)
	http.HandleFunc("/v1/move", moveHandler)
	log.Fatal(http.ListenAndServe(":8080", nil))
}
