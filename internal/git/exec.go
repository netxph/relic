package git

import (
	"fmt"
	"os/exec"
	"strings"
)

const logSep = "---GIT-SEP---"
const fieldSep = "---FIELD---"

// ExecGitClient runs git via os/exec with no shell interpolation.
type ExecGitClient struct{}

// Log returns commits in [from..to]. If from is empty, logs all commits up to to.
func (c ExecGitClient) Log(from, to string) ([]RawCommit, error) {
	format := "%H" + fieldSep + "%s" + fieldSep + "%b" + logSep
	var out []byte
	var err error
	if from == "" {
		out, err = exec.Command("git", "log", "--format="+format, to).Output()
	} else {
		rangeArg := from + ".." + to
		out, err = exec.Command("git", "log", "--format="+format, rangeArg).Output()
	}
	if err != nil {
		return nil, fmt.Errorf("git log failed: %w", err)
	}
	return parseRawOutput(string(out)), nil
}

func parseRawOutput(raw string) []RawCommit {
	var commits []RawCommit
	for _, entry := range strings.Split(raw, logSep) {
		entry = strings.TrimSpace(entry)
		if entry == "" {
			continue
		}
		parts := strings.SplitN(entry, fieldSep, 3)
		if len(parts) < 2 {
			continue
		}
		hash := strings.TrimSpace(parts[0])
		subject := strings.TrimSpace(parts[1])
		body := ""
		if len(parts) == 3 {
			body = strings.TrimSpace(parts[2])
		}
		commits = append(commits, RawCommit{Hash: hash, Subject: subject, Body: body})
	}
	return commits
}

// CheckGitAvailable returns an error if git is not on PATH.
func CheckGitAvailable() error {
	if _, err := exec.LookPath("git"); err != nil {
		return fmt.Errorf("relic requires git to be installed and available on PATH")
	}
	return nil
}

var _ GitClient = ExecGitClient{}
