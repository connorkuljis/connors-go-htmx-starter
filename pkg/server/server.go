package server

import (
	"bytes"
	"fmt"
	"html/template"
	"io/fs"
	"log"
	"log/slog"
	"net/http"
	"path/filepath"
)

type TemplateStore struct {
	Base       map[string]string
	Components map[string]string
	Views      map[string]string
}

// Server encapsulates all dependencies for the web Server.
// HTTP handlers access information via receiver types.
type Server struct {
	FileSystem           fs.FS
	StaticContentHandler http.Handler
	MuxRouter            *http.ServeMux
	Logger               *slog.Logger
	TemplateStore        TemplateStore

	Port string
}

// NewServer returns a new pointer Server struct.
//
// Server encapsulates all dependencies for the web Server.
// HTTP handlers access information via receiver types.
func NewServer(fileSystem fs.FS, logger *slog.Logger, port, templatesPath, staticPath string) (*Server, error) {
	templateStore, err := BuildTemplateStore(fileSystem, templatesPath)
	if err != nil {
		return nil, err
	}
	scfs, err := fs.Sub(fileSystem, staticPath)
	if err != nil {
		return nil, err
	}
	s := &Server{
		FileSystem:           fileSystem,
		MuxRouter:            http.NewServeMux(),
		Logger:               logger,
		Port:                 port,
		StaticContentHandler: http.FileServer(http.FS(scfs)),
		TemplateStore:        templateStore,
	}
	return s, nil
}

func BuildTemplateStore(filesystem fs.FS, templatesPath string) (TemplateStore, error) {
	var templateStore TemplateStore

	base, err := fs.ReadDir(filesystem, templatesPath)
	if err != nil {
		return templateStore, err
	}
	templateStore.Base = buildMap(base, templatesPath)

	componentsPath := filepath.Join(templatesPath, "components")
	components, err := fs.ReadDir(filesystem, componentsPath)
	if err != nil {
		return templateStore, err
	}
	templateStore.Components = buildMap(components, componentsPath)

	viewsPath := filepath.Join(templatesPath, "views")
	views, err := fs.ReadDir(filesystem, viewsPath)
	if err != nil {
		return templateStore, err
	}
	templateStore.Views = buildMap(views, viewsPath)

	return templateStore, nil
}

func buildMap(files []fs.DirEntry, path string) map[string]string {
	newMap := map[string]string{}

	for _, file := range files {
		filename := file.Name()
		if file.Type().IsRegular() {
			newMap[filename] = filepath.Join(path, filename)
		}
	}

	return newMap
}

// buildTemplates is a fast way to parse a collection of templates in the server filesystem.
//
// template files are provided as strings to be parsed from the filesystem
func (s *Server) ParseTemplates(name string, funcs template.FuncMap, templatefiles ...string) *template.Template {
	tmpl := template.New(name)
	if funcs != nil {
		tmpl.Funcs(funcs)
	}

	tmpl, err := tmpl.ParseFS(s.FileSystem, templatefiles...)
	if err != nil {
		err = fmt.Errorf("Error building template name='%s': %w", name, err)
		s.Logger.Error(err.Error())
		log.Fatal(err)
	}
	return tmpl
}

func (s *Server) ListenAndServe() error {
	log.Println("[ 💿 Spinning up server on http://localhost:" + s.Port + " ]")
	if err := http.ListenAndServe(":"+s.Port, s.MuxRouter); err != nil {
		return fmt.Errorf("Error starting server: %w", err)
	}
	return nil
}

//
// Utils
//
//

// safeTmplParse executes a given template to a bytes buffer. It returns the resulting buffer or nil, err if any error occurred.
//
// Templates are checked for missing keys to prevent partial data being written to the writer.
func (s *Server) SafeTmplExec(tmpl *template.Template, name string, data any) ([]byte, error) {
	var bufBytes bytes.Buffer
	tmpl.Option("missingkey=error")
	err := tmpl.ExecuteTemplate(&bufBytes, name, data)
	if err != nil {
		s.Logger.Error(err.Error())
		return bufBytes.Bytes(), err
	}
	return bufBytes.Bytes(), nil
}

// sendHTML writes a buffer a response writer as html
func (s *Server) SendHTML(w http.ResponseWriter, bufBytes []byte) {
	w.Header().Set("Content-Type", "text/html; charset=UTF-8")
	_, err := w.Write(bufBytes)
	if err != nil {
		s.Logger.Error(err.Error())
	}
}
