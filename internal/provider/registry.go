package provider

import "fmt"

var registry = map[string]Provider{
	"manual": ManualProvider{},
}

// Get returns the provider by name, or an error listing available names.
func Get(name string) (Provider, error) {
	p, ok := registry[name]
	if !ok {
		return nil, fmt.Errorf("unknown provider %q; available: manual", name)
	}
	return p, nil
}
