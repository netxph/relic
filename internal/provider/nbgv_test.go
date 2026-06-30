package provider

import (
	"fmt"
	"testing"
)

// mockNbgv replaces execNbgv for the duration of a test.
func mockNbgv(t *testing.T, fn func(args ...string) (string, error)) {
	t.Helper()
	old := execNbgv
	execNbgv = fn
	t.Cleanup(func() { execNbgv = old })
}

func mockGitLogVersionFile(t *testing.T, fn func() (string, error)) {
	t.Helper()
	old := execGitLogVersionFile
	execGitLogVersionFile = fn
	t.Cleanup(func() { execGitLogVersionFile = old })
}

func mockGitShow(t *testing.T, fn func(ref string) (string, error)) {
	t.Helper()
	old := execGitShow
	execGitShow = fn
	t.Cleanup(func() { execGitShow = old })
}

func mockGitParent(t *testing.T, fn func(hash string) (string, error)) {
	t.Helper()
	old := execGitParent
	execGitParent = fn
	t.Cleanup(func() { execGitParent = old })
}

// Task 6.1: no --version → resolves from nbgv get-version, To = HEAD
func TestNBGVProvider_ResolveNoVersion(t *testing.T) {
	mockNbgv(t, func(args ...string) (string, error) {
		switch args[0] {
		case "--version":
			return "3.6.133", nil
		case "get-version":
			return `{"SimpleVersion":"1.0.5","PrereleaseVersion":"-beta"}`, nil
		}
		return "", fmt.Errorf("unexpected: %v", args)
	})
	mockGitLogVersionFile(t, func() (string, error) { return "", nil })
	mockGitParent(t, func(hash string) (string, error) { return "", nil })

	result, err := NBGVProvider{}.Resolve(CLIFlags{})
	if err != nil {
		t.Fatal(err)
	}
	if result.Version == nil || *result.Version != "1.0.5-beta" {
		t.Errorf("expected Version=1.0.5-beta, got %v", result.Version)
	}
	if result.To == nil || *result.To != "HEAD" {
		t.Errorf("expected To=HEAD, got %v", result.To)
	}
}

// Task 6.2: --version supplied → provider ignores it for range; always resolves from HEAD
func TestNBGVProvider_ResolveWithVersion(t *testing.T) {
	var gotCmd string
	mockNbgv(t, func(args ...string) (string, error) {
		gotCmd = args[0]
		switch args[0] {
		case "--version":
			return "3.6.133", nil
		case "get-version":
			return `{"SimpleVersion":"1.0.5","PrereleaseVersion":"-beta"}`, nil
		}
		return "", fmt.Errorf("unexpected: %v", args)
	})
	mockGitLogVersionFile(t, func() (string, error) { return "", nil })
	mockGitParent(t, func(hash string) (string, error) { return "", nil })

	result, err := NBGVProvider{}.Resolve(CLIFlags{Version: "1.0.3-beta"})
	if err != nil {
		t.Fatal(err)
	}
	if gotCmd == "get-commits" {
		t.Error("expected nbgv get-commits NOT to be called")
	}
	if result.To == nil || *result.To != "HEAD" {
		t.Errorf("expected To=HEAD, got %v", result.To)
	}
}

// Task 6.3: findSeriesStart returns parent of oldest commit in the series
func TestFindSeriesStart(t *testing.T) {
	mockGitLogVersionFile(t, func() (string, error) {
		return "sha3\nsha2\nsha1", nil
	})
	contents := map[string]string{
		"sha3:version.json": `{"version":"1.0-beta"}`,
		"sha2:version.json": `{"version":"1.0-beta"}`,
		"sha1:version.json": `{"version":"0.9-beta"}`,
	}
	mockGitShow(t, func(ref string) (string, error) {
		if c, ok := contents[ref]; ok {
			return c, nil
		}
		return "", fmt.Errorf("not found: %s", ref)
	})
	mockGitParent(t, func(hash string) (string, error) {
		if hash == "sha2" {
			return "sha1parent", nil // parent of oldest 1.0 commit
		}
		return "", nil
	})

	from, err := findSeriesStart("1.0")
	if err != nil {
		t.Fatal(err)
	}
	// sha2 is the oldest 1.0 commit; its parent is the From boundary
	if from != "sha1parent" {
		t.Errorf("expected sha1parent, got %q", from)
	}
}

func TestFindSeriesStart_InitialRelease(t *testing.T) {
	mockGitLogVersionFile(t, func() (string, error) { return "", nil })

	from, err := findSeriesStart("1.0")
	if err != nil {
		t.Fatal(err)
	}
	if from != "" {
		t.Errorf("expected empty string for initial release, got %q", from)
	}
}

func TestFindSeriesStart_NoParent(t *testing.T) {
	mockGitLogVersionFile(t, func() (string, error) { return "sha1", nil })
	mockGitShow(t, func(ref string) (string, error) {
		return `{"version":"1.0-beta"}`, nil
	})
	mockGitParent(t, func(hash string) (string, error) {
		return "", nil // sha1 is the first commit, no parent
	})

	from, err := findSeriesStart("1.0")
	if err != nil {
		t.Fatal(err)
	}
	if from != "" {
		t.Errorf("expected empty string when no parent, got %q", from)
	}
}

func TestParseMajorMinor(t *testing.T) {
	cases := []struct{ input, want string }{
		{`{"version":"1.0-beta"}`, "1.0"},
		{`{"version":"1.0"}`, "1.0"},
		{"1.0.3-beta", "1.0"},
		{"1.0-beta", "1.0"},
		{"1.1", "1.1"},
	}
	for _, c := range cases {
		got := parseMajorMinor(c.input)
		if got != c.want {
			t.Errorf("parseMajorMinor(%q) = %q, want %q", c.input, got, c.want)
		}
	}
}

// Task 6.5: Get("nbgv") returns NBGVProvider
func TestGetNBGVProvider(t *testing.T) {
	p, err := Get("nbgv")
	if err != nil {
		t.Fatal(err)
	}
	if _, ok := p.(NBGVProvider); !ok {
		t.Errorf("expected NBGVProvider, got %T", p)
	}
}
