package uci

import (
	"strings"
	"time"
)

const uciTimeout = time.Second * 10

// TODO(guy) all these things should have timeout contexts.
// in particular, this would make implementation of time control trivial.

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

// PlayFrom sends the engine the current movelist and waits for it to make a
// move.
func (c *Client) PlayFrom(moves []string) (*BestMoveMsg, error) {
	err := c.sendMessage("position startpos moves " + strings.Join(moves, " "))
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
