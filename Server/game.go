package main

import (
	"net/http"
)

type Game struct {
	board [8][8]byte
}

func (s *Server) moveHandler(w http.ResponseWriter, r *http.Request) {

}

func (g *Game) makeMove(move string) {
	//move.
}
