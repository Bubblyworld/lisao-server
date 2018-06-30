package main

import (
	"clanpj/lisao/tasks"
	"flag"
	"sync"
)

// TODO generate this on a tag commit hook with go generate.
var versionString = "ponita-0.0.0"

func main() {
	flag.Parse()

	state := NewState()
	waitGroup := sync.WaitGroup{}

	waitGroup.Add(3)
	go state.PollGithubForever(&waitGroup)
	go runPool(state.buildsPool, &waitGroup)
	go runPool(state.tournamentPool, &waitGroup)

	engine := "/Users/guy/Workspace/lisao-bot/bin/uci"
	tournament := tasks.NewTournament(engine, engine)
	gameConfig := tasks.GameConfig{
		StartFEN: "rnbqkbnr/pppppppp/8/8/P7/8/1PPPPPPP/RNBQKBNR b KQkq a3 0 1", // Ware opening
	}

	tournament.AddGame(gameConfig)

	state.tournamentPool.PushWork(tournament)

	waitGroup.Wait()
}

func runPool(pool *tasks.Pool, waitGroup *sync.WaitGroup) {
	// This should never happen.
	defer waitGroup.Done()

	pool.Run()
}
