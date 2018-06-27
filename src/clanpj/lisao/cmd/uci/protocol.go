package uci

// SendUCI sends the inital UCI handshake.
// According to the UCI protocol docs, after receiving the "uci" message, the
// engine must identify itself with the "id" and "option" commands, after which
// it must reply with the "uciok" message to acknowledge the uci mode.
func (c *Client) SendUCI() error {
	// TODO(guy) expect replies somehow...
	return c.sendMessage("uci")
}
