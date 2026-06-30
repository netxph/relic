## ADDED Requirements

### Requirement: Release data model is complete
The system SHALL produce a structured release data model containing all information needed for template rendering, regardless of what the template chooses to display.

#### Scenario: Data model contains all fields
- **WHEN** release notes are generated for a range with mixed commit types
- **THEN** the data model contains: `version`, `date` (ISO 8601), `from` hash, `to` hash, `sections` (one per conventional commit type present), and `breaking_changes` (all breaking commits)

---

### Requirement: Default template renders clean release notes
The system SHALL ship a built-in default template that renders human-readable markdown release notes without commit hashes.

#### Scenario: Default output contains version header and sections
- **WHEN** `--template` flag is not provided
- **THEN** output starts with `## [<version>] - <date>` followed by non-empty sections only

#### Scenario: Scope rendered as bold prefix
- **WHEN** a commit has a scope (e.g., `fix(auth): token expiry`)
- **THEN** the rendered line is `- **auth:** token expiry`

#### Scenario: Commit without scope rendered as plain line
- **WHEN** a commit has no scope (e.g., `fix: token expiry`)
- **THEN** the rendered line is `- token expiry`

---

### Requirement: Default template shows relevant sections only
The system SHALL show only end-user-relevant sections in the default template. Internal commit types SHALL be hidden by default.

#### Scenario: Shown types appear in output
- **WHEN** commits of types `feat`, `fix`, `perf`, `revert` are present
- **THEN** those sections appear in the output under their respective labels: `Features`, `Bug Fixes`, `Performance`, `Reverts`

#### Scenario: Hidden types excluded from default output
- **WHEN** commits of types `chore`, `ci`, `build`, `test`, `style`, `refactor` are present
- **THEN** those commits do NOT appear in the default template output

#### Scenario: Breaking changes section appears first
- **WHEN** any breaking change commits are present
- **THEN** a `### ⚠ Breaking Changes` section appears before all other sections

#### Scenario: Empty sections are omitted
- **WHEN** no `fix` commits exist in the range
- **THEN** no `### Bug Fixes` section appears in the output

---

### Requirement: Custom template via flag
The system SHALL accept a `--template` flag pointing to a file path containing a Go `text/template` template. The full release data model SHALL be available to the custom template.

#### Scenario: Custom template renders hash
- **WHEN** `--template my.tmpl` is provided and the template references `.Commits.Hash`
- **THEN** the short commit hash appears in the output

#### Scenario: Missing template file produces an error
- **WHEN** `--template` points to a file that does not exist
- **THEN** the system exits with a non-zero code and a descriptive error

---

### Requirement: JSON output format
The system SHALL support a `--format json` flag that outputs the full release data model as JSON instead of rendered markdown.

#### Scenario: JSON output is valid and complete
- **WHEN** `--format json` is provided
- **THEN** output is valid JSON containing all fields of the release data model including commit hashes

---

### Requirement: Version flag with default
The system SHALL accept an optional `--version` flag. When not provided, it SHALL default to `0.0.1`.

#### Scenario: Version appears in output header
- **WHEN** `--version 1.2.3` is provided
- **THEN** output header contains `[1.2.3]`

#### Scenario: Default version when flag omitted
- **WHEN** `--version` is not provided
- **THEN** output header contains `[0.0.1]`

---

### Requirement: Output destination
The system SHALL write to stdout by default. An optional `--output` flag SHALL accept a file path and write the result to that file instead.

#### Scenario: Default output goes to stdout
- **WHEN** `--output` is not provided
- **THEN** release notes are written to stdout

#### Scenario: File output writes to specified path
- **WHEN** `--output release-notes.md` is provided
- **THEN** release notes are written to `release-notes.md` and nothing is written to stdout
