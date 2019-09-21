package main

func (s *Server) webHandler(w http.ResponseWriter, r *http.Request) {
	w.Write("Hello World!")
}

func (s *Server) dingHandler(w http.ResponseWriter, r *http.Request) {
	w.Write("Dong")
}
