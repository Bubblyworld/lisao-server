package tasks

import (
	"clanpj/lisao/cmd/uci"
	"errors"
	"log"
)

type Tournament struct {
	pathToWhiteEngine string
	pathToBlackEngine string

	gameConfigs []GameConfig
}

type GameConfig struct {
	StartFEN string
}

func NewTournament(pathToWhite, pathToBlack string) Tournament {
	return Tournament{
		pathToBlackEngine: pathToBlack,
		pathToWhiteEngine: pathToWhite,
	}
}

func (t *Tournament) AddGame(game GameConfig) {
	t.gameConfigs = append(t.gameConfigs, game)
}

func DoTournament(work interface{}) error {
	t, ok := work.(Tournament)
	if !ok {
		return errors.New("tournaments: received wrong type of work")
	}

	white := uci.NewClient(t.pathToWhiteEngine, nil)
	black := uci.NewClient(t.pathToBlackEngine, nil)

	var games []*uci.Game
	for _, config := range t.gameConfigs {
		game, err := uci.PlayGame(&white, &black, config.StartFEN)
		if err != nil {
			return err
		}

		games = append(games, game)
	}

	// TODO(guy) push the results into mysql in some format.
	for _, game := range games {
		pgn, err := game.GetPGN()
		if err != nil {
			return err
		}

		log.Print(pgn)
	}

	return nil
}
