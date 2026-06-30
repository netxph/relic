package renderer

import (
	"testing"

	"github.com/netxph/relic/internal/parser"
)

var testCommits = []parser.Commit{
	{Hash: "abc1234", Type: "feat", Scope: "ui", Description: "add dark mode"},
	{Hash: "def5678", Type: "fix", Description: "fix crash"},
	{Hash: "ghi9012", Type: "chore", Description: "update deps"},
	{Hash: "jkl3456", Type: "feat", Description: "new feature", Breaking: true, BreakingDescription: "old api removed"},
}

func TestBuildReleaseData_Grouping(t *testing.T) {
	data := BuildReleaseData("1.0.0", "abc1234", "HEAD", testCommits)
	if len(data.Sections) != 2 {
		t.Fatalf("expected 2 sections (feat, fix), got %d", len(data.Sections))
	}
	if data.Sections[0].Type != "feat" || len(data.Sections[0].Commits) != 2 {
		t.Errorf("unexpected feat section: %+v", data.Sections[0])
	}
	if data.Sections[1].Type != "fix" || len(data.Sections[1].Commits) != 1 {
		t.Errorf("unexpected fix section: %+v", data.Sections[1])
	}
}

func TestBuildReleaseData_BreakingChanges(t *testing.T) {
	data := BuildReleaseData("1.0.0", "abc1234", "HEAD", testCommits)
	if len(data.BreakingChanges) != 1 {
		t.Fatalf("expected 1 breaking change, got %d", len(data.BreakingChanges))
	}
}

func TestBuildReleaseData_HiddenTypesExcluded(t *testing.T) {
	data := BuildReleaseData("1.0.0", "abc1234", "HEAD", testCommits)
	for _, s := range data.Sections {
		if s.Type == "chore" {
			t.Error("chore should not appear in sections")
		}
	}
}

func TestBuildReleaseData_EmptySectionsOmitted(t *testing.T) {
	commits := []parser.Commit{
		{Hash: "abc1234", Type: "feat", Description: "something"},
	}
	data := BuildReleaseData("1.0.0", "abc1234", "HEAD", commits)
	for _, s := range data.Sections {
		if s.Type == "fix" {
			t.Error("fix section should not appear when no fix commits")
		}
	}
}
