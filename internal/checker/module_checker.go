package checker

import (
	"fmt"
	"strings"

	"github.com/sirrend/terrap-cli/internal/runner"
	"github.com/sirrend/terrap-cli/internal/state"
)

// ModuleVersion holds the name and version of a Terraform module.
type ModuleVersion struct {
	Name    string
	Version string
}

// GetModuleVersions runs `terraform version -json` style parsing to extract
// module sources from the lock file output produced by the runner.
func GetModuleVersions(r runner.Runner) (map[string]ModuleVersion, error) {
	out, err := r.Run("providers", "lock", "-help")
	_ = out
	if err != nil {
		// fallback: use version output
	}

	raw, err := r.Run("version")
	if err != nil {
		return nil, fmt.Errorf("get module versions: %w", err)
	}
	return parseModuleOutput(raw), nil
}

func parseModuleOutput(output string) map[string]ModuleVersion {
	result := make(map[string]ModuleVersion)
	for _, line := range splitLines(output) {
		line = strings.TrimSpace(line)
		if !strings.HasPrefix(line, "+ module.") {
			continue
		}
		parts := strings.SplitN(line, " ", 3)
		if len(parts) < 3 {
			continue
		}
		name := strings.TrimPrefix(parts[1], "module.")
		version := trimSuffix(parts[2])
		result[name] = ModuleVersion{Name: name, Version: version}
	}
	return result
}

// CheckModuleDrift compares current module versions against the saved record.
func CheckModuleDrift(r runner.Runner, stateDir string) ([]string, error) {
	current, err := GetModuleVersions(r)
	if err != nil {
		return nil, err
	}

	record, err := state.LoadVersionRecord(stateDir)
	if err != nil {
		return nil, err
	}

	var drifted []string
	for name, mv := range current {
		key := "module." + name
		if prev, ok := record.Modules[key]; ok && prev != mv.Version {
			drifted = append(drifted, fmt.Sprintf("%s: %s -> %s", name, prev, mv.Version))
		}
	}
	return drifted, nil
}
