package server

import (
	"net/http"

	"github.com/jmora/portdor/internal/registry"
	"github.com/jmora/portdor/internal/state"
	web "github.com/jmora/portdor/web"
)

type Server struct {
	reg       *registry.Registry
	state     *state.State
	statePath string
	mux       *http.ServeMux
}

func New(reg *registry.Registry, st *state.State, statePath string) *Server {
	s := &Server{reg: reg, state: st, statePath: statePath, mux: http.NewServeMux()}
	s.routes()
	return s
}

func (s *Server) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	s.mux.ServeHTTP(w, r)
}

func (s *Server) ListenAndServe(addr string) error {
	return http.ListenAndServe(addr, s)
}

func (s *Server) routes() {
	s.mux.Handle("GET /static/", http.FileServer(http.FS(web.Assets)))
	s.mux.HandleFunc("GET /", s.handleDashboard)
	s.mux.HandleFunc("GET /ui/groups", s.handleUIGroups)
	s.mux.HandleFunc("GET /ui/edit/{name}", s.handleUIEditRow)
	s.mux.HandleFunc("GET /ui/row/{name}", s.handleUIRow)
	s.mux.HandleFunc("POST /ui/services/{name}/stop", s.handleUIStop)
	s.mux.HandleFunc("POST /ui/services/{name}/kill", s.handleUIKill)
	s.mux.HandleFunc("POST /ui/services/{name}/restart", s.handleUIRestart)
	s.mux.HandleFunc("GET /api/status", s.handleStatus)
	s.mux.HandleFunc("GET /api/services", s.handleListServices)
	s.mux.HandleFunc("POST /api/services", s.handleRegisterService)
	s.mux.HandleFunc("GET /api/services/{name}", s.handleGetService)
	s.mux.HandleFunc("PUT /api/services/{name}", s.handleUpdateService)
	s.mux.HandleFunc("DELETE /api/services/{name}", s.handleUnregisterService)
	s.mux.HandleFunc("POST /api/services/{name}/stop", s.handleStopService)
	s.mux.HandleFunc("POST /api/services/{name}/kill", s.handleKillService)
	s.mux.HandleFunc("POST /api/services/{name}/restart", s.handleRestartService)
}
