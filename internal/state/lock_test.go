package state

import (
	"os"
	"testing"
)

func TestAcquireAndReleaseLock(t *testing.T) {
	dir := t.TempDir()

	if err := AcquireLock(dir, "plan"); err != nil {
		t.Fatalf("expected no error acquiring lock, got: %v", err)
	}
	defer ReleaseLock(dir)

	record, err := LoadLock(dir)
	if err != nil {
		t.Fatalf("expected no error loading lock, got: %v", err)
	}

	if record.Operation != "plan" {
		t.Errorf("expected operation 'plan', got '%s'", record.Operation)
	}
	if record.PID != os.Getpid() {
		t.Errorf("expected PID %d, got %d", os.Getpid(), record.PID)
	}
}

func TestAcquireLockWhenAlreadyLocked(t *testing.T) {
	dir := t.TempDir()

	if err := AcquireLock(dir, "plan"); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	defer ReleaseLock(dir)

	err := AcquireLock(dir, "apply")
	if err == nil {
		t.Fatal("expected error acquiring lock when already locked, got nil")
	}
}

func TestReleaseLockWhenNotLocked(t *testing.T) {
	dir := t.TempDir()

	// Releasing a non-existent lock should be a no-op, not an error.
	if err := ReleaseLock(dir); err != nil {
		t.Errorf("expected no error releasing non-existent lock, got: %v", err)
	}
}

func TestLoadLockMissing(t *testing.T) {
	dir := t.TempDir()

	_, err := LoadLock(dir)
	if err == nil {
		t.Fatal("expected error loading missing lock, got nil")
	}
}

// TestAcquireLockOperations verifies that different operation types can be
// used when acquiring a lock (e.g. "plan", "apply", "destroy").
func TestAcquireLockOperations(t *testing.T) {
	ops := []string{"plan", "apply", "destroy"}

	for _, op := range ops {
		t.Run(op, func(t *testing.T) {
			dir := t.TempDir()

			if err := AcquireLock(dir, op); err != nil {
				t.Fatalf("expected no error acquiring lock for op %q, got: %v", op, err)
			}
			defer ReleaseLock(dir)

			record, err := LoadLock(dir)
			if err != nil {
				t.Fatalf("expected no error loading lock, got: %v", err)
			}
			if record.Operation != op {
				t.Errorf("expected operation %q, got %q", op, record.Operation)
			}
		})
	}
}
