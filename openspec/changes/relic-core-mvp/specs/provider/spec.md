## ADDED Requirements

### Requirement: Provider interface contract
The system SHALL define a `Provider` interface that any provider implementation must satisfy. The interface SHALL return a `ProviderResult` containing optional version, from-ref, and to-ref values.

#### Scenario: Provider returns partial result
- **WHEN** a provider resolves only `version` and leaves `from` and `to` as nil
- **THEN** the resolution pipeline uses the provider's version and falls back to CLI flags for the missing fields

---

### Requirement: CLI flags take precedence over provider
The system SHALL merge CLI flags and provider results using the rule: explicit CLI flag overrides provider value.

#### Scenario: CLI version overrides provider version
- **WHEN** `--version 2.0.0` is supplied and the provider also resolves a version
- **THEN** `2.0.0` is used, not the provider's value

#### Scenario: Provider fills unset CLI values
- **WHEN** `--version` is not supplied and the provider resolves `version=1.5.0`
- **THEN** `1.5.0` is used as the version

---

### Requirement: Manual provider is the default
The system SHALL use the manual provider when no `--provider` flag is supplied. The manual provider is a no-op: it resolves nothing, and all values come from CLI flags or defaults.

#### Scenario: No provider flag uses manual provider
- **WHEN** `--provider` is not specified
- **THEN** the manual provider is active and CLI flags are the sole source of version and range

---

### Requirement: Provider selected via flag
The system SHALL accept a `--provider` flag to select the active provider by name.

#### Scenario: Unknown provider name produces an error
- **WHEN** `--provider unknown-name` is supplied
- **THEN** the system exits with a non-zero code and lists available provider names

---

### Requirement: Provider resolution does not affect core engine
The system SHALL ensure the core commit-parsing and rendering pipeline receives only resolved scalar values (`version string`, `from string`, `to string`). The core engine SHALL have no knowledge of providers.

#### Scenario: Core engine receives resolved values
- **WHEN** any provider is active
- **THEN** the core engine receives a fully resolved `version`, `from`, and `to` — it does not call the provider directly
