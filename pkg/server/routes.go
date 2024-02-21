package server

import (
	"net/http"
	"text/template"
)

// set up routes
func (s *server) Routes() {
	sData := siteData{
		Title:   "connors-go-htmx-starter",
		DevMode: true,
	}
	s.Router.Handle("/static/*", http.FileServer(http.FS(s.FileSystem)))

	s.Router.HandleFunc("/", s.handleIndex(s.getIndexTemplate(), sData))
	s.Router.HandleFunc("/projects/", s.handleProjects(s.getProjectsTemplate(), sData))
	s.Router.HandleFunc("/api/projects/", s.handleApiProjects(s.getApiProjectsPartial()))
}

func (s *server) handleIndex(tmpl *template.Template, sData siteData) http.HandlerFunc {
	type ViewData struct {
		SiteData siteData

		PageTitle string
		Username  string
		Option    string
		Offset    string
	}

	data := ViewData{
		SiteData: sData,

		PageTitle: "Index",
		Username:  "connorkuljis",
		Option:    "dev",
		Offset:    "-2",
	}

	return func(w http.ResponseWriter, r *http.Request) {
		tmpl.ExecuteTemplate(w, "root", data)
	}
}

func (s *server) handleProjects(tmpl *template.Template, sData siteData) http.HandlerFunc {
	type ViewData struct {
		SiteData siteData

		PageTitle string
		Username  string
	}

	data := ViewData{
		SiteData: sData,

		PageTitle: "Projects",
		Username:  "connorkuljis",
	}

	return func(w http.ResponseWriter, r *http.Request) {
		tmpl.ExecuteTemplate(w, "root", data)
	}
}
