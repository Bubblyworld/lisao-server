package uci

import (
	"fmt"
	"strings"
)

// TODO(guy) all these things should have timeout contexts, and we should
// log all errors, especially if there are false positives for ErrIncorrectlyFormatted
// I think this would make implementation of time control trivial.

// SendUCI sends the inital UCI handshake and waits for the id, config state.
func (c *Client) DoHandshake() error {
	err := c.sendMessage("uci")
	if err != nil {
		return err
	}

	uciOk := UciOKMsg{}

	for {
		line, err := c.GetLine()
		if err != nil {
			return err
		}

		// Throw input away until we get a ready ACK.
		// TODO(guy) support ID and Option messages here.
		err = c.ParseLine(line, &uciOk)
		if err == nil {
			return nil
		}
	}
}

// NewGame sends the ucinewgame alert to the engine.
func (c *Client) NewGame() error {
	return c.sendMessage("ucinewgame")
}

// EnsureReadiness sends the isready alert to the engine and waits for an ack.
func (c *Client) EnsureReadiness() error {
	err := c.sendMessage("isready")
	if err != nil {
		return err
	}

	readyOk := ReadyOKMsg{}

	for {
		line, err := c.GetLine()
		if err != nil {
			return err
		}

		// Throw input away until we get a ready ACK.
		err = c.ParseLine(line, &readyOk)
		if err == nil {
			return nil
		}
	}
}

// SetOption sets a UCI option in the engine.
// TODO(guy): keep a list of valid options and error if invalid is passed in
func (c *Client) SetOption(name, value string) error {
	msg := fmt.Sprintf("setoption name %s value %s", name, value)

	return c.sendMessage(msg)
}

// PlayFrom sets the engine's starting position for a search.
func (c *Client) SetPosition(startFEN string, moves []string) error {
	msg := "position"
	if startFEN == "startpos" {
		msg += " startpos"
	} else {
		msg += " fen " + startFEN
	}

	if len(moves) > 0 {
		msg += " moves " + strings.Join(moves, " ")
	}

	return c.sendMessage(msg)
}

type SearchOptions struct {
	Depth int
}

// Search instructs the engine to look for the best move in the position.
func (c *Client) Search(opts SearchOptions) (*BestMoveMsg, error) {
	msg := "go"
	if opts.Depth != 0 {
		msg += fmt.Sprintf(" depth %d", opts.Depth)
	}

	err := c.sendMessage(msg)
	if err != nil {
		return nil, err
	}

	bestMove := BestMoveMsg{}

	for {
		line, err := c.GetLine()
		if err != nil {
			return nil, err
		}

		// Throw input away until we get a best move.
		err = c.ParseLine(line, &bestMove)
		if err == nil {
			return &bestMove, nil
		}
	}
}
