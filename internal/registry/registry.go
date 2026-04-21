package registry

import (
	"errors"
	"sync"
	"time"
)

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

var (
	ErrNotFound  = errors.New("service not found")
	ErrDuplicate = errors.New("service name already registered")
)

type UpdateFields struct {
	Name    *string
	Project *string
	Port    *int
	Command *string
	Cwd     *string
}

type Registry struct {
	mu       sync.RWMutex
	services map[string]*Service
}

func New() *Registry {
	return &Registry{services: make(map[string]*Service)}
}

func (r *Registry) Register(svc Service) error {
	r.mu.Lock()
	defer r.mu.Unlock()
	if _, exists := r.services[svc.Name]; exists {
		return ErrDuplicate
	}
	svc.UpdatedAt = time.Now()
	svc.Status = StatusUnknown
	r.services[svc.Name] = &svc
	return nil
}

func (r *Registry) Unregister(name string) error {
	r.mu.Lock()
	defer r.mu.Unlock()
	if _, exists := r.services[name]; !exists {
		return ErrNotFound
	}
	delete(r.services, name)
	return nil
}

func (r *Registry) Get(name string) (Service, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()
	svc, exists := r.services[name]
	if !exists {
		return Service{}, ErrNotFound
	}
	return *svc, nil
}

func (r *Registry) List() []Service {
	r.mu.RLock()
	defer r.mu.RUnlock()
	out := make([]Service, 0, len(r.services))
	for _, svc := range r.services {
		out = append(out, *svc)
	}
	return out
}

func (r *Registry) Update(name string, fields UpdateFields) error {
	r.mu.Lock()
	defer r.mu.Unlock()
	svc, exists := r.services[name]
	if !exists {
		return ErrNotFound
	}
	if fields.Name != nil {
		delete(r.services, name)
		svc.Name = *fields.Name
		r.services[svc.Name] = svc
	}
	if fields.Project != nil {
		svc.Project = *fields.Project
	}
	if fields.Port != nil {
		svc.Port = *fields.Port
	}
	if fields.Command != nil {
		svc.Command = *fields.Command
	}
	if fields.Cwd != nil {
		svc.Cwd = *fields.Cwd
	}
	svc.UpdatedAt = time.Now()
	return nil
}

func (r *Registry) SetPID(name string, pid int, status Status) error {
	r.mu.Lock()
	defer r.mu.Unlock()
	svc, exists := r.services[name]
	if !exists {
		return ErrNotFound
	}
	svc.PID = pid
	svc.Status = status
	svc.UpdatedAt = time.Now()
	return nil
}

func (r *Registry) Snapshot() []Service {
	return r.List()
}

func (r *Registry) Restore(services []Service) {
	r.mu.Lock()
	defer r.mu.Unlock()
	for _, svc := range services {
		s := svc
		r.services[s.Name] = &s
	}
}
