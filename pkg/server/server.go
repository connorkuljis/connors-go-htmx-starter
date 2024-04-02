package server

import (
	"html/template"
	"io/fs"
	"log"
	"net/http"
)

// Server encapsulates all dependencies for the web Server.
// HTTP handlers access information via receiver types.
type Server struct {
	FileSystem fs.FS // in-memory or disk
	Router     *http.ServeMux
	Fragments  Fragments

	Port         string
	StaticDir    string // location of static assets
	TemplatesDir string // location of html templates, makes template parsing less verbose.
}

// NewServer returns a new pointer Server struct.
//
// Server encapsulates all dependencies for the web Server.
// HTTP handlers access information via receiver types.
func NewServer(fileSystem fs.FS, port string, static string, templates string) (*Server, error) {
	s := &Server{
		FileSystem:   fileSystem,
		Router:       http.NewServeMux(),
		Port:         port,
		TemplatesDir: templates,
		StaticDir:    static,
	}

	var err error
	s.Fragments, err = InitFragments(s.FileSystem, s.TemplatesDir)
	if err != nil {
		return s, err
	}

	return s, nil
}

// Routes instatiates http Handlers and associated patterns on the server.
func (s *Server) Routes() {
	s.Router.Handle("/static/", http.FileServer(http.FS(s.FileSystem)))
	s.Router.HandleFunc("/", s.HandleIndex())
}

// buildTemplates is a fast way to parse a collection of templates in the server filesystem.
//
// template files are provided as strings to be parsed from the filesystem
func (s *Server) BuildTemplates(name string, funcs template.FuncMap, templates ...string) *template.Template {
	// give the template a name
	tmpl := template.New(name)

	// custom template functions if exists
	if funcs != nil {
		tmpl.Funcs(funcs)
	}

	// generate a template from the files in the server fs (usually embedded)
	tmpl, err := tmpl.ParseFS(s.FileSystem, templates...)
	if err != nil {
		log.Fatal(err)
	}

	return tmpl
}
