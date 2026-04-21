package server

import (
	"encoding/json"
	"errors"
	"net/http"
	"strings"

	"github.com/jmora/portdor/internal/process"
	"github.com/jmora/portdor/internal/registry"
	"github.com/jmora/portdor/internal/state"
)

func writeJSON(w http.ResponseWriter, code int, v any) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(code)
	json.NewEncoder(w).Encode(v)
}

func writeError(w http.ResponseWriter, code int, msg string) {
	writeJSON(w, code, map[string]string{"error": msg})
}

func (s *Server) persist() {
	s.state.Services = s.reg.Snapshot()
	state.Save(s.state, s.statePath)
}

func (s *Server) handleStatus(w http.ResponseWriter, r *http.Request) {
	writeJSON(w, http.StatusOK, map[string]string{"status": "ok"})
}

func (s *Server) handleListServices(w http.ResponseWriter, r *http.Request) {
	s.reg.CheckAll()
	services := s.reg.List()
	writeJSON(w, http.StatusOK, map[string]any{"services": services})
}

func (s *Server) handleGetService(w http.ResponseWriter, r *http.Request) {
	name := r.PathValue("name")
	svc, err := s.reg.Get(name)
	if errors.Is(err, registry.ErrNotFound) {
		writeError(w, http.StatusNotFound, "service not found")
		return
	}
	info, _ := s.reg.Check(name)
	writeJSON(w, http.StatusOK, map[string]any{"service": svc, "health": info})
}

func (s *Server) handleRegisterService(w http.ResponseWriter, r *http.Request) {
	var req struct {
		Name    string `json:"name"`
		Command string `json:"command"`
		Cwd     string `json:"cwd"`
		Port    int    `json:"port"`
		Project string `json:"project"`
	}
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		writeError(w, http.StatusBadRequest, "invalid JSON")
		return
	}
	if req.Name == "" || req.Command == "" {
		writeError(w, http.StatusBadRequest, "name and command are required")
		return
	}

	svc := registry.Service{
		Name:    req.Name,
		Command: req.Command,
		Cwd:     req.Cwd,
		Port:    req.Port,
		Project: req.Project,
	}

	if req.Port > 0 {
		if pid := registry.PIDForPort(req.Port); pid > 0 {
			svc.PID = pid
			svc.Status = registry.StatusRunning
		}
	}

	if err := s.reg.Register(svc); err != nil {
		if errors.Is(err, registry.ErrDuplicate) {
			writeError(w, http.StatusConflict, "service name already registered")
			return
		}
		writeError(w, http.StatusInternalServerError, err.Error())
		return
	}

	s.persist()
	writeJSON(w, http.StatusCreated, map[string]any{"service": svc})
}

func (s *Server) handleUpdateService(w http.ResponseWriter, r *http.Request) {
	name := r.PathValue("name")

	var fields registry.UpdateFields

	contentType := r.Header.Get("Content-Type")
	if strings.Contains(contentType, "application/x-www-form-urlencoded") {
		r.ParseForm()
		if v := r.FormValue("name"); v != "" && v != name {
			fields.Name = &v
		}
		if v := r.FormValue("project"); v != "" {
			fields.Project = &v
		}
	} else {
		var req struct {
			Name    *string `json:"name"`
			Project *string `json:"project"`
			Port    *int    `json:"port"`
			Command *string `json:"command"`
			Cwd     *string `json:"cwd"`
		}
		if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
			writeError(w, http.StatusBadRequest, "invalid JSON")
			return
		}
		fields = registry.UpdateFields{
			Name: req.Name, Project: req.Project,
			Port: req.Port, Command: req.Command, Cwd: req.Cwd,
		}
	}

	if err := s.reg.Update(name, fields); errors.Is(err, registry.ErrNotFound) {
		writeError(w, http.StatusNotFound, "service not found")
		return
	}
	s.persist()
	w.WriteHeader(http.StatusNoContent)
}

func (s *Server) handleUnregisterService(w http.ResponseWriter, r *http.Request) {
	name := r.PathValue("name")
	if err := s.reg.Unregister(name); errors.Is(err, registry.ErrNotFound) {
		writeError(w, http.StatusNotFound, "service not found")
		return
	}
	s.persist()
	w.WriteHeader(http.StatusNoContent)
}

func (s *Server) handleStopService(w http.ResponseWriter, r *http.Request) {
	name := r.PathValue("name")
	svc, err := s.reg.Get(name)
	if errors.Is(err, registry.ErrNotFound) {
		writeError(w, http.StatusNotFound, "service not found")
		return
	}
	if err := process.Stop(svc.PID); err != nil {
		writeError(w, http.StatusInternalServerError, err.Error())
		return
	}
	s.reg.SetPID(name, svc.PID, registry.StatusStopped)
	s.persist()
	w.WriteHeader(http.StatusNoContent)
}

func (s *Server) handleKillService(w http.ResponseWriter, r *http.Request) {
	name := r.PathValue("name")
	svc, err := s.reg.Get(name)
	if errors.Is(err, registry.ErrNotFound) {
		writeError(w, http.StatusNotFound, "service not found")
		return
	}
	if err := process.Kill(svc.PID); err != nil {
		writeError(w, http.StatusInternalServerError, err.Error())
		return
	}
	s.reg.SetPID(name, svc.PID, registry.StatusStopped)
	s.persist()
	w.WriteHeader(http.StatusNoContent)
}

func (s *Server) handleRestartService(w http.ResponseWriter, r *http.Request) {
	name := r.PathValue("name")
	svc, err := s.reg.Get(name)
	if errors.Is(err, registry.ErrNotFound) {
		writeError(w, http.StatusNotFound, "service not found")
		return
	}
	newPID, err := process.Restart(svc.PID, svc.Command, svc.Cwd)
	if err != nil {
		writeError(w, http.StatusInternalServerError, err.Error())
		return
	}
	s.reg.SetPID(name, newPID, registry.StatusRunning)
	s.persist()
	writeJSON(w, http.StatusOK, map[string]any{"pid": newPID})
}
