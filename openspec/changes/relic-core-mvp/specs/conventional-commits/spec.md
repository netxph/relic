## ADDED Requirements

### Requirement: Parse conventional commit subject line
The system SHALL parse each commit subject line according to the Conventional Commits specification, extracting type, optional scope, and description.

#### Scenario: Commit with type and description
- **WHEN** commit subject is `feat: add dark mode`
- **THEN** the parsed result has `type=feat`, `scope=nil`, `description="add dark mode"`

#### Scenario: Commit with type, scope, and description
- **WHEN** commit subject is `fix(auth): resolve token expiry`
- **THEN** the parsed result has `type=fix`, `scope="auth"`, `description="resolve token expiry"`

---

### Requirement: Detect breaking changes via exclamation mark
The system SHALL detect breaking changes when the commit subject contains `!` after the type or scope.

#### Scenario: Breaking change via exclamation mark
- **WHEN** commit subject is `feat!: remove legacy endpoint`
- **THEN** the parsed commit has `breaking=true`

#### Scenario: Breaking change with scope and exclamation mark
- **WHEN** commit subject is `feat(api)!: remove legacy endpoint`
- **THEN** the parsed commit has `breaking=true`

---

### Requirement: Detect breaking changes via footer
The system SHALL detect breaking changes when the commit body contains a `BREAKING CHANGE:` footer token.

#### Scenario: Breaking change via footer
- **WHEN** commit body contains `BREAKING CHANGE: drops support for v1 config`
- **THEN** the parsed commit has `breaking=true` and the footer value is captured as the breaking change description

---

### Requirement: Non-conventional commits are skipped silently
The system SHALL silently skip commits whose subject lines do not conform to the Conventional Commits format. No error or warning is emitted for non-conforming commits.

#### Scenario: Merge commit is skipped
- **WHEN** commit subject is `Merge branch 'main' into feature/dark-mode`
- **THEN** the commit is excluded from the parsed output without error

#### Scenario: Free-form commit is skipped
- **WHEN** commit subject is `WIP stuff`
- **THEN** the commit is excluded from the parsed output without error

---

### Requirement: Commit data includes short hash
The system SHALL attach the short commit hash (7 characters) to each parsed commit in the data model, regardless of whether it is rendered in the default template.

#### Scenario: Short hash present in data model
- **WHEN** a conventional commit is successfully parsed
- **THEN** the resulting commit object contains a `hash` field with exactly 7 characters
