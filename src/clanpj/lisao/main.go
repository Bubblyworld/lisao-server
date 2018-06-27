package main

import (
	"clanpj/lisao/cmd"
	"flag"
	"io"
	"log"
	"os"
)

// TODO generate this on a tag commit hook with go generate.
var versionString = "ponita-0.0.0"

func main() {
	flag.Parse()

	// state := NewState()
	// waitGroup := sync.WaitGroup{}
	//
	// waitGroup.Add(1)
	// go state.PollGithubForever(&waitGroup)
	//
	// waitGroup.Wait()

	logger := cmd.NewLogWriter()
	command := cmd.
		NewCommand("ls /Users/guy").
		LogTo(logger)

	stdIn := command.GetStdIn()
	stdOut := command.GetStdOut()
	go io.Copy(stdIn, os.Stdin)
	go io.Copy(os.Stdout, stdOut)

	err := command.Do()
	if err != nil {
		log.Printf("Error running command: %v", err)
	}
}
