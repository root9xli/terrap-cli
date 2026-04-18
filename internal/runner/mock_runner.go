package runner

// MockRunner is a test double for TerraformRunner.
type MockRunner struct {
	Calls   []MockCall
	Outputs map[string]string
	Errors  map[string]error
}

// MockCall records a single invocation.
type MockCall struct {
	Subcommand string
	Args       []string
}

// NewMockRunner creates an initialised MockRunner.
func NewMockRunner() *MockRunner {
	return &MockRunner{
		Outputs: make(map[string]string),
		Errors:  make(map[string]error),
	}
}

// Run records the call and returns the configured stub output/error.
func (m *MockRunner) Run(subcommand string, args ...string) (string, error) {
	m.Calls = append(m.Calls, MockCall{Subcommand: subcommand, Args: args})
	if err, ok := m.Errors[subcommand]; ok {
		return "", err
	}
	if out, ok := m.Outputs[subcommand]; ok {
		return out, nil
	}
	return "", nil
}

// CalledWith returns true if subcommand was called at least once.
func (m *MockRunner) CalledWith(subcommand string) bool {
	for _, c := range m.Calls {
		if c.Subcommand == subcommand {
			return true
		}
	}
	return false
}

// CallCount returns how many times subcommand was called.
func (m *MockRunner) CallCount(subcommand string) int {
	count := 0
	for _, c := range m.Calls {
		if c.Subcommand == subcommand {
			count++
		}
	}
	return count
}
