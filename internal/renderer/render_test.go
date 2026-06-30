package renderer

import (
	"strings"
	"testing"

	"github.com/netxph/relic/internal/parser"
)

func TestRender_Integration(t *testing.T) {
	commits := []parser.Commit{
		{Hash: "abc1234", Type: "feat", Scope: "ui", Description: "add dark mode"},
		{Hash: "def5678", Type: "fix", Description: "fix crash"},
		{Hash: "ghi9012", Type: "feat", Description: "new feature", Breaking: true, BreakingDescription: "old api removed"},
	}
	data := BuildReleaseData("1.2.3", "abc1234", "HEAD", commits)

	var sb strings.Builder
	if err := Render(data, "", &sb); err != nil {
		t.Fatal(err)
	}
	out := sb.String()

	if !strings.Contains(out, "## [1.2.3]") {
		t.Errorf("missing version header; got:\n%s", out)
	}
	if !strings.Contains(out, "⚠ Breaking Changes") {
		t.Errorf("missing breaking changes section; got:\n%s", out)
	}
	if !strings.Contains(out, "**ui:** add dark mode") {
		t.Errorf("missing scoped commit; got:\n%s", out)
	}
	if !strings.Contains(out, "- fix crash") {
		t.Errorf("missing plain commit; got:\n%s", out)
	}
}
