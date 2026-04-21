package server

import (
	"fmt"
	"html/template"
	"net/http"
	"sort"

	"github.com/jmora/portdor/internal/process"
	"github.com/jmora/portdor/internal/registry"
	web "github.com/jmora/portdor/web"
)

type serviceView struct {
	Name    string
	Port    int
	Project string
	Status  string
	MemPct  float64
	CPUPct  float64
}

type groupView struct {
	Name     string
	Services []serviceView
}

type dashboardView struct {
	Meta   string
	Groups []groupView
}

var tmpl *template.Template

func init() {
	var err error
	tmpl, err = template.ParseFS(web.Assets,
		"templates/layout.html",
		"templates/dashboard.html",
		"templates/partials/group.html",
		"templates/partials/row.html",
		"templates/partials/edit.html",
	)
	if err != nil {
		panic("template parse error: " + err.Error())
	}
}

func (s *Server) buildGroups() []groupView {
	services := s.reg.List()

	grouped := map[string][]serviceView{}
	for _, svc := range services {
		info, _ := s.reg.Check(svc.Name)
		key := svc.Project
		if key == "" {
			key = "ungrouped"
		}
		sv := serviceView{
			Name:    svc.Name,
			Port:    svc.Port,
			Project: svc.Project,
			Status:  string(svc.Status),
			MemPct:  info.MemPct,
			CPUPct:  info.CPUPct,
		}
		grouped[key] = append(grouped[key], sv)
	}

	var groups []groupView
	for project, svcs := range grouped {
		sort.Slice(svcs, func(i, j int) bool { return svcs[i].Name < svcs[j].Name })
		groups = append(groups, groupView{Name: project, Services: svcs})
	}
	sort.Slice(groups, func(i, j int) bool {
		if groups[j].Name == "ungrouped" {
			return true
		}
		if groups[i].Name == "ungrouped" {
			return false
		}
		return groups[i].Name < groups[j].Name
	})
	return groups
}

func (s *Server) handleDashboard(w http.ResponseWriter, r *http.Request) {
	if r.URL.Path != "/" {
		http.NotFound(w, r)
		return
	}
	groups := s.buildGroups()
	total := 0
	running := 0
	for _, g := range groups {
		for _, sv := range g.Services {
			total++
			if sv.Status == "running" {
				running++
			}
		}
	}
	data := dashboardView{
		Meta:   fmt.Sprintf("%d services · %d running", total, running),
		Groups: groups,
	}
	tmpl.ExecuteTemplate(w, "layout.html", data)
}

func (s *Server) handleUIGroups(w http.ResponseWriter, r *http.Request) {
	groups := s.buildGroups()
	for _, g := range groups {
		tmpl.ExecuteTemplate(w, "group", g)
	}
}

func (s *Server) handleUIStop(w http.ResponseWriter, r *http.Request) {
	name := r.PathValue("name")
	if svc, err := s.reg.Get(name); err == nil && svc.PID > 0 {
		process.Stop(svc.PID)
		s.reg.SetPID(name, svc.PID, registry.StatusStopped)
		s.persist()
	}
	s.handleUIGroups(w, r)
}

func (s *Server) handleUIKill(w http.ResponseWriter, r *http.Request) {
	name := r.PathValue("name")
	if svc, err := s.reg.Get(name); err == nil && svc.PID > 0 {
		process.Kill(svc.PID)
		s.reg.SetPID(name, svc.PID, registry.StatusStopped)
		s.persist()
	}
	s.handleUIGroups(w, r)
}

func (s *Server) handleUIRestart(w http.ResponseWriter, r *http.Request) {
	name := r.PathValue("name")
	if svc, err := s.reg.Get(name); err == nil {
		if newPID, err := process.Restart(svc.PID, svc.Command, svc.Cwd); err == nil {
			s.reg.SetPID(name, newPID, registry.StatusRunning)
			s.persist()
		}
	}
	s.handleUIGroups(w, r)
}

func (s *Server) handleUIEditRow(w http.ResponseWriter, r *http.Request) {
	name := r.PathValue("name")
	svc, err := s.reg.Get(name)
	if err != nil {
		http.NotFound(w, r)
		return
	}
	sv := serviceView{Name: svc.Name, Port: svc.Port, Project: svc.Project, Status: string(svc.Status)}
	tmpl.ExecuteTemplate(w, "edit", sv)
}

func (s *Server) handleUIRow(w http.ResponseWriter, r *http.Request) {
	name := r.PathValue("name")
	svc, err := s.reg.Get(name)
	if err != nil {
		http.NotFound(w, r)
		return
	}
	sv := serviceView{Name: svc.Name, Port: svc.Port, Project: svc.Project, Status: string(svc.Status)}
	tmpl.ExecuteTemplate(w, "row", sv)
}
