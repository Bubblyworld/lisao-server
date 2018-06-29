package main

import (
	"flag"
	"sync"
)

// TODO generate this on a tag commit hook with go generate.
var versionString = "ponita-0.0.0"

func main() {
	flag.Parse()

	state := NewState()
	waitGroup := sync.WaitGroup{}

	waitGroup.Add(1)
	go state.PollGithubForever(&waitGroup)

	waitGroup.Wait()
}
