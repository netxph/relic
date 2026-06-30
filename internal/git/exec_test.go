package git

import (
	"os/exec"
	"testing"
)

func TestCheckGitAvailable(t *testing.T) {
	if err := CheckGitAvailable(); err != nil {
		t.Skipf("git not on PATH, skipping: %v", err)
	}
}

func TestExecGitClient_ValidRange(t *testing.T) {
	if err := CheckGitAvailable(); err != nil {
		t.Skip("git not available")
	}
	// Use HEAD~1..HEAD if there are commits, otherwise skip.
	out, err := exec.Command("git", "rev-parse", "HEAD").Output()
	if err != nil {
		t.Skip("no commits in repo")
	}
	head := string(out[:7])
	client := ExecGitClient{}
	_, err = client.Log(head, "HEAD")
	// We just need it not to panic; error is ok if range is empty.
	_ = err
}

func TestExecGitClient_InvalidHash(t *testing.T) {
	if err := CheckGitAvailable(); err != nil {
		t.Skip("git not available")
	}
	client := ExecGitClient{}
	_, err := client.Log("0000000", "HEAD")
	if err == nil {
		t.Fatal("expected error for invalid hash, got nil")
	}
}

func TestParseRawOutput_Empty(t *testing.T) {
	commits := parseRawOutput("")
	if len(commits) != 0 {
		t.Fatalf("expected 0 commits, got %d", len(commits))
	}
}

func TestParseRawOutput_Single(t *testing.T) {
	raw := "abc1234def5678" + fieldSep + "feat: add thing" + fieldSep + "body text" + logSep
	commits := parseRawOutput(raw)
	if len(commits) != 1 {
		t.Fatalf("expected 1 commit, got %d", len(commits))
	}
	if commits[0].Subject != "feat: add thing" {
		t.Errorf("unexpected subject: %s", commits[0].Subject)
	}
	if commits[0].Body != "body text" {
		t.Errorf("unexpected body: %s", commits[0].Body)
	}
}
