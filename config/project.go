package config

import (
	"fmt"
	"os"
	"path/filepath"

	"github.com/pelletier/go-toml/v2"
)

type ProjectConfig struct {
	Project  ProjectMeta     `toml:"project"`
	Services []ServiceConfig `toml:"services"`
}

type ProjectMeta struct {
	Name string `toml:"name"`
}

type ServiceConfig struct {
	Name    string `toml:"name"`
	Port    int    `toml:"port"`
	Command string `toml:"command"`
	Cwd     string `toml:"cwd"`
}

func LoadProject(path string) (*ProjectConfig, error) {
	data, err := os.ReadFile(path)
	if err != nil {
		return nil, fmt.Errorf("read %s: %w", path, err)
	}
	var cfg ProjectConfig
	if err := toml.Unmarshal(data, &cfg); err != nil {
		return nil, fmt.Errorf("parse %s: %w", path, err)
	}
	dir := filepath.Dir(path)
	for i := range cfg.Services {
		if cfg.Services[i].Cwd == "" {
			cfg.Services[i].Cwd = dir
		} else if !filepath.IsAbs(cfg.Services[i].Cwd) {
			cfg.Services[i].Cwd = filepath.Join(dir, cfg.Services[i].Cwd)
		}
	}
	return &cfg, nil
}

func FindProjectConfig(dir string) (string, error) {
	candidate := filepath.Join(dir, "portdor.toml")
	if _, err := os.Stat(candidate); err == nil {
		return candidate, nil
	}
	return "", fmt.Errorf("portdor.toml not found in %s", dir)
}
