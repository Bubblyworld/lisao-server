package uci

import (
	"errors"
	"log"

	"github.com/notnil/chess"
)

type Outcome int

const (
	NoOutcome Outcome = 0
	WhiteWon  Outcome = 1
	BlackWon  Outcome = 2
	Draw      Outcome = 3
)

func convertOutcome(o chess.Outcome) Outcome {
	switch o {
	case chess.WhiteWon:
		return WhiteWon

	case chess.BlackWon:
		return BlackWon

	case chess.Draw:
		return Draw

	default:
		return NoOutcome
	}
}

type Game struct {
	startFEN string
	moves    []string
	Outcome  Outcome
}

func convertGame(g Game) (*chess.Game, error) {
	chessFEN, err := chess.FEN(g.startFEN)
	if err != nil {
		return nil, err
	}

	game := chess.NewGame(chessFEN)
	chess.UseNotation(chess.LongAlgebraicNotation{})(game)

	for _, move := range g.moves {
		err := game.MoveStr(move)
		if err != nil {
			return nil, err
		}
	}

	return game, nil
}

func (g Game) GetLatestFEN() (string, error) {
	game, err := convertGame(g)
	if err != nil {
		return "", err
	}

	return game.FEN(), nil
}

func (g Game) GetPGN() (string, error) {
	return convertGameToPGN(g)
}

// PlayGame starts both clients and has them play a game with each other.
// TODO(guy) support for start options, time management, etc
func PlayGame(white, black *Client, startFEN string) (*Game, error) {
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

	return playGame(white, black, startFEN)
}

func playGame(white, black *Client, startFEN string) (*Game, error) {
	var game Game
	game.startFEN = startFEN

	chessGame, err := convertGame(game)
	if err != nil {
		return nil, err
	}

	for {
		currentPlayer := white
		if chessGame.Position().Turn() == chess.Black {
			currentPlayer = black
		}

		bestMove, err := currentPlayer.PlayFrom(startFEN, game.moves)
		if err != nil {
			return nil, err
		}

		// Validate the move and check endgame conditions.
		// TODO(guy) support draws
		err = chessGame.MoveStr(bestMove.Move)
		if err != nil {
			return nil, err
		}

		game.moves = append(game.moves, bestMove.Move)

		if chessGame.Outcome() != chess.NoOutcome {
			break
		}
	}

	game.Outcome = convertOutcome(chessGame.Outcome())
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
