package config_test

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/jmora/portdor/config"
)

const sampleTOML = `
[project]
name = "myapp"

[[services]]
name    = "api"
port    = 3000
command = "npm run dev"
cwd     = "./api"

[[services]]
name    = "worker"
command = "python manage.py worker"
cwd     = "./worker"
`

func writeTOML(t *testing.T, content string) string {
	t.Helper()
	dir := t.TempDir()
	path := filepath.Join(dir, "portdor.toml")
	os.WriteFile(path, []byte(content), 0644)
	return path
}

func TestLoadProjectConfig(t *testing.T) {
	path := writeTOML(t, sampleTOML)
	cfg, err := config.LoadProject(path)
	if err != nil {
		t.Fatalf("LoadProject: %v", err)
	}
	if cfg.Project.Name != "myapp" {
		t.Errorf("expected project name=myapp, got %s", cfg.Project.Name)
	}
	if len(cfg.Services) != 2 {
		t.Fatalf("expected 2 services, got %d", len(cfg.Services))
	}
	if cfg.Services[0].Name != "api" {
		t.Errorf("expected first service name=api, got %s", cfg.Services[0].Name)
	}
	if cfg.Services[0].Port != 3000 {
		t.Errorf("expected port=3000, got %d", cfg.Services[0].Port)
	}
}

func TestLoadProjectConfigMissingFile(t *testing.T) {
	_, err := config.LoadProject("/nonexistent/portdor.toml")
	if err == nil {
		t.Fatal("expected error for missing file")
	}
}

func TestFindProjectConfig(t *testing.T) {
	dir := t.TempDir()
	path := filepath.Join(dir, "portdor.toml")
	os.WriteFile(path, []byte(sampleTOML), 0644)

	found, err := config.FindProjectConfig(dir)
	if err != nil {
		t.Fatalf("FindProjectConfig: %v", err)
	}
	if found != path {
		t.Errorf("expected %s, got %s", path, found)
	}
}

func TestFindProjectConfigNotFound(t *testing.T) {
	dir := t.TempDir()
	_, err := config.FindProjectConfig(dir)
	if err == nil {
		t.Fatal("expected error when no portdor.toml found")
	}
}
