package cmd

import (
	"encoding/json"
	"fmt"
	"os"
	"strings"

	"github.com/spf13/cobra"

	igit "github.com/netxph/relic/internal/git"
	"github.com/netxph/relic/internal/parser"
	"github.com/netxph/relic/internal/provider"
	"github.com/netxph/relic/internal/renderer"
)

var (
	flagRange    string
	flagVersion  string
	flagProvider string
	flagFormat   string
	flagTemplate string
	flagOutput   string
)

var rootCmd = &cobra.Command{
	Use:   "relic",
	Short: "Relic — release notes generator for any versioning strategy",
	RunE:  run,
}

func init() {
	rootCmd.Flags().StringVar(&flagRange, "range", "", "Commit range: <hash> or <hash>..<hash> (required)")
	rootCmd.Flags().StringVar(&flagVersion, "version", "", "Version string to embed in release notes")
	rootCmd.Flags().StringVar(&flagProvider, "provider", "manual", "Provider to use for version/range resolution")
	rootCmd.Flags().StringVar(&flagFormat, "format", "markdown", "Output format: markdown or json")
	rootCmd.Flags().StringVar(&flagTemplate, "template", "", "Path to a custom Go text/template file")
	rootCmd.Flags().StringVar(&flagOutput, "output", "", "Write output to file instead of stdout")
}

func run(cmd *cobra.Command, args []string) error {
	// 8.2 git availability check
	if err := igit.CheckGitAvailable(); err != nil {
		return err
	}

	// 8.1 resolve provider
	p, err := provider.Get(flagProvider)
	if err != nil {
		return err
	}
	provResult, err := p.Resolve(provider.CLIFlags{Version: flagVersion})
	if err != nil {
		return err
	}

	// 8.7 --range required only when provider did not fill From/To
	if flagRange == "" && (provResult.From == nil || provResult.To == nil) {
		return fmt.Errorf("--range is required")
	}

	// 8.3 parse range if provided
	var parsedFrom, parsedTo string
	if flagRange != "" {
		parsed, err := igit.ParseRange(flagRange)
		if err != nil {
			return err
		}
		parsedFrom, parsedTo = parsed.From, parsed.To
	}

	// merge: CLI > provider > defaults
	resolved := provider.Resolve(
		provider.CLIFlags{Version: flagVersion, From: parsedFrom, To: parsedTo},
		provResult,
		provider.ResolvedInput{Version: "0.0.1"},
	)

	// git log
	client := igit.ExecGitClient{}
	raws, err := client.Log(resolved.From, resolved.To)
	if err != nil {
		return err
	}

	// parse commits
	commits := parser.Parse(raws)

	// build data model
	data := renderer.BuildReleaseData(resolved.Version, resolved.From, resolved.To, commits)

	// choose output destination
	out := cmd.OutOrStdout()
	if flagOutput != "" {
		f, err := os.Create(flagOutput)
		if err != nil {
			return fmt.Errorf("cannot open output file: %w", err)
		}
		defer f.Close()
		out = f
	}

	// 8.4 json format
	if strings.EqualFold(flagFormat, "json") {
		return json.NewEncoder(out).Encode(data)
	}

	// 8.5 custom template
	tmplStr := ""
	if flagTemplate != "" {
		tmplStr, err = renderer.LoadTemplate(flagTemplate)
		if err != nil {
			return fmt.Errorf("cannot load template %q: %w", flagTemplate, err)
		}
	}

	return renderer.Render(data, tmplStr, out)
}

// Execute runs the root command.
func Execute() {
	if err := rootCmd.Execute(); err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
}
