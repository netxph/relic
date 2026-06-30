package provider

// ManualProvider is the default no-op provider: everything comes from CLI flags.
type ManualProvider struct{}

func (ManualProvider) Resolve(_ CLIFlags) (ProviderResult, error) {
	return ProviderResult{}, nil
}

var _ Provider = ManualProvider{}
