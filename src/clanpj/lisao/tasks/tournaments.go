package tasks

import (
	"clanpj/lisao/cmd/uci"
	"errors"
	"log"
)

type Tournament struct {
	pathToWhiteEngine string
	pathToBlackEngine string

	numberOfGames int
}

func NewTournament(pathToWhite, pathToBlack string, numGames int) Tournament {
	return Tournament{
		pathToBlackEngine: pathToBlack,
		pathToWhiteEngine: pathToWhite,
		numberOfGames:     numGames,
	}
}

func DoTournament(work interface{}) error {
	t, ok := work.(Tournament)
	if !ok {
		return errors.New("tournaments: received wrong type of work")
	}

	white := uci.NewClient(t.pathToWhiteEngine, nil)
	black := uci.NewClient(t.pathToBlackEngine, nil)

	var games []*uci.Game
	for i := 0; i < t.numberOfGames; i++ {
		game, err := uci.PlayGame(&white, &black)
		if err != nil {
			return err
		}

		games = append(games, game)
	}

	// TODO(guy) push the results into mysql in some format.
	for _, game := range games {
		log.Print(game.GetPGN())
	}

	return nil
}
