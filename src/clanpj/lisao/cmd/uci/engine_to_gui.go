package uci

// TODO(guy) implement reply types and marshalling
func (c *Client) GetMessage() (string, error) {
	return c.getMessage()
}
