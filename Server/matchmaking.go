package main

import (
	"encoding/json"
	"net/http"
)

type MatchMeRequest struct {
	clientId string
}

type MatchMeResponse struct {
	err string
	haveMatch bool
	areWhite bool
}

func RespondOK(w http.ResponseWriter) {
	res := MatchMeResponse {
		"", false, false
	}

	w.Write(json.Marshal(&res))
}

func RespondErr(w http.ResponseWriter, err string) {
	res := MatchMeResponse {
		err, false, false
	}

	w.Write(json.Marshal(&res))
}

func RespondFound(w http.ResponseWriter) {
	res := MatchMeResponse {
		"", true, true // TODO ask game what color we are
	}

	w.Write(json.Marshal(&res))
}


func matchHandler(server *Server, w http.ResponseWriter, r *http.Request) {
	req, err := json.Unmarshal(r.GetBody())
	if err != nil {
		ResponseErr(err)
	}

	if looking[req.clientId] {
		if len(looking) > 1 {
			// Match found! Take the first 2 people
			ResponseFound(w)

			players := make([]string, 2)
			ct := 0
			for clientId, _ := range server.looking {
				delete(server.looking, clientId)
				players[ct++] = clientId
				if ct == 2 {
					break
				}
			}
		} else {
		}
	}
}
