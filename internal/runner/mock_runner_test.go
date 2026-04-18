package runner

import (
	"errors"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestMockRunnerRecordsCalls(t *testing.T) {
	m := NewMockRunner()
	m.Outputs["plan"] = "No changes."

	out, err := m.Run("plan", "-detailed-exitcode")
	require.NoError(t, err)
	assert.Equal(t, "No changes.", out)
	assert.True(t, m.CalledWith("plan"))
	assert.Equal(t, 1, m.CallCount("plan"))
}

func TestMockRunnerReturnsError(t *testing.T) {
	m := NewMockRunner()
	m.Errors["destroy"] = errors.New("destroy failed")

	_, err := m.Run("destroy")
	assert.EqualError(t, err, "destroy failed")
}

func TestMockRunnerCallCount(t *testing.T) {
	m := NewMockRunner()
	m.Run("init")
	m.Run("init")
	m.Run("plan")

	assert.Equal(t, 2, m.CallCount("init"))
	assert.Equal(t, 1, m.CallCount("plan"))
	assert.Equal(t, 0, m.CallCount("destroy"))
}

func TestMockRunnerCalledWith(t *testing.T) {
	m := NewMockRunner()
	assert.False(t, m.CalledWith("version"))
	m.Run("version")
	assert.True(t, m.CalledWith("version"))
}
