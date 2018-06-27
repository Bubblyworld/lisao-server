package main

import (
	"bufio"
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
	defer logger.Close()

	stdIn := bytes.NewBuffer([]byte{})

	command := cmd.NewCommand("/Users/guy/Workspace/lisao-bot/bin/uci")
	command.Stdout = logger
	command.Stderr = logger
	command.Stdin =

	bufStdIn := bufio.NewWriter(stdIn)
	bufStdIn.WriteString("uci\n")
	bufStdIn.Flush()

	err := command.Do()
	if err != nil {
		log.Print(err)
	}
}
