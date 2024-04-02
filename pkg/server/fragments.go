package server

import (
	"io/fs"
	"path/filepath"
)

type Fragments struct {
	Base       map[string]string
	Components map[string]string
	Views      map[string]string
}

// InitFragments traverses the base, components and views directory in the given filesystem and returns a Fragments structure, or an error if an error occurs.
func InitFragments(filesystem fs.FS, templatesDir string) (Fragments, error) {
	var tmpls Fragments
	var err error

	tmpls.Base, err = loadfiles(filesystem, templatesDir)
	if err != nil {
		return tmpls, err
	}

	tmpls.Components, err = loadfiles(filesystem, filepath.Join(templatesDir, "components"))
	if err != nil {
		return tmpls, err
	}

	tmpls.Views, err = loadfiles(filesystem, filepath.Join(templatesDir, "views"))
	if err != nil {
		return tmpls, err
	}

	return tmpls, nil
}

// loadfiles reads the filepath of all regular files into a map, keyed by the filename
func loadfiles(filesystem fs.FS, path string) (map[string]string, error) {
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

// IndexTemplate returns all the html fragments in the server for the index view.
func (s *Server) IndexTemplate() []string {
	return []string{
		s.Fragments.Base["root.html"],
		s.Fragments.Base["layout.html"],
		s.Fragments.Base["head.html"],
		s.Fragments.Components["footer.html"],
		s.Fragments.Components["dev-tool.html"],
		s.Fragments.Components["nav.html"],
		s.Fragments.Components["header.html"],
		s.Fragments.Views["index.html"],
	}
}
