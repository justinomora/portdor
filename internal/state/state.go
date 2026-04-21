// internal/state/state.go
package state

import (
	"encoding/json"
	"os"
	"path/filepath"
	"time"

	"github.com/jmora/portdor/internal/registry"
)

type State struct {
	Services    []registry.Service `json:"services"`
	LastUpdated time.Time          `json:"last_updated"`
}

func DefaultPath() string {
	home, _ := os.UserHomeDir()
	return filepath.Join(home, ".portdor", "state.json")
}

func Load(path string) (*State, error) {
	data, err := os.ReadFile(path)
	if err != nil {
		return nil, err
	}
	var s State
	if err := json.Unmarshal(data, &s); err != nil {
		return nil, err
	}
	return &s, nil
}

func Save(s *State, path string) error {
	s.LastUpdated = time.Now()
	data, err := json.MarshalIndent(s, "", "  ")
	if err != nil {
		return err
	}
	if err := os.MkdirAll(filepath.Dir(path), 0755); err != nil {
		return err
	}
	return os.WriteFile(path, data, 0644)
}
