package parser

import (
	"regexp"
	"strings"

	"github.com/netxph/relic/internal/git"
)

var subjectRe = regexp.MustCompile(`^(\w+)(\([\w-]+\))?(!)?: (.+)$`)

const breakingFooter = "BREAKING CHANGE:"

// Parse converts raw commits into conventional commits, silently skipping non-conforming ones.
func Parse(raws []git.RawCommit) []Commit {
	var commits []Commit
	for _, raw := range raws {
		c, ok := parseOne(raw)
		if ok {
			commits = append(commits, c)
		}
	}
	return commits
}

func parseOne(raw git.RawCommit) (Commit, bool) {
	m := subjectRe.FindStringSubmatch(raw.Subject)
	if m == nil {
		return Commit{}, false
	}
	// m[1]=type, m[2]=scope(with parens), m[3]=!, m[4]=description
	typ := m[1]
	scope := strings.Trim(m[2], "()")
	breaking := m[3] == "!"
	desc := m[4]

	shortHash := raw.Hash
	if len(shortHash) > 7 {
		shortHash = shortHash[:7]
	}

	breakingDesc := ""
	if idx := strings.Index(raw.Body, breakingFooter); idx != -1 {
		breaking = true
		breakingDesc = strings.TrimSpace(raw.Body[idx+len(breakingFooter):])
		// trim to first newline
		if nl := strings.Index(breakingDesc, "\n"); nl != -1 {
			breakingDesc = breakingDesc[:nl]
		}
	}

	return Commit{
		Hash:                shortHash,
		Type:                typ,
		Scope:               scope,
		Description:         desc,
		Breaking:            breaking,
		BreakingDescription: breakingDesc,
		Body:                raw.Body,
	}, true
}
