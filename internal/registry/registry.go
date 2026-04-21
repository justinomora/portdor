// internal/registry/registry.go
package registry

import "time"

type Status string

const (
	StatusRunning Status = "running"
	StatusStopped Status = "stopped"
	StatusCrashed Status = "crashed"
	StatusUnknown Status = "unknown"
)

type Service struct {
	Name      string    `json:"name"`
	Command   string    `json:"command"`
	Cwd       string    `json:"cwd"`
	Port      int       `json:"port,omitempty"`
	Project   string    `json:"project,omitempty"`
	PID       int       `json:"pid,omitempty"`
	Status    Status    `json:"status"`
	UpdatedAt time.Time `json:"updated_at"`
}
