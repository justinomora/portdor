package server_test

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/jmora/portdor/internal/registry"
	"github.com/jmora/portdor/internal/server"
	"github.com/jmora/portdor/internal/state"
)

func newTestServer(t *testing.T) *server.Server {
	t.Helper()
	reg := registry.New()
	st := &state.State{}
	statePath := t.TempDir() + "/state.json"
	return server.New(reg, st, statePath)
}

func TestGetStatus(t *testing.T) {
	srv := newTestServer(t)
	req := httptest.NewRequest("GET", "/api/status", nil)
	w := httptest.NewRecorder()
	srv.ServeHTTP(w, req)

	if w.Code != http.StatusOK {
		t.Errorf("expected 200, got %d", w.Code)
	}
}

func TestListServicesEmpty(t *testing.T) {
	srv := newTestServer(t)
	req := httptest.NewRequest("GET", "/api/services", nil)
	w := httptest.NewRecorder()
	srv.ServeHTTP(w, req)

	if w.Code != http.StatusOK {
		t.Errorf("expected 200, got %d", w.Code)
	}

	var resp struct{ Services []registry.Service }
	json.NewDecoder(w.Body).Decode(&resp)
	if len(resp.Services) != 0 {
		t.Errorf("expected empty services, got %d", len(resp.Services))
	}
}

func TestRegisterAndList(t *testing.T) {
	srv := newTestServer(t)

	body, _ := json.Marshal(map[string]interface{}{
		"name": "api", "port": 3000, "command": "npm run dev", "project": "myapp",
	})
	req := httptest.NewRequest("POST", "/api/services", bytes.NewReader(body))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()
	srv.ServeHTTP(w, req)

	if w.Code != http.StatusCreated {
		t.Errorf("expected 201, got %d: %s", w.Code, w.Body.String())
	}

	req2 := httptest.NewRequest("GET", "/api/services", nil)
	w2 := httptest.NewRecorder()
	srv.ServeHTTP(w2, req2)

	var resp struct{ Services []registry.Service }
	json.NewDecoder(w2.Body).Decode(&resp)
	if len(resp.Services) != 1 {
		t.Fatalf("expected 1 service, got %d", len(resp.Services))
	}
	if resp.Services[0].Name != "api" {
		t.Errorf("expected name=api, got %s", resp.Services[0].Name)
	}
}

func TestUnregister(t *testing.T) {
	srv := newTestServer(t)

	body, _ := json.Marshal(map[string]interface{}{"name": "api", "port": 3000, "command": "npm run dev"})
	req := httptest.NewRequest("POST", "/api/services", bytes.NewReader(body))
	req.Header.Set("Content-Type", "application/json")
	srv.ServeHTTP(httptest.NewRecorder(), req)

	req2 := httptest.NewRequest("DELETE", "/api/services/api", nil)
	w2 := httptest.NewRecorder()
	srv.ServeHTTP(w2, req2)

	if w2.Code != http.StatusNoContent {
		t.Errorf("expected 204, got %d", w2.Code)
	}
}

func TestGetServiceNotFound(t *testing.T) {
	srv := newTestServer(t)
	req := httptest.NewRequest("GET", "/api/services/nonexistent", nil)
	w := httptest.NewRecorder()
	srv.ServeHTTP(w, req)

	if w.Code != http.StatusNotFound {
		t.Errorf("expected 404, got %d", w.Code)
	}
}
