package renderer

import (
	"time"

	"github.com/netxph/relic/internal/parser"
)

// sectionOrder defines display order and labels for visible types.
var sectionOrder = []struct {
	typ   string
	label string
}{
	{"feat", "Features"},
	{"fix", "Bug Fixes"},
	{"perf", "Performance"},
	{"revert", "Reverts"},
}

// BuildReleaseData groups commits into sections and extracts breaking changes.
func BuildReleaseData(version, from, to string, commits []parser.Commit) ReleaseData {
	data := ReleaseData{
		Version: version,
		Date:    time.Now().Format("2006-01-02"),
		From:    from,
		To:      to,
	}

	byType := make(map[string][]parser.Commit)
	for _, c := range commits {
		byType[c.Type] = append(byType[c.Type], c)
		if c.Breaking {
			data.BreakingChanges = append(data.BreakingChanges, c)
		}
	}

	for _, s := range sectionOrder {
		if cs, ok := byType[s.typ]; ok && len(cs) > 0 {
			data.Sections = append(data.Sections, Section{
				Type:    s.typ,
				Label:   s.label,
				Commits: cs,
			})
		}
	}

	return data
}
