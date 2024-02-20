package main

import (
	"net/http"
	"text/template"
)

type SiteData struct {
	Title   string
	DevMode bool
}

var siteData = SiteData{
	Title:   "connors-go-htmx-starter",
	DevMode: true,
}

func (s *Server) routes() {
	// static content embed fs in server struct
	s.Router.Handle("/static/*", http.FileServer(http.FS(s.FileSystem)))

	s.Router.HandleFunc("/", s.handleIndex(s.compile("index", indexView, nil)))
}

func (s *Server) handleIndex(tmpl *template.Template) http.HandlerFunc {
	type ViewData struct {
		SiteData  SiteData
		PageTitle string
		Username  string
		Option    string
		Offset    string
	}

	data := ViewData{
		SiteData:  siteData,
		PageTitle: "Index",
		Username:  "connorkuljis",
		Option:    "dev",
		Offset:    "-2",
	}

	return func(w http.ResponseWriter, r *http.Request) {
		tmpl.ExecuteTemplate(w, "root", data)
	}
}
