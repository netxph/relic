package provider

import (
	"encoding/json"
	"fmt"
	"os/exec"
	"strings"
)

// NBGVProvider resolves version and commit range using Nerdbank.GitVersioning.
type NBGVProvider struct{}

// execNbgv runs the nbgv CLI. Var allows test override.
var execNbgv = func(args ...string) (string, error) {
	out, err := exec.Command("nbgv", args...).Output()
	return strings.TrimSpace(string(out)), err
}

// execGitLogVersionFile lists commits that changed version.json. Var allows test override.
var execGitLogVersionFile = func() (string, error) {
	out, err := exec.Command("git", "log", "--follow", "--format=%H", "--", "version.json").Output()
	return strings.TrimSpace(string(out)), err
}

// execGitShow reads a file at a specific commit ref. Var allows test override.
var execGitShow = func(ref string) (string, error) {
	out, err := exec.Command("git", "show", ref).Output()
	return strings.TrimSpace(string(out)), err
}

// execGitParent resolves the parent SHA of a commit. Returns "", nil when the
// commit has no parent (initial repo commit). Var allows test override.
var execGitParent = func(hash string) (string, error) {
	out, err := exec.Command("git", "rev-parse", "--verify", hash+"^").Output()
	if err != nil {
		return "", nil // no parent
	}
	return strings.TrimSpace(string(out)), nil
}

func (NBGVProvider) Resolve(flags CLIFlags) (ProviderResult, error) {
	if err := checkNbgvAvailable(); err != nil {
		return ProviderResult{}, err
	}

	version := flags.Version
	toRef := "HEAD"

	if version == "" {
		v, err := nbgvGetVersion()
		if err != nil {
			return ProviderResult{}, err
		}
		version = v
	} else {
		sha, err := nbgvGetCommits(version)
		if err != nil {
			return ProviderResult{}, err
		}
		toRef = sha
	}

	majorMinor := parseMajorMinor(version)
	fromRef, err := findSeriesStart(majorMinor)
	if err != nil {
		return ProviderResult{}, err
	}

	return ProviderResult{
		Version: &version,
		From:    &fromRef,
		To:      &toRef,
	}, nil
}

func checkNbgvAvailable() error {
	// ponytail: reuses execNbgv so tests that mock execNbgv also bypass this check
	if _, err := execNbgv("--version"); err != nil {
		return fmt.Errorf("nbgv not found in PATH; install Nerdbank.GitVersioning CLI or use --provider manual")
	}
	return nil
}

func nbgvGetVersion() (string, error) {
	out, err := execNbgv("get-version", "--format", "json")
	if err != nil {
		return "", fmt.Errorf("nbgv get-version failed: %w", err)
	}
	var result struct {
		SimpleVersion     string
		PrereleaseVersion string
	}
	if err := json.Unmarshal([]byte(out), &result); err != nil {
		return "", fmt.Errorf("nbgv get-version: failed to parse output: %w", err)
	}
	if result.SimpleVersion == "" {
		return "", fmt.Errorf("nbgv get-version: SimpleVersion is empty")
	}
	// PrereleaseVersion includes "-g{hash}" for non-public builds; strip it.
	return result.SimpleVersion + stripGitHash(result.PrereleaseVersion), nil
}

// stripGitHash removes the trailing "-g{hex}" suffix nbgv appends to
// PrereleaseVersion for non-public builds.
// e.g. "-beta-g40e0303fa7" → "-beta", "-beta" → "-beta", "" → ""
func stripGitHash(s string) string {
	i := strings.LastIndex(s, "-g")
	if i == -1 || len(s[i+2:]) == 0 {
		return s
	}
	for _, c := range s[i+2:] {
		if !((c >= '0' && c <= '9') || (c >= 'a' && c <= 'f')) {
			return s // not a hex hash, leave unchanged
		}
	}
	return s[:i]
}

func nbgvGetCommits(version string) (string, error) {
	out, err := execNbgv("get-commits", version)
	if err != nil {
		return "", fmt.Errorf("nbgv get-commits %s failed: %w", version, err)
	}
	if out == "" {
		return "", fmt.Errorf("nbgv get-commits %s: no commits found; version.json may use a 3-part version — not supported", version)
	}
	return strings.TrimSpace(strings.SplitN(out, "\n", 2)[0]), nil
}

// findSeriesStart walks version.json git history to find the oldest commit that
// introduced the given major.minor series, then returns its parent SHA as the
// From boundary — so the series-start commit itself is included in the range.
// Returns "" for the initial commit (no parent) or no version.json history.
// Tasks 4.1 + 4.2.
func findSeriesStart(majorMinor string) (string, error) {
	out, err := execGitLogVersionFile()
	if err != nil || out == "" {
		return "", nil // no version.json history → initial release
	}

	var seriesStart string
	for _, hash := range strings.Split(out, "\n") {
		hash = strings.TrimSpace(hash)
		if hash == "" {
			continue
		}
		content, err := execGitShow(hash + ":version.json")
		if err != nil {
			continue
		}
		if parseMajorMinor(content) == majorMinor {
			seriesStart = hash // keep updating — we want the oldest match
		} else {
			break
		}
	}

	if seriesStart == "" {
		return "", nil
	}

	// Use the parent as From so the series-start commit is inside the range.
	return execGitParent(seriesStart)
}

// parseMajorMinor extracts "major.minor" from either a version.json content
// string (JSON) or a plain version string.
// Examples: `{"version":"1.0-beta"}` → "1.0", "1.0.3-beta" → "1.0", "1.1" → "1.1"
func parseMajorMinor(input string) string {
	var v struct {
		Version string `json:"version"`
	}
	if err := json.Unmarshal([]byte(strings.TrimSpace(input)), &v); err == nil && v.Version != "" {
		input = v.Version
	}
	// strip prerelease suffix, take first two dot-separated parts
	s := strings.SplitN(input, "-", 2)[0]
	parts := strings.SplitN(s, ".", 3)
	if len(parts) < 2 {
		return s
	}
	return parts[0] + "." + parts[1]
}

var _ Provider = NBGVProvider{}
