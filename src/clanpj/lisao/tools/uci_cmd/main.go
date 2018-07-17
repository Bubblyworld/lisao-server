// Package main is a command-line utility for testing UCI engines. The main
// purpose of this is to enable refactoring of the Lisao engine by automatically
// checking that the evaluation of a group of positions remains the same.
package main

import (
	"bufio"
	"flag"
	"io"
	"log"
	"strings"

	"clanpj/lisao/cmd/uci"
)

var enginePath = flag.String("engine", "", "Path the the UCI engine binary.")
var startingFEN = flag.String("fen", "", "Initial FEN position for the engine.")
var options = flag.String("options", "depth=7", "Comma-separated list of engine options to set, in the format \"OP1=VAL1,OP2=VAL2\".")

func main() {
	flag.Parse()

	engine := uci.NewClient(*enginePath, NewInfoLogger())
	fatalOnErr(engine.Start)
	fatalOnErr(engine.DoHandshake)
	fatalOnErr(engine.EnsureReadiness)

	for _, option := range strings.Split(*options, ",") {
		option = strings.TrimSpace(option)
		optionTokens := strings.Split(option, "=")

		if len(optionTokens) != 2 {
			log.Fatal("Options should be specified in the format \"OP1=VAL1,OP2=VAL2\".")
		}

		name := strings.TrimSpace(optionTokens[0])
		value := strings.TrimSpace(optionTokens[1])
		fatalOnErr(setOptionFn(engine, name, value))
	}

	fatalOnErr(engine.EnsureReadiness)
	fatalOnErr(playFromFn(engine, *startingFEN))
	fatalOnErr(engine.Stop)
}

func fatalOnErr(fn func() error) {
	if err := fn(); err != nil {
		log.Fatal(err.Error())
	}
}

func setOptionFn(engine uci.Client, name, value string) func() error {
	return func() error {
		return engine.SetOption(name, value)
	}
}

func playFromFn(engine uci.Client, fen string) func() error {
	return func() error {
		_, err := engine.PlayFrom(fen, nil)

		return err
	}
}

// InfoLogger is an io.Writer that logs UCI "info" outputs.
type InfoLogger struct {
	pipeReader *io.PipeReader
	pipeWriter *io.PipeWriter
	doneCh     chan error
}

func NewInfoLogger() InfoLogger {
	r, w := io.Pipe()

	infoLogger := InfoLogger{
		pipeReader: r,
		pipeWriter: w,
		doneCh:     make(chan error),
	}

	go infoLogger.copy()
	return infoLogger
}

func (lw InfoLogger) copy() {
	bufReader := bufio.NewReader(lw.pipeReader)

	for {
		line, err := bufReader.ReadString('\n')
		if err != nil {
			lw.doneCh <- err
			return
		}

		tokens := strings.Fields(line)
		if len(tokens) > 0 && tokens[0] == "info" {
			log.Print(line)
		}
	}
}

func (lw InfoLogger) Write(data []byte) (int, error) {
	return lw.pipeWriter.Write(data)
}

func (lw InfoLogger) Close() error {
	err := lw.pipeWriter.Close()
	if err != nil {
		return err
	}

	return <-lw.doneCh
}
