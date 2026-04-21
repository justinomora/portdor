package registry_test

import (
	"testing"

	"github.com/jmora/portdor/internal/registry"
)

func TestRegister(t *testing.T) {
	r := registry.New()
	svc := registry.Service{Name: "api", Port: 3000, Command: "npm run dev", Project: "myapp"}

	if err := r.Register(svc); err != nil {
		t.Fatalf("Register: %v", err)
	}

	services := r.List()
	if len(services) != 1 {
		t.Fatalf("expected 1 service, got %d", len(services))
	}
	if services[0].Name != "api" {
		t.Errorf("expected name=api, got %s", services[0].Name)
	}
}

func TestRegisterDuplicate(t *testing.T) {
	r := registry.New()
	svc := registry.Service{Name: "api", Port: 3000, Command: "npm run dev"}
	r.Register(svc)

	err := r.Register(svc)
	if err == nil {
		t.Fatal("expected error for duplicate name, got nil")
	}
}

func TestUnregister(t *testing.T) {
	r := registry.New()
	r.Register(registry.Service{Name: "api", Port: 3000, Command: "npm run dev"})
	r.Register(registry.Service{Name: "frontend", Port: 3001, Command: "npm start"})

	if err := r.Unregister("api"); err != nil {
		t.Fatalf("Unregister: %v", err)
	}

	services := r.List()
	if len(services) != 1 {
		t.Fatalf("expected 1 service, got %d", len(services))
	}
	if services[0].Name != "frontend" {
		t.Errorf("expected frontend to remain, got %s", services[0].Name)
	}
}

func TestUnregisterNotFound(t *testing.T) {
	r := registry.New()
	err := r.Unregister("nonexistent")
	if err == nil {
		t.Fatal("expected error for unknown service, got nil")
	}
}

func TestUpdate(t *testing.T) {
	r := registry.New()
	r.Register(registry.Service{Name: "api", Port: 3000, Command: "npm run dev", Project: "old"})

	if err := r.Update("api", registry.UpdateFields{Project: strPtr("newproject")}); err != nil {
		t.Fatalf("Update: %v", err)
	}

	svc, err := r.Get("api")
	if err != nil {
		t.Fatalf("Get: %v", err)
	}
	if svc.Project != "newproject" {
		t.Errorf("expected project=newproject, got %s", svc.Project)
	}
	if svc.Port != 3000 {
		t.Errorf("port should not have changed, got %d", svc.Port)
	}
}

func TestGet(t *testing.T) {
	r := registry.New()
	r.Register(registry.Service{Name: "api", Port: 3000, Command: "npm run dev"})

	svc, err := r.Get("api")
	if err != nil {
		t.Fatalf("Get: %v", err)
	}
	if svc.Name != "api" {
		t.Errorf("expected name=api, got %s", svc.Name)
	}
}

func TestGetNotFound(t *testing.T) {
	r := registry.New()
	_, err := r.Get("nonexistent")
	if err == nil {
		t.Fatal("expected error for unknown service, got nil")
	}
}

func strPtr(s string) *string { return &s }
