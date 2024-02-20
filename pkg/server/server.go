package server

import (
	"io/fs"
	"net/http"
)

// Server encapsulates all dependencies for the web server.
// HTTP handlers access information via receiver types.
type server struct {
	Port         string
	Router       *http.ServeMux
	TemplatesDir string // location of html templates, makes template parsing less verbose.
	StaticDir    string // location of static assets
	FileSystem   fs.FS  // in-memory or disk
}

func NewServer(port string, router *http.ServeMux, templatesDir string, staticDir string, fileSystem fs.FS) server {
	return server{
		Router:       router,
		Port:         port,
		TemplatesDir: templatesDir,
		StaticDir:    staticDir,
		FileSystem:   fileSystem,
	}
}
