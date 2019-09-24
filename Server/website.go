package main

import (
	"fmt"
	"net/http"
	"time"
)

var nextCheck time.Time

func (s *Server) webHandler(w http.ResponseWriter, r *http.Request) {
	w.Write([]byte(s.poorMansHTML()))
}

func (s *Server) poorMansHTML() string {
	baseURL := "https://storage.cloud.google.com/chessgo/"
	linuxURL := "master/linux_chessgo_master"
	osxURL := "master/osx-chessgo-master.app"

	return fmt.Sprintf(`Welcome to ChessGo. Here are our beta clients to try:<br>
	<a href="%s">ChessGo for Mac OSX</a><br>
	<a href="%s">ChessGo for Linux</a><br><br>
	Email chessgoinfo@gmail.com for more info.
	`, baseURL+linuxURL, baseURL+osxURL)
}

func (s *Server) dingHandler(w http.ResponseWriter, r *http.Request) {
	w.Write([]byte("dong"))
}
