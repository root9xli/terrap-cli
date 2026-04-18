package state

import (
	"os"
	"testing"
)

func TestSaveAndLoad(t *testing.T) {
	t.Setenv("HOME", t.TempDir())

	original := &InitData{
		WorkingDir: "/tmp/myproject",
		Backend:    "s3",
		Vars:       map[string]string{"env": "dev"},
	}

	if err := original.Save(); err != nil {
		t.Fatalf("Save() error: %v", err)
	}

	loaded, err := Load()
	if err != nil {
		t.Fatalf("Load() error: %v", err)
	}

	if loaded.WorkingDir != original.WorkingDir {
		t.Errorf("WorkingDir mismatch: got %q, want %q", loaded.WorkingDir, original.WorkingDir)
	}
	if loaded.Backend != original.Backend {
		t.Errorf("Backend mismatch: got %q, want %q", loaded.Backend, original.Backend)
	}
	if loaded.Vars["env"] != original.Vars["env"] {
		t.Errorf("Vars mismatch: got %v, want %v", loaded.Vars, original.Vars)
	}
}

func TestDelete(t *testing.T) {
	t.Setenv("HOME", t.TempDir())

	data := &InitData{WorkingDir: "/tmp/test"}
	if err := data.Save(); err != nil {
		t.Fatalf("Save() error: %v", err)
	}

	if err := Delete(); err != nil {
		t.Fatalf("Delete() error: %v", err)
	}

	path, _ := StateFilePath()
	if _, err := os.Stat(path); !os.IsNotExist(err) {
		t.Error("expected state file to be deleted")
	}
}

func TestLoadMissing(t *testing.T) {
	t.Setenv("HOME", t.TempDir())

	_, err := Load()
	if err == nil {
		t.Error("expected error when state file is missing")
	}
}
