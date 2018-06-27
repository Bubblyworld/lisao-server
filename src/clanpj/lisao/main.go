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
	if err := client.Start(); err != nil {
		log.Print(err)
	}

	if err := client.SendUCI(); err != nil {
		log.Print(err)
	}

	client.Stop()
}
