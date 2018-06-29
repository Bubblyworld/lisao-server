package main

import (
	"clanpj/lisao/cmd/uci"
	"log"
)

// TODO generate this on a tag commit hook with go generate.
var versionString = "ponita-0.0.0"

func main() {
	// flag.Parse()
	//
	// state := NewState()
	// waitGroup := sync.WaitGroup{}
	//
	// waitGroup.Add(1)
	// go state.PollGithubForever(&waitGroup)
	//
	// waitGroup.Wait()

	white := uci.NewClient("/Users/guy/Workspace/lisao-bot/bin/uci", nil)
	black := uci.NewClient("/Users/guy/Workspace/lisao-bot/bin/uci", nil)

	game, err := uci.PlayGame(&white, &black)
	if err != nil {
		log.Fatal(err)
	}

	log.Printf("Moves: %+v", game)

	game, err = uci.PlayGame(&white, &black)
	if err != nil {
		log.Fatal(err)
	}

	log.Printf("Moves: %+v", game)
}
