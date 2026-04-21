// internal/state/state_test.go
package state_test

import (
	"path/filepath"
	"testing"
	"time"

	"github.com/jmora/portdor/internal/registry"
	"github.com/jmora/portdor/internal/state"
)

func TestSaveAndLoad(t *testing.T) {
	dir := t.TempDir()
	path := filepath.Join(dir, "state.json")

	original := &state.State{
		Services: []registry.Service{
			{Name: "api", Port: 3000, Command: "npm run dev", Project: "myapp", Status: registry.StatusRunning},
		},
		LastUpdated: time.Now().Truncate(time.Second),
	}

	if err := state.Save(original, path); err != nil {
		t.Fatalf("Save: %v", err)
	}

	loaded, err := state.Load(path)
	if err != nil {
		t.Fatalf("Load: %v", err)
	}

	if len(loaded.Services) != 1 {
		t.Fatalf("expected 1 service, got %d", len(loaded.Services))
	}
	if loaded.Services[0].Name != "api" {
		t.Errorf("expected name=api, got %s", loaded.Services[0].Name)
	}
	if loaded.Services[0].Port != 3000 {
		t.Errorf("expected port=3000, got %d", loaded.Services[0].Port)
	}
}

func TestLoadMissingFile(t *testing.T) {
	_, err := state.Load("/nonexistent/path/state.json")
	if err == nil {
		t.Fatal("expected error for missing file, got nil")
	}
}

func TestLoadEmptyState(t *testing.T) {
	dir := t.TempDir()
	path := filepath.Join(dir, "state.json")

	empty := &state.State{}
	if err := state.Save(empty, path); err != nil {
		t.Fatalf("Save: %v", err)
	}

	loaded, err := state.Load(path)
	if err != nil {
		t.Fatalf("Load: %v", err)
	}
	if len(loaded.Services) != 0 {
		t.Errorf("expected 0 services, got %d", len(loaded.Services))
	}
}
