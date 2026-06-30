package renderer

import "github.com/netxph/relic/internal/parser"

// Section holds commits of a single type grouped for rendering.
type Section struct {
	Type   string
	Label  string
	Commits []parser.Commit
}

// ReleaseData is the full model passed to templates.
type ReleaseData struct {
	Version        string
	Date           string // ISO 8601 (YYYY-MM-DD)
	From           string
	To             string
	Sections       []Section
	BreakingChanges []parser.Commit
}
