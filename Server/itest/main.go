package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"github.com/hashcacher/ChessGoNeue/Server/v2/inmemory"
	"io/ioutil"
	"math/rand"
	"net/http"
	"os"
	"strconv"
	"sync"
)

var HOST string = os.Getenv("ITEST_HOST")

func setup() {
	if HOST == "" {
		HOST = "http://localhost:8080"
	}
}

func main() {
	setup()
	fmt.Println(0)
	gameID, firstPlayer := createGame("123", "456")
	fmt.Println(1)
	makeRandomMoves(10, firstPlayer, gameID, "123", "456")
	fmt.Println(2)
	board, _ := getBoard("123", gameID)
	fmt.Println(3)
	fmt.Println(board)
}
func randomMove() string {
	move := strconv.Itoa(rand.Intn(8)) + "," + strconv.Itoa(rand.Intn(8))
	move += ">" + strconv.Itoa(rand.Intn(8)) + "," + strconv.Itoa(rand.Intn(8))

	return move
}

func makeRandomMoves(nMoves int, curPlayer string, gameID int, secret1 string, secret2 string) error {
	var err error
	var receivingPlayer string
	for i := 0; i < nMoves; i++ {
		if curPlayer == secret1 {
			receivingPlayer = secret2
		} else {
			receivingPlayer = secret1
		}

		move := randomMove()

		// Make and receive moves in lockstep
		var wg sync.WaitGroup
		wg.Add(2)

		// Receive Move
		go func() {
			defer wg.Done()
			fmt.Println("BLAH0")
			receivedMove, _ := getMove(receivingPlayer, gameID, move)
			fmt.Printf("BLAH %+v\n", receivedMove)
			if receivedMove.Move != move {
				fmt.Printf("received move %s didnt match move %s", receivedMove, move)
			}
		}()

		// Make move
		go func() {
			defer wg.Done()
			for {
				makeMoveRes, err := makeMove(curPlayer, gameID, move)
				fmt.Printf("%+v\n", makeMoveRes)
				if err != nil || makeMoveRes.Success == true {
					break
				}
				move = randomMove()
			}
		}()
		wg.Wait()

		if curPlayer == secret1 {
			curPlayer = secret2
		} else {
			curPlayer = secret1
		}
	}

	return err
}

func getMove(secret string, gameID int, move string) (inmemory.GetMoveResponse, error) {
	reqJSON, _ := json.Marshal(inmemory.MoveRequest{
		Secret: secret,
		GameID: gameID,
	})

	req, err := http.Post(HOST+"/v1/getMove", "application/json", bytes.NewBuffer(reqJSON))
	bodyBytes, err := ioutil.ReadAll(req.Body)

	var res inmemory.GetMoveResponse
	err = json.Unmarshal(bodyBytes, &res)

	return res, err
}

func makeMove(secret string, gameID int, move string) (inmemory.MakeMoveResponse, error) {
	reqJSON, _ := json.Marshal(inmemory.MoveRequest{
		Secret: secret,
		GameID: gameID,
		Move:   move,
	})

	req, err := http.Post(HOST+"/v1/makeMove", "application/json", bytes.NewBuffer(reqJSON))
	bodyBytes, err := ioutil.ReadAll(req.Body)

	var res inmemory.MakeMoveResponse
	err = json.Unmarshal(bodyBytes, &res)

	if err != nil {
		return inmemory.MakeMoveResponse{}, err
	}

	return res, nil
}

func getBoard(secret string, gameID int) (string, error) {
	reqJSON, _ := json.Marshal(inmemory.MoveRequest{
		Secret: secret,
		GameID: gameID,
	})

	req, err := http.Post(HOST+"/v1/getBoard", "application/json", bytes.NewBuffer(reqJSON))
	bodyBytes, err := ioutil.ReadAll(req.Body)

	var res struct {
		Board string
	}
	err = json.Unmarshal(bodyBytes, &res)

	return res.Board, err
}

func createGame(secret1 string, secret2 string) (int, string) {
	request1, _ := json.Marshal(inmemory.MatchMeRequest{
		Secret: secret1,
	})
	request2, _ := json.Marshal(inmemory.MatchMeRequest{
		Secret: secret2,
	})

	var wg sync.WaitGroup
	wg.Add(2)

	var response1, response2 inmemory.MatchMeResponse
	go func() {
		defer wg.Done()
		res, _ := http.Post(HOST+"/v1/matchMe", "application/json", bytes.NewBuffer(request1))
		bodyBytes, _ := ioutil.ReadAll(res.Body)
		json.Unmarshal(bodyBytes, &response1)
	}()
	go func() {
		defer wg.Done()
		res, err := http.Post(HOST+"/v1/matchMe", "application/json", bytes.NewBuffer(request2))
		if err != nil {
			fmt.Println(err)
		}
		bodyBytes, _ := ioutil.ReadAll(res.Body)
		json.Unmarshal(bodyBytes, &response2)
	}()

	fmt.Printf("waiting...")
	wg.Wait()
	fmt.Printf("res 1: %+v", response1)
	fmt.Printf("res 2: %+v", response2)

	var firstPlayer string
	if response1.AreWhite {
		firstPlayer = secret1
	} else {
		firstPlayer = secret2
	}

	return response1.GameID, firstPlayer
}
