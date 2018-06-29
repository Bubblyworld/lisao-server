package uci

import (
	"errors"
	"log"
)

type Game struct {
	moves []string
}

// PlayGame starts both clients and has them play a game with each other.
// TODO(guy) support for start options, time management, can we have clients
// replay each other? I.e. reset their pipes and such.
func PlayGame(white, black *Client) (*Game, error) {
	defer func() {
		err := stopEngines(white, black)
		if err != nil {
			log.Printf("uci: Error stopping engines: %v", err)
		}
	}()

	if err := startEngines(white, black); err != nil {
		return nil, err
	}

	if err := handshakeEngines(white, black); err != nil {
		return nil, err
	}

	if err := readyEngines(white, black); err != nil {
		return nil, err
	}

	return playGame(white, black)
}

func playGame(white, black *Client) (*Game, error) {
	var game Game

	// TODO(guy) break on game over, not random number
	for moveNum := 1; moveNum < 10; moveNum++ {
		currentPlayer := white
		if moveNum%2 == 0 {
			currentPlayer = black
		}

		bestMove, err := currentPlayer.PlayFrom(game.moves)
		if err != nil {
			return nil, err
		}

		game.moves = append(game.moves, bestMove.Move)
	}

	return &game, nil
}

func startEngines(white, black *Client) error {
	return combineErrors(white.Start(), black.Start())
}

func stopEngines(white, black *Client) error {
	return combineErrors(white.Stop(), black.Stop())
}

func handshakeEngines(white, black *Client) error {
	return combineErrors(white.DoHandshake(), black.DoHandshake())
}

func readyEngines(white, black *Client) error {
	return combineErrors(white.EnsureReadiness(), black.EnsureReadiness())
}

func combineErrors(a, b error) error {
	if a == nil && b == nil {
		return nil
	}

	if a == nil {
		return b
	}

	if b == nil {
		return a
	}

	return errors.New(a.Error() + " // " + b.Error())
}
