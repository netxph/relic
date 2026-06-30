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
- Provider system — CLI flags today, extensible for tag-based or CI resolution later

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

### 5. Use in CI (GitHub Actions)

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
| `--range` | *(required)* | Commit range: `<hash>` or `<from>..<to>` |
| `--version` | `0.0.1` | Version string embedded in the output |
| `--format` | `markdown` | Output format: `markdown` or `json` |
| `--output` | *(stdout)* | Write output to this file path |
| `--template` | *(built-in)* | Path to a custom Go `text/template` file |
| `--provider` | `manual` | Provider for range/version resolution (see below) |

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

Controls how `--range` and `--version` can be resolved from sources other than explicit flags.

| Provider | Description |
|----------|-------------|
| `manual` | Default. All values must be supplied via CLI flags. |

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
