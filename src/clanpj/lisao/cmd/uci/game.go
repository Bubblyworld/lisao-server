package uci

// A game keeps track of board state and timing, manages communication
// between the competing engines and makes sure nobody cheats.
type Game struct {
	white *Client
	black *Client

	moves []string
	// TODO(guy) time management
}
