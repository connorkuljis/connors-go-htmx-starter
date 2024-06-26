package webserver

import (
	"fmt"
	"io/fs"
	"log/slog"
	"net/http"
	"path/filepath"
)

// WebServer encapsulates all dependencies for the web WebServer.
type WebServer struct {
	MuxRouter        *http.ServeMux
	FileSystem       fs.FS
	StaticFileServer http.Handler
	Logger           *slog.Logger
	TemplateStore    TemplateStore

	Debug bool
	Port  string
}

// TemplateStore stores paths of HTML files categorized into base, components, and views.
type TemplateStore struct {
	Base       map[string]string
	Components map[string]string
	Views      map[string]string
}

// NewWebServer returns a new pointer to a WebServer struct.
// It initializes the server with the provided file system, logger, port, templates path, and static files path.
// If there is an error during initialization, it returns the error.
func NewWebServer(fsys fs.FS, logger *slog.Logger, port, templatesPath, staticPath string) (*WebServer, error) {
	router := http.NewServeMux() // Create a new HTTP request multiplexer.

	// Initialize the TemplateStore with the provided file system and templates path.
	ts, err := newTemplateStore(fsys, templatesPath)
	if err != nil {
		return nil, err
	}

	// Create a sub file system for serving static files.
	scfs, err := fs.Sub(fsys, staticPath)
	if err != nil {
		return nil, err
	}

	// Create a file server handler for serving static files.
	staticFileServer := http.FileServer(http.FS(scfs))

	// Initialize the WebServer struct with the provided dependencies.
	s := &WebServer{
		MuxRouter:        router,
		FileSystem:       fsys,
		Logger:           logger,
		Port:             port,
		StaticFileServer: staticFileServer,
		TemplateStore:    ts,
	}

	return s, nil
}

// newTemplateStore creates a TemplateStore structure from a given filesystem and templates file path.
//
// It organizes templates into three categories: base, components, and views.
// The function returns a populated TemplateStore and any error encountered during the process.
func newTemplateStore(filesystem fs.FS, templatesFilePath string) (TemplateStore, error) {
	var ts TemplateStore

	// Define paths for each template category
	pathBase := templatesFilePath
	pathComponents := filepath.Join(templatesFilePath, "components")
	pathViews := filepath.Join(templatesFilePath, "views")

	// Read directory contents for each template category
	base, err := fs.ReadDir(filesystem, pathBase)
	if err != nil {
		return ts, fmt.Errorf("failed to read base templates directory: %w", err)
	}
	components, err := fs.ReadDir(filesystem, pathComponents)
	if err != nil {
		return ts, fmt.Errorf("failed to read components directory: %w", err)
	}
	views, err := fs.ReadDir(filesystem, pathViews)
	if err != nil {
		return ts, fmt.Errorf("failed to read views directory: %w", err)
	}

	// Create maps of filename to filepath for each template category
	baseMap := mapNameToLocation(templatesFilePath, base)
	componentsMap := mapNameToLocation(pathComponents, components)
	viewsMap := mapNameToLocation(pathViews, views)

	// Populate the TemplateStore with the created maps
	ts = TemplateStore{
		Base:       baseMap,
		Components: componentsMap,
		Views:      viewsMap,
	}

	return ts, nil
}

// mapNameToLocation creates a map that associates file names with their full paths.
//
// It takes a parent path and a slice of directory entries, and returns a map where
// the keys are file names and the values are the corresponding full file paths.
// Only regular files are included in the resulting map.
func mapNameToLocation(parentPath string, files []fs.DirEntry) map[string]string {
	store := map[string]string{}

	for _, file := range files {
		if file.Type().IsRegular() {
			name := file.Name()
			location := filepath.Join(parentPath, name)

			store[name] = location
		}
	}

	return store
}
