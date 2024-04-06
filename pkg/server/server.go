package server

import (
	"bytes"
	"errors"
	"html/template"
	"io/fs"
	"log"
	"net/http"
	"path/filepath"
)

// Server encapsulates all dependencies for the web Server.
// HTTP handlers access information via receiver types.
type Server struct {
	FileSystem fs.FS
	Router     *http.ServeMux
	Fragments  Fragments

	Port string
}

type Fragments struct {
	Base       map[string]string
	Components map[string]string
	Views      map[string]string
}

const (
	TemplatesDir  = "www/templates"
	StaticDir     = "www/static"
	ComponentsDir = "components"
	ViewsDir      = "views"
)

// NewServer returns a new pointer Server struct.
//
// Server encapsulates all dependencies for the web Server.
// HTTP handlers access information via receiver types.
func NewServer(fileSystem fs.FS, port string) (*Server, error) {
	fragments, err := InitFragments(fileSystem)
	if err != nil {
		return nil, err
	}

	s := &Server{
		FileSystem: fileSystem,
		Router:     http.NewServeMux(),
		Port:       port,
		Fragments:  fragments,
	}

	return s, nil
}

// Routes instatiates http Handlers and associated patterns on the server.
func (s *Server) Routes() error {
	scfs, err := fs.Sub(s.FileSystem, StaticDir) // static content sub fs from the server's embedded fs
	if err != nil {
		return err
	}

	s.Router.Handle("/static/", http.StripPrefix("/static/", http.FileServer(http.FS(scfs))))
	s.Router.HandleFunc("/", s.HandleIndex())

	return nil
}

func (s *Server) ListenAndServe() error {
	log.Println("[ 💿 Spinning up server on http://localhost:" + s.Port + " ]")
	if err := http.ListenAndServe(":"+s.Port, s.Router); err != nil {
		log.Println("Error starting server.")
		return err
	}

	return nil
}

// InitFragments traverses the base, components and views directory in the given filesystem and returns a Fragments structure, or an error if an error occurs.
func InitFragments(filesystem fs.FS) (Fragments, error) {
	var tmpls Fragments
	var err error

	templatesPath := TemplatesDir
	componentsPath := filepath.Join(TemplatesDir, ComponentsDir)
	viewsPath := filepath.Join(TemplatesDir, ViewsDir)

	tmpls.Base, err = mapNameToPath(filesystem, templatesPath)
	if err != nil {
		return tmpls, err
	}

	tmpls.Components, err = mapNameToPath(filesystem, componentsPath)
	if err != nil {
		return tmpls, err
	}

	tmpls.Views, err = mapNameToPath(filesystem, viewsPath)
	if err != nil {
		return tmpls, err
	}

	return tmpls, nil
}

// mapNameToPath reads the filepath of all regular files into a map, keyed by the filename
func mapNameToPath(filesystem fs.FS, path string) (map[string]string, error) {
	mm := make(map[string]string)

	files, err := fs.ReadDir(filesystem, path)
	if err != nil {
		return mm, err
	}

	for _, file := range files {
		if !file.IsDir() {
			name := file.Name()               // "index.html"
			path := filepath.Join(path, name) // "www/static/templates/views/index/html"
			mm[name] = path                   // "index.html" => "www/static/templates/views/index/html"
		}
	}

	return mm, nil
}

// buildTemplates is a fast way to parse a collection of templates in the server filesystem.
//
// template files are provided as strings to be parsed from the filesystem
func (s *Server) BuildTemplates(name string, funcs template.FuncMap, templates ...string) *template.Template {
	for _, template := range templates {
		if template == "" {
			log.Fatal(errors.New("Error building template for (" + name + "): an empty template was detected..."))
		}
	}
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

// safeTmplParse executes a given template to a bytes buffer. It returns the resulting buffer or nil, err if any error occurred.
//
// Templates are checked for missing keys to prevent partial data being written to the writer.
func SafeTmplExec(tmpl *template.Template, name string, data any) ([]byte, error) {
	var buf bytes.Buffer
	tmpl.Option("missingkey=error")
	err := tmpl.ExecuteTemplate(&buf, name, data)
	if err != nil {
		log.Print(err)
		return buf.Bytes(), err
	}
	return buf.Bytes(), nil
}

// sendHTML writes a buffer a response writer as html
func SendHTML(w http.ResponseWriter, buf []byte) {
	w.Header().Set("Content-Type", "text/html; charset=UTF-8")
	_, err := w.Write(buf)
	if err != nil {
		log.Println(err)
	}
}
