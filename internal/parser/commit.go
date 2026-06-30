package parser

// Commit holds the parsed data from a single conventional commit.
type Commit struct {
	Hash               string
	Type               string
	Scope              string
	Description        string
	Breaking           bool
	BreakingDescription string
	Body               string
}
