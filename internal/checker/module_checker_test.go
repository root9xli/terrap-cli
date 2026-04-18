package checker

import (
	"testing"

	"github.com/sirrend/terrap-cli/internal/runner"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestParseModuleOutputEmpty(t *testing.T) {
	result := parseModuleOutput("")
	assert.Empty(t, result)
}

func TestParseModuleOutputSingleModule(t *testing.T) {
	input := "+ module.vpc v2.0.0"
	result := parseModuleOutput(input)
	require.Len(t, result, 1)
	assert.Equal(t, "vpc", result["vpc"].Name)
	assert.Equal(t, "v2.0.0", result["vpc"].Version)
}

func TestParseModuleOutputMultipleModules(t *testing.T) {
	input := "+ module.vpc v2.0.0\n+ module.eks v1.3.1"
	result := parseModuleOutput(input)
	assert.Len(t, result, 2)
	assert.Equal(t, "v2.0.0", result["vpc"].Version)
	assert.Equal(t, "v1.3.1", result["eks"].Version)
}

func TestParseModuleOutputIgnoresNonModuleLines(t *testing.T) {
	input := "Terraform v1.5.0\n+ module.vpc v2.0.0\non linux/amd64"
	result := parseModuleOutput(input)
	assert.Len(t, result, 1)
}

func TestGetModuleVersionsRunnerError(t *testing.T) {
	mr := runner.NewMockRunner()
	mr.SetError("version", assert.AnError)
	_, err := GetModuleVersions(mr)
	require.Error(t, err)
}

func TestGetModuleVersionsSuccess(t *testing.T) {
	mr := runner.NewMockRunner()
	mr.SetOutput("providers lock -help", "")
	mr.SetOutput("version", "+ module.vpc v3.0.0")
	versions, err := GetModuleVersions(mr)
	require.NoError(t, err)
	assert.Equal(t, "v3.0.0", versions["vpc"].Version)
}
