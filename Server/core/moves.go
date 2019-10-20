package core

import (
	"errors"
	"fmt"
	"time"
)

// MakeMove validates a user and then performs a move
func (i *GamesInteractor) MakeMove(secret string, gameID int, move string) error {
	user, err := i.users.FindBySecret(secret, "")
	if err != nil {
		return errors.New("--couldnt find user by that id")
	}

	game := i.getGameForUser(user.ID, gameID)
	if game == nil {
		return errors.New("--couldnt find game")
	}

	if game.WhiteTurn && game.WhiteUser != user.ID ||
		!game.WhiteTurn && game.BlackUser != user.ID {
		Debug(fmt.Sprintf("--Wrong turn for user %s for game: %+v\n", secret, game))

		return errors.New("--not your turn")
	}

	err = i.games.MakeMove(game, &user, move)
	if err != nil {
		return err
	}

	// Update clocks
	lagCompensation, _ := time.ParseDuration("250ms")
	if user.ID == game.WhiteUser {
		game.BlackTurnStarted = time.Now().Add(lagCompensation)
		game.WhiteLeft = game.WhiteLeft - time.Now().Sub(game.WhiteTurnStarted)
	} else {
		game.WhiteTurnStarted = time.Now().Add(lagCompensation)
		game.BlackLeft = game.BlackLeft - time.Now().Sub(game.BlackTurnStarted)
	}

	// If the player time's out, find out
	go i.games.ListenForTimeout(game, user.ID, i.CompleteGame)

	return nil
}
