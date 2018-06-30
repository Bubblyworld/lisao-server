package uci

import (
	"fmt"

	"github.com/notnil/chess"
)

const startFEN = "rnbqkbnr/pppppppp/8/8/8/8/PPPPPPPP/RNBQKBNR w KQkq - 0 1"

// convertGameToPGN does what it says on the box - the chess libraries I'm using
// don't seem to handle arbitrary start positions well for generatings PGNs,
// which is probably because they don't store starting positions. So we'll just
// do it ourselves.
func convertGameToPGN(game Game) (string, error) {
	chessGame, err := convertGame(game)
	if err != nil {
		return "", err
	}

	moves := chessGame.Moves()
	positions := chessGame.Positions()

	pgn := getTags(positions[0])
	fullMoveIndex := 1
	halfMoveIndex := 0
	notation := chess.AlgebraicNotation{}

	// No moves, early out.
	if len(moves) == 0 {
		return pgn + "*", nil
	}

	// Handle starting on black's turn.
	if positions[0].Turn() == chess.Black {
		move := moves[halfMoveIndex]
		pgn += "1..." + notation.Encode(positions[0], move) + " "
		halfMoveIndex++
		fullMoveIndex++
	}

	for ; halfMoveIndex < len(moves); halfMoveIndex++ {
		move := moves[halfMoveIndex]
		position := positions[halfMoveIndex]

		if positions[halfMoveIndex].Turn() == chess.White {
			fullMoveStr := fmt.Sprint(fullMoveIndex)
			pgn += fullMoveStr + "." + notation.Encode(position, move) + " "
		} else {
			pgn += notation.Encode(position, move) + " "
			fullMoveIndex++
		}
	}

	return pgn + "*", nil
}

func getTags(startingPosition *chess.Position) string {
	fen := startingPosition.String()
	if fen == startFEN {
		return ""
	}

	return fmt.Sprintf("[SetUp \"1\"]\n[FEN \"%s\"]\n\n", fen)
}
