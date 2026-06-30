## 1. Project Setup

- [x] 1.1 Initialize Go module (`go mod init github.com/<user>/relic`)
- [x] 1.2 Add `cobra` dependency (`go get github.com/spf13/cobra`)
- [x] 1.3 Create project directory structure (`cmd/`, `internal/git/`, `internal/parser/`, `internal/renderer/`, `internal/provider/`, `templates/`)
- [x] 1.4 Create root `main.go` wiring cobra root command
- [x] 1.5 Add `.goreleaser.yaml` for multi-platform binary builds (Windows, Linux, macOS amd64/arm64)
- [x] 1.6 Add `Makefile` with `build`, `test`, `release` targets

## 2. Git Client

- [x] 2.1 Define `GitClient` interface in `internal/git/client.go` (`Log(from, to string) ([]RawCommit, error)`)
- [x] 2.2 Implement `ExecGitClient` using `os/exec` to run `git log --format=...` with argument slices (no shell interpolation)
- [x] 2.3 Implement git availability check at startup (`exec.LookPath("git")`)
- [x] 2.4 Write `ExecGitClient` tests covering: valid range, single hash to HEAD, invalid hash error, git not found

## 3. Commit Range

- [x] 3.1 Implement range parser in `internal/git/range.go` — accepts `<hash>` and `<hash>..<hash>` formats
- [x] 3.2 Validate minimum hash length (7 characters) with descriptive error
- [x] 3.3 Resolve single hash to `<hash>..HEAD` internally
- [x] 3.4 Write range parser tests covering: single hash, explicit range, too-short hash, empty input

## 4. Conventional Commits Parser

- [x] 4.1 Define `Commit` struct in `internal/parser/commit.go` (`Hash`, `Type`, `Scope`, `Description`, `Breaking`, `BreakingDescription`, `Body`)
- [x] 4.2 Implement subject line parser (regex: `^(\w+)(\([\w-]+\))?(!)?: (.+)$`)
- [x] 4.3 Implement breaking change detection via `!` suffix
- [x] 4.4 Implement breaking change detection via `BREAKING CHANGE:` footer in commit body
- [x] 4.5 Silently skip non-conforming commit subjects
- [x] 4.6 Attach 7-character short hash to each parsed commit
- [x] 4.7 Write parser tests covering: type only, type+scope, breaking `!`, breaking footer, non-conventional skip, merge commit skip

## 5. Provider Interface

- [x] 5.1 Define `Provider` interface and `ProviderResult` struct in `internal/provider/provider.go`
- [x] 5.2 Implement `ManualProvider` (no-op, returns empty `ProviderResult`)
- [x] 5.3 Implement `Resolve(cliFlags, providerResult, defaults) → ResolvedInput` merge function (CLI > provider > default)
- [x] 5.4 Implement provider registry and `--provider` flag lookup with descriptive error for unknown names
- [x] 5.5 Write resolution tests covering: CLI overrides provider, provider fills missing CLI, default applied when both empty

## 6. Release Notes Data Model & Renderer

- [x] 6.1 Define `ReleaseData` struct in `internal/renderer/model.go` (`Version`, `Date`, `From`, `To`, `Sections []Section`, `BreakingChanges []Commit`)
- [x] 6.2 Define `Section` struct (`Type`, `Label`, `Commits []Commit`)
- [x] 6.3 Implement `BuildReleaseData` — groups parsed commits into sections, separates breaking changes, sets date to today (ISO 8601)
- [x] 6.4 Define section ordering and label map (`feat→Features`, `fix→Bug Fixes`, `perf→Performance`, `revert→Reverts`)
- [x] 6.5 Write build tests covering: grouping by type, breaking changes extracted, empty sections omitted, hidden types excluded

## 7. Default Template

- [x] 7.1 Create `templates/default.tmpl` with clean markdown output
- [x] 7.2 Template renders: `## [version] - date` header, `### ⚠ Breaking Changes` section (if any), then visible sections
- [x] 7.3 Template renders scope as `**scope:** description` when scope present, plain `description` when absent
- [x] 7.4 Template omits empty sections
- [x] 7.5 Embed default template in binary using `//go:embed templates/default.tmpl`
- [x] 7.6 Write renderer integration test: full pipeline from raw commits to rendered markdown output

## 8. CLI Wiring

- [x] 8.1 Implement `generate` command (or root command) in `cmd/root.go` with flags: `--range` (required), `--version` (default `0.0.1`), `--provider` (default `manual`), `--format` (`markdown`/`json`), `--template` (file path), `--output` (file path)
- [x] 8.2 Wire git availability check before any execution
- [x] 8.3 Wire range parsing → git log → commit parsing → provider resolution → render pipeline
- [x] 8.4 Implement `--format json` output path using `encoding/json`
- [x] 8.5 Implement `--template` file loading and validation (error if file not found)
- [x] 8.6 Implement `--output` file writing (default stdout)
- [x] 8.7 Validate `--range` is provided; emit `Error: --range is required` if missing

## 9. End-to-End Validation

- [ ] 9.1 Run relic against this repository's own commit history as a smoke test
- [ ] 9.2 Verify `--format json` output is valid JSON with all expected fields
- [ ] 9.3 Verify custom `--template` flag with a test template that surfaces commit hashes
- [ ] 9.4 Verify unknown `--provider` flag produces a clear error listing available providers
- [ ] 9.5 Verify missing `--range` produces `Error: --range is required`
- [ ] 9.6 Verify binary builds successfully for Windows, Linux, macOS via `goreleaser build --snapshot`
