## Context

Relic is a new greenfield CLI tool with no existing codebase. The goal is a release notes generator that decouples versioning from commit range boundary detection — unlike all existing tools (git-cliff, release-please, conventional-changelog) which are tag-centric. Target users are teams using Nerdbank.GitVersioning (NBGV) or any computed versioning strategy, primarily in .NET/Azure DevOps shops but not limited to them.

## Goals / Non-Goals

**Goals:**
- Single static binary — no runtime, SDK, or interpreter required on CI agents
- Azure Pipelines compatible via `curl` download or `go install`
- Version-source agnostic: accept version as a string parameter, not derived from git tags
- Extensible provider interface: MVP ships with manual provider, NBGV provider safely addable later without core changes
- Conventional commits support with clean default output template
- Structured JSON output for machine consumption

**Non-Goals:**
- NBGV provider implementation (MVP)
- Git-tags auto-detection provider (MVP)
- GitHub / GitLab API publishing
- CHANGELOG.md file management or appending
- Version bumping or release branch management
- Monorepo support (MVP)

## Decisions

### 1. Go over C#

**Decision**: Implement relic in Go.

**Rationale**: Go compiles to a single static binary per platform with zero runtime dependencies. The primary usage context (Azure Pipelines, GitHub Actions) benefits greatly — no .NET SDK required on agents, no Node.js, no Python. `go install` and goreleaser cover all distribution paths cleanly. C# dotnet global tools require the SDK to be present.

**Alternatives considered**: C# (dotnet global tool) — rejected because it requires .NET SDK on agents; Rust — valid but steeper learning curve for the author; Python — requires interpreter, not suited for single-binary distribution.

---

### 2. cobra for CLI

**Decision**: Use `github.com/spf13/cobra` for CLI argument parsing.

**Rationale**: The de facto standard for Go CLI tools (kubectl, gh, docker all use it). Handles flags, subcommands (future: `relic generate`, `relic providers list`), help text, and shell completion out of the box. The project will grow — cobra handles that gracefully.

**Alternatives considered**: stdlib `flag` — too limited for future subcommand expansion; `urfave/cli` — valid but cobra has broader ecosystem and contributor familiarity.

---

### 3. `text/template` for templating

**Decision**: Use Go stdlib `text/template` as the template engine.

**Rationale**: Zero external dependencies. Sufficient expressive power for release notes (loops, conditionals, variable interpolation). Every Go developer already knows it. Custom templates are `.tmpl` files using standard Go template syntax.

**Alternatives considered**: Scriban — C#-native, not applicable in Go; Handlebars.NET — same; `raymond` (Go Handlebars) — additional dependency for no meaningful gain over stdlib.

---

### 4. Shell out to `git` for MVP

**Decision**: Execute `git log` via `os/exec` rather than using `go-git`.

**Rationale**: Simplest path for MVP. `git` is universally available in environments where relic would run (CI agents always have git). Eliminates a significant dependency (go-git) from the binary. Migration to go-git later is localized to one module — the core engine never calls git directly.

**Alternatives considered**: `go-git` (pure Go git) — adds ~10MB to binary, more complexity; brings portability benefits (no git required) but premature for MVP.

---

### 5. Provider interface defined in MVP

**Decision**: Define the `Provider` interface in MVP even though only `ManualProvider` ships.

**Rationale**: Defining the interface now means NBGV and git-tags providers can be added later as new files with zero changes to core logic. The resolution pipeline (`merge(cliFlags, providerResult)`) is also written once. Deferring the interface would require touching core code when the first real provider is added.

```
Provider interface
  Resolve() → ProviderResult{ Version *string, From *string, To *string }

Resolution precedence: CLI flag > provider result > default
```

---

### 6. Commit range required in MVP

**Decision**: `--range` is a required flag in MVP. No auto-detection.

**Rationale**: Auto-detecting range (via last tag, version.json change, etc.) requires assumptions that differ per project. Explicit is always correct. Providers (git-tags, NBGV) can supply range values later — making it optional when a provider resolves it. In MVP with only ManualProvider, requiring it avoids wrong output silently.

---

### 7. goreleaser for distribution

**Decision**: Use goreleaser for multi-platform binary releases.

**Rationale**: Standard in Go ecosystem. Single config file produces Windows, Linux, macOS binaries and publishes to GitHub Releases on tag push. Also supports Homebrew tap, Scoop bucket, and Docker image generation with minimal config.

## Risks / Trade-offs

| Risk | Mitigation |
|------|-----------|
| `git` not on PATH | Detect at startup, emit clear error: `relic requires git to be installed and available on PATH` |
| Short hash ambiguity (< 7 chars) | Validate minimum 7 characters on input, emit descriptive error |
| `text/template` whitespace handling | Default template carefully managed; document whitespace control (`{{-`) for custom template authors |
| Shell injection via commit messages | Use `exec.Command` with argument slices, never string interpolation with `sh -c` |
| go-git migration later | Git operations isolated behind a single `GitClient` interface from day one — swappable without touching core |

## Open Questions

- None blocking for MVP. NBGV provider design deferred.
