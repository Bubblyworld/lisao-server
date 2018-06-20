package main

import (
	"log"

	"clanpj/lisao/github"
)

// TODO generate this on a tag commit hook with go generate.
var versionString = "rapidash-0.0.0"

func main() {
	log.Printf("Starting Lisao DevServer %s.", versionString)

	githubClient := github.NewClient("Bubblyworld", "lisao-bot")
	refs, err := githubClient.GetRefs()
	if err != nil {
		log.Fatalf("Error getting all refs: %v", err)
	}

	log.Printf("%+v", refs)
}
