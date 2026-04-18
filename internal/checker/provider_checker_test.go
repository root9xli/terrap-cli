package checker

import (
	"errors"
	"testing"

	"github.com/sirrend/terrap-cli/internal/runner"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestParseProviderOutputEmpty(t *testing.T) {
	providers := parseProviderOutput("")
	assert.Empty(t, providers)
}

func TestParseProviderOutputSingleProvider(t *testing.T) {
	input := "provider[registry.terraform.io/hashicorp/aws] 4.67.0\n"
	providers := parseProviderOutput(input)
	require.Len(t, providers, 1)
	assert.Equal(t, "registry.terraform.io/hashicorp/aws", providers[0].Name)
	assert.Equal(t, "4.67.0", providers[0].Version)
}

func TestParseProviderOutputMultipleProviders(t *testing.T) {
	input := "provider[registry.terraform.io/hashicorp/aws] 4.67.0\n" +
		"provider[registry.terraform.io/hashicorp/google] 5.1.0\n"
	providers := parseProviderOutput(input)
	assert.Len(t, providers, 2)
}

func TestGetProviderVersionsRunnerError(t *testing.T) {
	mock := runner.NewMockRunner()
	mock.SetError(errors.New("terraform not found"))
	_, err := GetProviderVersions(mock, "/tmp")
	require.Error(t, err)
	assert.Contains(t, err.Error(), "failed to list providers")
}

func TestGetProviderVersionsSuccess(t *testing.T) {
	mock := runner.NewMockRunner()
	mock.SetOutput("provider[registry.terraform.io/hashicorp/null] 3.2.1\n")
	providers, err := GetProviderVersions(mock, "/tmp")
	require.NoError(t, err)
	require.Len(t, providers, 1)
	assert.Equal(t, "registry.terraform.io/hashicorp/null", providers[0].Name)
	assert.Equal(t, "3.2.1", providers[0].Version)
	assert.True(t, mock.CalledWith("providers"))
}

func TestTrimSuffix(t *testing.T) {
	assert.Equal(t, "hello", trimSuffix("hello]", "]"))
	assert.Equal(t, "hello", trimSuffix("hello", "x"))
	assert.Equal(t, "", trimSuffix("]", "]"))
}
