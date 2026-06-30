# relic

**relic** is a CLI tool that generates release notes from [Conventional Commits](https://www.conventionalcommits.org/). Point it at a git commit range, get clean markdown (or JSON) out.

```
relic --range v0.1.0..v0.2.0 --version 0.2.0
```

---

## Features

- Parses [Conventional Commits](https://www.conventionalcommits.org/) from any git repo
- Groups commits by type: Features, Bug Fixes, Performance, etc.
- Highlights breaking changes (`!` suffix or `BREAKING CHANGE:` footer)
- Outputs **Markdown** (default) or **JSON**
- Supports custom [Go text/template](https://pkg.go.dev/text/template) files
- Writes to stdout or a file via `--output`
- Provider system — `manual` for explicit ranges, `nbgv` for automatic range resolution via [Nerdbank.GitVersioning](https://github.com/dotnet/Nerdbank.GitVersioning)

---

## Installation

### Pre-built binaries

Download the latest release from [Releases](https://github.com/netxph/relic/releases) and put the binary on your `PATH`.

### Build from source

```bash
git clone https://github.com/netxph/relic.git
cd relic
make build          # produces ./relic
```

Requires Go 1.21+.

---

## Getting Started

### 1. Generate release notes for a range

```bash
relic --range abc1234..def5678 --version 1.2.0
```

Output goes to stdout as markdown:

```markdown
## [1.2.0] - 2026-06-30

### ✨ Features
- **auth:** add OAuth2 login

### 🐛 Bug Fixes
- fix null pointer in user lookup
```

### 2. Write to a file

```bash
relic --range v0.1.0..HEAD --version 1.2.0 --output CHANGELOG.md
```

### 3. JSON output

```bash
relic --range v0.1.0..HEAD --version 1.2.0 --format json
```

### 4. Custom template

```bash
relic --range v0.1.0..HEAD --version 1.2.0 --template ./my-template.tmpl
```

Templates receive a `ReleaseData` struct — see [Command Reference](#command-reference) for the full schema.

### 5. Use with Nerdbank.GitVersioning

If your repo uses [nbgv](https://github.com/dotnet/Nerdbank.GitVersioning), relic can resolve the commit range and version automatically — no tags or `--range` required.

```bash
# Current version + all commits in the active major.minor series
relic --provider nbgv

# Same, but override the displayed version label
relic --provider nbgv --version 1.0.0
```

**How it works:**

- Version is read from `nbgv get-version` (`SimpleVersion` + `PrereleaseVersion`, hash stripped).
- The commit range covers every commit since the current major.minor series began (i.e. since `version.json` was set to the current `major.minor`).
- When `nbgv prepare-release` creates a new series, the next run automatically starts fresh from that bump commit.

**Requirements:** `nbgv` must be installed and on your `PATH` (`dotnet tool install -g nbgv`).

### 6. Use in CI (GitHub Actions)

```yaml
- name: Generate release notes
  run: |
    relic --range ${{ github.event.before }}..HEAD \
          --version ${{ github.ref_name }} \
          --output RELEASE_NOTES.md
```

---

## Command Reference

```
relic [flags]
```

| Flag | Default | Description |
|------|---------|-------------|
| `--range` | *(optional)* | Commit range: `<hash>` or `<from>..<to>`. Required when not using a provider that resolves the range. |
| `--version` | *(empty)* | Version label in the output. Overrides the provider's resolved version. |
| `--format` | `markdown` | Output format: `markdown` or `json` |
| `--output` | *(stdout)* | Write output to this file path |
| `--template` | *(built-in)* | Path to a custom Go `text/template` file |
| `--provider` | `manual` | Provider for range/version resolution: `manual` or `nbgv` |

### `--range`

Specifies which commits to include.

| Form | Example | Meaning |
|------|---------|---------|
| Single hash | `--range abc1234` | From that commit to `HEAD` |
| Range | `--range v0.1.0..v0.2.0` | From tag/hash to tag/hash |
| Tag to HEAD | `--range v0.1.0..HEAD` | Everything since a tag |

### `--format`

- `markdown` — human-readable changelog block (default)
- `json` — machine-readable `ReleaseData` object, useful for downstream tooling

### `--template`

A Go [`text/template`](https://pkg.go.dev/text/template) file. Receives a `ReleaseData` value:

```
ReleaseData
├── Version         string
├── Date            string       (YYYY-MM-DD)
├── From            string       (start of range)
├── To              string       (end of range)
├── BreakingChanges []Commit
└── Sections        []Section
    ├── Type        string       (feat, fix, perf, …)
    ├── Label       string       (human label)
    └── Commits     []Commit
        ├── Hash               string
        ├── Type               string
        ├── Scope              string
        ├── Description        string
        ├── Breaking           bool
        └── BreakingDescription string
```

See [`internal/renderer/templates/default.tmpl`](internal/renderer/templates/default.tmpl) for a working example.

### `--provider`

Controls how `--range` and `--version` are resolved.

| Provider | Description |
|----------|-------------|
| `manual` | Default. `--range` and `--version` must be supplied explicitly. |
| `nbgv` | Resolves version and commit range automatically from [Nerdbank.GitVersioning](https://github.com/dotnet/Nerdbank.GitVersioning). Requires `nbgv` on `PATH`. |

#### `nbgv` provider

The `nbgv` provider reads `version.json` git history to determine the start of the current major.minor series and uses `HEAD` as the end of the range. `--version` overrides the displayed version label without affecting the resolved range.

```bash
relic --provider nbgv                    # version from nbgv, full series range
relic --provider nbgv --version 1.0.0   # label override only
```

**Series boundary:** relic walks `git log -- version.json` to find the oldest commit where `version.json` matched the current `major.minor`, then uses its parent as `From`. This means the commit that introduced the series is included in the output.

```
[version.json → "1.0-beta"]   ← series start (included)
  feat: user auth              1.0.1-beta
  feat: billing                1.0.2-beta  ← all in release notes
  fix:  crash on login         1.0.3-beta
                               ↑ HEAD

nbgv prepare-release          ← 1.0 sealed, 1.1 series begins
```

---

## Conventional Commits Support

relic recognises the following commit types:

| Type | Section label |
|------|--------------|
| `feat` | ✨ Features |
| `fix` | 🐛 Bug Fixes |
| `perf` | ⚡ Performance |
| `refactor` | ♻️ Refactors |
| `docs` | 📚 Documentation |
| `test` | 🧪 Tests |
| `chore` | 🔧 Chores |
| `build` | 📦 Build |
| `ci` | 🤖 CI |

Breaking changes (type with `!`, e.g. `feat!:`, or a `BREAKING CHANGE:` footer) are always surfaced at the top of the output under **⚠ Breaking Changes**.

Commits that do not match the Conventional Commits format are silently skipped.

---

## License

[MIT](LICENSE)
