// Package main is a tool for converting opening info and movelists into a
// CSV of opening names to FEN strings, for testing bots in various positions.
package main

import (
	"flag"
	"io"
	"log"
	"os"
	"strings"
)

var pathToPGN = flag.String("pgn_path", "", "Path to a file containing a list of PGNS to parse.")

type PGN struct {
	Eco       string
	Opening   string
	Variation string

	Moves []string
}

func main() {
	flag.Parse()

	fenFile, err := os.Open(*pathToPGN)
	if err != nil {
		log.Fatalf("Error opening \"%s\": %v", *pathToPGN, err)
	}
	defer fenFile.Close()

	pgnList, err := parsePGNs(fenFile)
	if err != nil {
		log.Fatalf("Error parsing PGNs: %v", err)
	}

	for _, pgn := range pgnList {
		log.Printf("%+v", pgn)
	}
}

const (
	NewSymbolState      = 1
	CloseCommentState   = 2
	ReadTagState        = 3
	ReadMoveNumberState = 4
	ReadFirstMoveState  = 5
	ReadSecondMoveState = 6
)

const numerics = "0123456789"
const alphas = "abcdefgh"

func parsePGNs(file *os.File) ([]PGN, error) {
	data := make([]byte, 1024)

	var err error
	var result []PGN
	var currentPGN PGN
	var currentState = NewSymbolState
	var runningBuffer string
	for {
		_, err = file.Read(data)
		if err != nil {
			break
		}

		for _, b := range data {
			switch currentState {
			case NewSymbolState:
				switch b {
				case '{':
					currentState = CloseCommentState

				case '*':
					result = append(result, currentPGN)
					currentPGN = PGN{}

				case '[':
					currentState = ReadTagState
				}

				if strings.Contains(numerics, string(b)) {
					currentState = ReadMoveNumberState
				}

			case CloseCommentState:
				if b == '}' {
					currentState = NewSymbolState
				}

			case ReadTagState:
				if b == ']' {
					currentState = NewSymbolState
					currentPGN = parseTag(currentPGN, runningBuffer)
					runningBuffer = ""
					continue
				}

				runningBuffer += string(b)

			case ReadMoveNumberState:
				if b == '.' {
					currentState = ReadFirstMoveState
				}

			case ReadFirstMoveState:
				fallthrough

			case ReadSecondMoveState:
				if b == ' ' {
					if runningBuffer != "" {
						currentPGN.Moves = append(currentPGN.Moves, strings.TrimSpace(runningBuffer))
						runningBuffer = ""

						if currentState == ReadFirstMoveState {
							currentState = ReadSecondMoveState
							continue
						}

						currentState = NewSymbolState
						continue
					}
				}

				if b == '*' {
					result = append(result, currentPGN)
					currentPGN = PGN{}
					currentState = NewSymbolState
					continue
				}

				runningBuffer += string(b)

			default:
				log.Print("Invalid state - ignorning till end of file.")
			}
		}
	}

	if err != nil && err != io.EOF {
		return nil, err
	}
	return result, nil
}

func parseTag(pgn PGN, tag string) PGN {
	parts := strings.Split(tag, "\"")
	label := strings.Trim(parts[0], " ")
	value := strings.Trim(parts[1], " ")

	switch label {
	case "ECO":
		pgn.Eco = value

	case "Opening":
		pgn.Opening = value

	case "Variation":
		pgn.Variation = value

	default:
		log.Printf("Unknown label: %s", label)
	}

	return pgn
}
