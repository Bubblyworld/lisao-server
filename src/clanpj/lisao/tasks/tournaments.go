package tasks

import (
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
	tournament, ok := work.(Tournament)
	if !ok {
		return errors.New("tournaments: received wrong type of work, should be Tournament")
	}

	log.Print(tournament)

	return nil
}
