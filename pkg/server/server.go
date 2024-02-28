package server

import (
	"io/fs"
	"net/http"
)

// Server encapsulates all dependencies for the web server.
// HTTP handlers access information via receiver types.
type server struct {
	Router     *http.ServeMux
	FileSystem fs.FS // in-memory or disk

	Port         string
	TemplatesDir string // location of html templates, makes template parsing less verbose.
	StaticDir    string // location of static assets

	Title          string
	DevModeEnabled bool
}

func NewServer(router *http.ServeMux, fileSystem fs.FS, port string, templatesDir string, staticDir string, title string, devModeEnabled bool) server {
	return server{
		Router:         router,
		FileSystem:     fileSystem,
		Port:           port,
		TemplatesDir:   templatesDir,
		StaticDir:      staticDir,
		Title:          title,
		DevModeEnabled: devModeEnabled,
	}
}
