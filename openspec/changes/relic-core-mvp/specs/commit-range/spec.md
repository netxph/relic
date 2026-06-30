## ADDED Requirements

### Requirement: Single hash range resolves to HEAD
The system SHALL accept a single short commit hash as the `--range` value and interpret it as the start of the range with `HEAD` as the end.

#### Scenario: Single hash produces range to HEAD
- **WHEN** user provides `--range abc1234`
- **THEN** the system queries commits from `abc1234` to `HEAD` (exclusive of `abc1234`, inclusive of HEAD)

---

### Requirement: Double-hash range resolves to explicit boundary
The system SHALL accept `<from>..<to>` syntax as the `--range` value and use both hashes as explicit boundaries.

#### Scenario: Explicit range produces bounded commit list
- **WHEN** user provides `--range abc1234..def5678`
- **THEN** the system queries commits from `abc1234` to `def5678` (exclusive of `abc1234`, inclusive of `def5678`)

---

### Requirement: Range is required
The system SHALL require `--range` to be provided. It SHALL NOT auto-detect or default the range.

#### Scenario: Missing range produces an error
- **WHEN** user runs `relic` without `--range`
- **THEN** the system exits with a non-zero code and displays: `Error: --range is required`

---

### Requirement: Short hash minimum length
The system SHALL require commit hashes to be at least 7 characters long.

#### Scenario: Hash too short is rejected
- **WHEN** user provides a hash shorter than 7 characters in `--range`
- **THEN** the system exits with a non-zero code and displays a descriptive error indicating minimum hash length

---

### Requirement: Invalid hash is rejected
The system SHALL validate that hashes provided in `--range` exist in the current repository.

#### Scenario: Unknown hash produces an error
- **WHEN** user provides a hash that does not resolve to a commit in the repository
- **THEN** the system exits with a non-zero code and displays a descriptive error identifying the invalid hash

---

### Requirement: Git availability is checked at startup
The system SHALL verify that `git` is available on PATH before executing any range operation.

#### Scenario: git not found produces a clear error
- **WHEN** `git` is not available on PATH
- **THEN** the system exits with a non-zero code and displays: `relic requires git to be installed and available on PATH`
