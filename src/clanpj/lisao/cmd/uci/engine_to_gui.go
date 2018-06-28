package uci

import (
	"errors"
	"strings"
)

var ErrWrongNumberOfArgs = errors.New("wrong number of args")
var ErrIncorrectlyFormatted = errors.New("incorrectly formatted")

// UciUnmarshaler is an interface implemented by types that can unmarshal
// themselves from a UCI engine message, which consists of an array of strings.
//  e.g:  "id lisao" => ["id", "lisao"]
type UciUnmarshaler interface {
	Unmarshal([]string) error
}

// GetLine waits for the first non-empty line from the client and returns it.
// TODO(guy) this is dangerous, should probably take a timeout context
func (c *Client) GetLine() (string, error) {
	for {
		line, err := c.getLine()
		if err != nil {
			return "", err
		}

		if line != "" {
			return line, nil
		}
	}
}

func (c *Client) ParseLine(line string, msg UciUnmarshaler) error {
	// The UCI protocol states that if a line is malformed, one should ignore
	// tokens until either a match can be found or a match is impossible. It's
	// possible to perform a matching in better than O(n^2) time, but I can't
	// be bothered to write a DP for this. We return the first error if no match
	// can be found as it's usually indicative of the problem.
	var err error
	args := strings.Fields(line)
	for len(args) > 0 {
		newErr := msg.Unmarshal(args)
		if newErr == nil {
			return nil
		}

		if err == nil {
			err = newErr
		}

		args = args[1:]
	}

	return err
}

type UciOKMsg struct{}

func (m *UciOKMsg) Unmarshal(args []string) error {
	if len(args) != 1 {
		return ErrWrongNumberOfArgs
	}

	if args[0] != "uciok" {
		return ErrIncorrectlyFormatted
	}

	return nil
}

type ReadyOKMsg struct{}

func (m *ReadyOKMsg) Unmarshal(args []string) error {
	if len(args) != 1 {
		return ErrWrongNumberOfArgs
	}

	if args[0] != "readyok" {
		return ErrIncorrectlyFormatted
	}

	return nil
}

type IDMsg struct {
	Name  string
	Value string
}

func (m *IDMsg) Unmarshal(args []string) error {
	if len(args) < 3 {
		return ErrWrongNumberOfArgs
	}

	if args[0] != "id" {
		return ErrIncorrectlyFormatted
	}

	m.Name = args[1]
	m.Value = strings.Join(args[2:], " ")

	return nil
}

type BestMoveMsg struct {
	Move      string
	Pondering string
}

func (m *BestMoveMsg) Unmarshal(args []string) error {
	if len(args) != 2 && len(args) != 4 {
		return ErrWrongNumberOfArgs
	}

	if args[0] != "bestmove" {
		return ErrIncorrectlyFormatted
	}

	if len(args) == 4 && args[2] != "ponder" {
		return ErrIncorrectlyFormatted
	}

	m.Move = args[1]
	if len(args) == 4 {
		m.Pondering = args[3]
	}

	return nil
}
