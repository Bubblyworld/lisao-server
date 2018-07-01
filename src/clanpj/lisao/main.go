package main

import (
	"clanpj/lisao/tasks"
	"flag"
	"sync"
)

// TODO generate this on a tag commit hook with go generate.
var versionString = "ponita-0.0.0"

var testFENS = []string{
	"rnbqkbnr/pp2pppp/2p5/3p2P1/8/8/PPPPPPBP/RNBQK1NR b KQkq - 0 3",         // Grob, spike attack
	"rnbqkbnr/pppp1B1p/8/8/4Ppp1/5N2/PPPP2PP/RNBQK2R b KQkq - 0 5",          // KGA, wild muzio gambit
	"rnbqkbnr/pp1p1ppp/2p5/8/2B1Pp2/8/PPPP2PP/RNBQK1NR w KQkq - 0 4",        // KGA, ruy lopez defence
	"r1bqk2r/2ppbppp/p1n2n2/1p2p3/P3P3/1B3N2/1PPP1PPP/RNBQ1RK1 b kq a3 0 7", // RL, wing attack
	"rnbqkbnr/pppp2pp/5p2/4p3/4P3/5N2/PPPP1PPP/RNBQKB1R w KQkq - 0 3",       // Damiano defence. (oh god)
}

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

	for _, fen := range testFENS {
		gameConfig := tasks.GameConfig{fen}
		tournament.AddGame(gameConfig)
	}

	state.tournamentPool.PushWork(tournament)

	waitGroup.Wait()
}

func runPool(pool *tasks.Pool, waitGroup *sync.WaitGroup) {
	// This should never happen.
	defer waitGroup.Done()

	pool.Run()
}
