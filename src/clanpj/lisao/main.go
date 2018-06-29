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

	client := uci.NewClient("/Users/guy/Workspace/lisao-bot/bin/uci")
	defer client.Stop()

	if err := client.Start(); err != nil {
		log.Print(err)
	}

	if err := client.DoHandshake(); err != nil {
		log.Print(err)
	}

	if err := client.NewGame(); err != nil {
		log.Print(err)
	}

	if err := client.EnsureReadiness(); err != nil {
		log.Print(err)
	}

	// Should be ready to play! Let's go.
	moves := []string{}
	for i := 0; i < 10; i++ {
		bestMove, err := client.PlayFrom(moves)
		if err != nil {
			log.Print(err)
			return
		}

		moves = append(moves, bestMove.Move)
	}

	log.Print()
	log.Print(moves)

	// msg, err := client.GetLine()
	// if err != nil {
	// 	log.Printf("Error getting message: %v", err)
	// }
	//
	// fmt.Println(msg)
}
