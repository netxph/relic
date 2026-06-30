package parser

import (
	"testing"

	"github.com/netxph/relic/internal/git"
)

func raw(hash, subject, body string) git.RawCommit {
	return git.RawCommit{Hash: hash + "extra", Subject: subject, Body: body}
}

func TestParse_TypeOnly(t *testing.T) {
	commits := Parse([]git.RawCommit{raw("abc1234", "feat: add dark mode", "")})
	if len(commits) != 1 {
		t.Fatalf("expected 1, got %d", len(commits))
	}
	c := commits[0]
	if c.Type != "feat" || c.Scope != "" || c.Description != "add dark mode" {
		t.Errorf("unexpected: %+v", c)
	}
	if c.Hash != "abc1234" {
		t.Errorf("expected 7-char hash, got %q", c.Hash)
	}
}

func TestParse_TypeAndScope(t *testing.T) {
	commits := Parse([]git.RawCommit{raw("abc1234", "fix(auth): resolve token expiry", "")})
	if len(commits) != 1 {
		t.Fatal("expected 1")
	}
	c := commits[0]
	if c.Type != "fix" || c.Scope != "auth" || c.Description != "resolve token expiry" {
		t.Errorf("unexpected: %+v", c)
	}
}

func TestParse_BreakingBang(t *testing.T) {
	commits := Parse([]git.RawCommit{raw("abc1234", "feat!: remove legacy endpoint", "")})
	if !commits[0].Breaking {
		t.Error("expected breaking=true")
	}
}

func TestParse_BreakingBangWithScope(t *testing.T) {
	commits := Parse([]git.RawCommit{raw("abc1234", "feat(api)!: remove legacy endpoint", "")})
	if !commits[0].Breaking {
		t.Error("expected breaking=true")
	}
}

func TestParse_BreakingFooter(t *testing.T) {
	body := "Some details.\n\nBREAKING CHANGE: drops support for v1 config"
	commits := Parse([]git.RawCommit{raw("abc1234", "feat: something", body)})
	c := commits[0]
	if !c.Breaking {
		t.Error("expected breaking=true")
	}
	if c.BreakingDescription != "drops support for v1 config" {
		t.Errorf("unexpected breaking desc: %q", c.BreakingDescription)
	}
}

func TestParse_NonConventionalSkipped(t *testing.T) {
	raws := []git.RawCommit{
		raw("abc1234", "WIP stuff", ""),
		raw("def5678", "Merge branch 'main' into feature/dark-mode", ""),
	}
	commits := Parse(raws)
	if len(commits) != 0 {
		t.Errorf("expected 0, got %d", len(commits))
	}
}

func TestParse_ShortHash(t *testing.T) {
	commits := Parse([]git.RawCommit{raw("abc1234", "chore: update deps", "")})
	if commits[0].Hash != "abc1234" {
		t.Errorf("expected 7-char hash, got %q", commits[0].Hash)
	}
}
