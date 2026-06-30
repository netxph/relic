package provider

import (
	"strings"
	"testing"
)

func strPtr(s string) *string { return &s }

func TestResolve_CLIOverridesProvider(t *testing.T) {
	flags := CLIFlags{Version: "2.0.0"}
	result := ProviderResult{Version: strPtr("1.5.0")}
	out := Resolve(flags, result, ResolvedInput{Version: "0.0.1"})
	if out.Version != "2.0.0" {
		t.Errorf("expected 2.0.0, got %s", out.Version)
	}
}

func TestResolve_ProviderFillsMissingCLI(t *testing.T) {
	flags := CLIFlags{}
	result := ProviderResult{Version: strPtr("1.5.0")}
	out := Resolve(flags, result, ResolvedInput{Version: "0.0.1"})
	if out.Version != "1.5.0" {
		t.Errorf("expected 1.5.0, got %s", out.Version)
	}
}

func TestResolve_DefaultAppliedWhenBothEmpty(t *testing.T) {
	flags := CLIFlags{}
	result := ProviderResult{}
	out := Resolve(flags, result, ResolvedInput{Version: "0.0.1"})
	if out.Version != "0.0.1" {
		t.Errorf("expected 0.0.1, got %s", out.Version)
	}
}

func TestGet_UnknownProvider(t *testing.T) {
	_, err := Get("nbgv")
	if err == nil {
		t.Fatal("expected error for unknown provider")
	}
	if !strings.Contains(err.Error(), "available") {
		t.Errorf("error should list available providers: %v", err)
	}
}

func TestGet_ManualProvider(t *testing.T) {
	p, err := Get("manual")
	if err != nil {
		t.Fatal(err)
	}
	result, err := p.Resolve()
	if err != nil {
		t.Fatal(err)
	}
	if result.Version != nil || result.From != nil || result.To != nil {
		t.Error("manual provider should return empty result")
	}
}
