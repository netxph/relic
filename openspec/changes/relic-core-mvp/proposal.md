## Why

There is no release notes generator that works well with version-source-agnostic workflows — every existing tool (git-cliff, release-please, conventional-changelog) is either tag-centric or owns the versioning process entirely. Teams using Nerdbank.GitVersioning (NBGV) or other computed-version strategies are left without a well-fitting tool. Relic fills that gap: bring your own version, bring your own commit range, get clean release notes out.

## What Changes

This is a greenfield project. The MVP establishes the core relic CLI:

- **New**: Go CLI tool (`relic`) distributed as a single static binary
- **New**: Commit range input — accepts `<hash>` (to HEAD) or `<hash>..<hash>` (explicit range)
- **New**: Optional `--version` flag, defaults to `0.0.1` when omitted
- **New**: Conventional commit parser — reads git log in range, categorizes by type
- **New**: Template-based release notes renderer using Go's `text/template`
- **New**: Default clean template — shows Features, Bug Fixes, Performance, Reverts, Breaking Changes; hides internal types (chore, ci, build, test, style, refactor)
- **New**: `--format json` flag for structured machine-readable output
- **New**: Provider interface — extensible abstraction for version/range resolution; ships with manual provider only in MVP
- **New**: Multi-platform binary releases via goreleaser (Windows, Linux, macOS)

## Capabilities

### New Capabilities

- `commit-range`: Commit range syntax, validation, and git log querying (`<hash>` and `<hash>..<hash>` formats, short hash)
- `conventional-commits`: Parsing conventional commit messages into structured data (type, scope, description, breaking changes)
- `release-notes`: Data model, template rendering pipeline, default template, JSON output format
- `provider`: Provider interface contract and manual provider implementation (MVP); designed for future NBGV and git-tags providers

### Modified Capabilities

## Impact

- New repository — no existing code affected
- Relic's own version is managed by Nerdbank.GitVersioning (NBGV) — `version.json` is the version source of truth; relic eats its own cooking
- Runtime dependency: `git` must be available on PATH (MVP shells out via `os/exec`)
- External dependencies: `cobra` (CLI), `go-git` considered but deferred — MVP shells out to git for simplicity
- Distribution: NuGet/npm not applicable; goreleaser publishes GitHub Release binaries; `go install` supported
- Azure Pipelines: compatible via `curl` download of binary or `go install` step; no .NET SDK required
