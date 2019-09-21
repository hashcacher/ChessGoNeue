package main

import "net/http"

func (s *Server) webHandler(w http.ResponseWriter, r *http.Request) {
	w.Write([]byte("Hello World!"))
}

func (s *Server) dingHandler(w http.ResponseWriter, r *http.Request) {
	w.Write([]byte("Dong"))
}
