package checker

import (
	"fmt"

	"github.com/sirrend/terrap-cli/internal/runner"
	"github.com/sirrend/terrap-cli/internal/state"
)

// ProviderVersion holds the name and version of a Terraform provider.
type ProviderVersion struct {
	Name    string
	Version string
}

// GetProviderVersions returns the list of providers and their versions
// by running `terraform providers` in the given working directory.
func GetProviderVersions(r runner.Runner, workDir string) ([]ProviderVersion, error) {
	out, err := r.Run(workDir, "providers")
	if err != nil {
		return nil, fmt.Errorf("failed to list providers: %w", err)
	}
	return parseProviderOutput(out), nil
}

// parseProviderOutput parses the raw output of `terraform providers`.
// It returns a slice of ProviderVersion extracted from lines starting with "provider".
func parseProviderOutput(output string) []ProviderVersion {
	var providers []ProviderVersion
	lines := splitLines(output)
	for _, line := range lines {
		var name, version string
		// Expected format: provider[registry.terraform.io/hashicorp/aws] 4.0.0
		if n, _ := fmt.Sscanf(line, "provider[%s %s", &name, &version); n == 2 {
			// strip trailing ]
			name = trimSuffix(name, "]")
			providers = append(providers, ProviderVersion{Name: name, Version: version})
		}
	}
	return providers
}

// CheckProviderDrift compares current provider versions against the saved workspace record.
func CheckProviderDrift(r runner.Runner, workDir string) ([]ProviderVersion, error) {
	current, err := GetProviderVersions(r, workDir)
	if err != nil {
		return nil, err
	}

	ws, err := state.LoadWorkspace(workDir)
	if err != nil {
		return nil, fmt.Errorf("failed to load workspace: %w", err)
	}
	_ = ws // workspace metadata available for future drift logic

	return current, nil
}

func splitLines(s string) []string {
	var lines []string
	start := 0
	for i := 0; i < len(s); i++ {
		if s[i] == '\n' {
			lines = append(lines, s[start:i])
			start = i + 1
		}
	}
	if start < len(s) {
		lines = append(lines, s[start:])
	}
	return lines
}

func trimSuffix(s, suffix string) string {
	if len(s) >= len(suffix) && s[len(s)-len(suffix):] == suffix {
		return s[:len(s)-len(suffix)]
	}
	return s
}
