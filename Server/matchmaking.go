package main

import (
	"encoding/json"
	"io/ioutil"
	"log"
	"net/http"
)

// MatchMeRequest is the format of the request sent to this endpoint
type MatchMeRequest struct {
	ClientID string `json:"clientID"`
}

// MatchMeResponse is the format of the response sent from this endpoint
type MatchMeResponse struct {
	Err       error `json:"err"`
	HaveMatch bool  `json:"haveMatch"`
	AreWhite  bool  `json:"areWhite"`
}

// RespondOK means OK
func RespondOK(w http.ResponseWriter) {
	res := MatchMeResponse{
		nil, false, false,
	}

	resJSON, err := json.Marshal(res)

	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte("Unknown server error occurred"))
		return
	}
	w.Write(resJSON)
	log.Printf("Responded with ok ")
}

// RespondErr means not OK
func RespondErr(w http.ResponseWriter, err error) {
	res := MatchMeResponse{
		err, false, false,
	}

	resJSON, err := json.Marshal(res)

	if err != nil {
		w.Write([]byte("Unknown server error occurred"))
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	w.WriteHeader(http.StatusInternalServerError)
	w.Write(resJSON)
	log.Printf("Responded with err ")
}

// RespondFound means match found!
func RespondFound(w http.ResponseWriter) {
	res := MatchMeResponse{
		nil, true, true, // TODO ask game what color we are
	}

	resJSON, err := json.Marshal(res)

	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte("Unknown server error occurred"))
		return
	}
	w.Write(resJSON)
	log.Printf("Responded with found ")
}

func (s *Server) matchMeHandler(w http.ResponseWriter, r *http.Request) {
	bodyBytes, err := ioutil.ReadAll(r.Body)
	if err != nil {
		RespondErr(w, err)
		return
	}

	var matchMeReq MatchMeRequest
	err = json.Unmarshal(bodyBytes, &matchMeReq)
	if err != nil {
		RespondErr(w, err)
		return
	}

	// If the client is already registered
	if s.looking[matchMeReq.ClientID] != nil {
		RespondOK(w) // Just ignore it
	} else {
		// Create channel that will get response
		resChan := make(chan bool)
		finishChan := make(chan bool)
		s.looking[matchMeReq.ClientID] = resChan

		// Start a coroutine that sends a response once ready
		go func(responseWriter http.ResponseWriter, resChan chan bool, finishChan chan bool) {
			<-resChan
			RespondFound(responseWriter)
			finishChan <- true
		}(w, resChan, finishChan)

		// Remove the channel if conn closes
		// XXX doesnt seem to fire
		closeNotify := w.(http.CloseNotifier).CloseNotify()
		go func(clientID string, finishChan chan bool) {
			<-closeNotify
			delete(s.looking, clientID)
			finishChan <- true
		}(matchMeReq.ClientID, finishChan)

		// Check if we have 2 people to match
		if len(s.looking) > 1 {
			ct := 0
			for clientID, resChan := range s.looking {
				// Drop two messages in the channel:
				// One to complete the handler and one for the goroutine
				resChan <- true
				ct++

				delete(s.looking, clientID)

				if ct == 2 {
					break
				}
			}
		}
		<-finishChan // Wait for something to trigger finish
	}
}
