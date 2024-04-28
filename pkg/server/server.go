package server

import (
	"bytes"
	"fmt"
	"html/template"
	"io/fs"
	"log"
	"net/http"
	"os"
	"path/filepath"
)

const Base = "base"

// Server encapsulates all dependencies for the web Server.
// HTTP handlers access information via receiver types.
type Server struct {
	FileSystem           fs.FS
	StaticContentHandler http.Handler
	MuxRouter            *http.ServeMux
	TemplateFragments    TemplateFragments

	Port string
}

type TemplateFragments map[string]map[string]string

// NewServer returns a new pointer Server struct.
//
// Server encapsulates all dependencies for the web Server.
// HTTP handlers access information via receiver types.
func NewServer(fileSystem fs.FS, port, templatesPath, staticPath string) (*Server, error) {
	templateFragments, err := ExtractTemplateFragmentsFromFilesystem(fileSystem, templatesPath)
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
		Port:                 port,
		StaticContentHandler: http.FileServer(http.FS(scfs)),
		TemplateFragments:    templateFragments,
	}
	return s, nil
}

func (s *Server) ListenAndServe() error {
	log.Println("[ 💿 Spinning up server on http://localhost:" + s.Port + " ]")
	if err := http.ListenAndServe(":"+s.Port, s.MuxRouter); err != nil {
		return fmt.Errorf("Error starting server: %w", err)
	}
	return nil
}

// ExtractTemplateFragmentsFromFilesystem traverses the base, components and views directory in the given filesystem and returns a Fragments structure, or an error if an error occurs.
func ExtractTemplateFragmentsFromFilesystem(filesystem fs.FS, templatesPath string) (TemplateFragments, error) {
	var err error
	templateFragments := make(TemplateFragments, 0)
	regularFiles, err := readRegularFiles(filesystem, templatesPath)
	if err != nil {
		return templateFragments, err
	}
	templateFragments[Base] = regularFiles
	topLevelDirs, err := readTopLevelDirs(filesystem, templatesPath)
	if err != nil {
		return templateFragments, err
	}
	for _, dir := range topLevelDirs {
		k := dir.Name()
		targetPath := filepath.Join(templatesPath, k)
		v, err := readRegularFiles(filesystem, targetPath)
		if err != nil {
			return templateFragments, err
		}
		templateFragments[k] = v
	}
	return templateFragments, nil
}

func readTopLevelDirs(filesystem fs.FS, targetPath string) ([]os.DirEntry, error) {
	var topLevelDirs []os.DirEntry
	dirs, err := fs.ReadDir(filesystem, targetPath)
	if err != nil {
		return topLevelDirs, err
	}
	for _, dir := range dirs {
		if dir.IsDir() {
			topLevelDirs = append(topLevelDirs, dir)
		}
	}
	return topLevelDirs, nil
}

// readRegularFiles takes a filesystem and a targetPath and returns a map[string]string, or an error if an error occurs.
//
// First a new empty map is created, then we iterate over each dir. If the dir
// is a regular file, we assign a new entry into the map where the filename is the
// key and the joined filepath string is the value
func readRegularFiles(filesystem fs.FS, targetPath string) (map[string]string, error) {
	newMap := make(map[string]string)
	dirs, err := fs.ReadDir(filesystem, targetPath)
	if err != nil {
		return newMap, err
	}
	for _, dir := range dirs {
		if dir.Type().IsRegular() {
			k := dir.Name()
			v := filepath.Join(targetPath, k)
			newMap[k] = v
		}
	}
	return newMap, nil
}

// buildTemplates is a fast way to parse a collection of templates in the server filesystem.
//
// template files are provided as strings to be parsed from the filesystem
func (s *Server) BuildTemplates(name string, funcs template.FuncMap, templates ...string) *template.Template {
	tmpl := template.New(name)
	if funcs != nil {
		tmpl.Funcs(funcs)
	}
	tmpl, err := tmpl.ParseFS(s.FileSystem, templates...)
	if err != nil {
		err = fmt.Errorf("Error building template name='%s': %w", name, err)
		log.Fatal(err)
	}
	return tmpl
}

// safeTmplParse executes a given template to a bytes buffer. It returns the resulting buffer or nil, err if any error occurred.
//
// Templates are checked for missing keys to prevent partial data being written to the writer.
func SafeTmplExec(tmpl *template.Template, name string, data any) ([]byte, error) {
	var bufBytes bytes.Buffer
	tmpl.Option("missingkey=error")
	err := tmpl.ExecuteTemplate(&bufBytes, name, data)
	if err != nil {
		log.Print(err)
		return bufBytes.Bytes(), err
	}
	return bufBytes.Bytes(), nil
}

// sendHTML writes a buffer a response writer as html
func SendHTML(w http.ResponseWriter, bufBytes []byte) {
	w.Header().Set("Content-Type", "text/html; charset=UTF-8")
	_, err := w.Write(bufBytes)
	if err != nil {
		log.Println(err)
	}
}
