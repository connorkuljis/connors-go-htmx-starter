package server

import (
	"log"
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
	s.Router.HandleFunc("/", middleWareEx1(s.handleIndex(s.indexViewTmpl(), sData)))
}

// outputs the request method type
func middleWareEx1(next http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		log.Printf("Request method: %s", r.Method)
		next.ServeHTTP(w, r)
	}
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
