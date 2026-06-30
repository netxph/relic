package provider

// ProviderResult holds optional values resolved by a provider.
// Nil fields mean "not resolved by provider; fall back to CLI or default."
type ProviderResult struct {
	Version *string
	From    *string
	To      *string
}

// Provider resolves optional version and range values from an external source.
type Provider interface {
	Resolve(flags CLIFlags) (ProviderResult, error)
}

// ResolvedInput holds the final, fully-merged values passed to the core engine.
type ResolvedInput struct {
	Version string
	From    string
	To      string
}

// CLIFlags holds the raw CLI flag values (empty string = not set by user).
type CLIFlags struct {
	Version string
	From    string
	To      string
}

// Resolve merges CLIFlags > providerResult > defaults.
func Resolve(flags CLIFlags, result ProviderResult, defaults ResolvedInput) ResolvedInput {
	out := defaults

	// Provider fills what CLI left empty.
	if flags.Version != "" {
		out.Version = flags.Version
	} else if result.Version != nil {
		out.Version = *result.Version
	}

	if flags.From != "" {
		out.From = flags.From
	} else if result.From != nil {
		out.From = *result.From
	}

	if flags.To != "" {
		out.To = flags.To
	} else if result.To != nil {
		out.To = *result.To
	}

	return out
}
