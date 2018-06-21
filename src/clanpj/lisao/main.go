package main

import (
	"flag"
	"log"

	"clanpj/lisao/tasks"
)

// TODO generate this on a tag commit hook with go generate.
var versionString = "ponita-0.0.0"

func main() {
	flag.Parse()

	log.Printf("Starting Lisao DevServer %s.", versionString)

	buildInfo := tasks.NewBuildInfo("8ce3f9f32b2abfa1672158c1e0160f7eeb13cf2d",
		"clanpj/lisao/mains/lichess", "/Users/guy/lichess")

	err := tasks.DoBuild(buildInfo)
	log.Print(err)
}
