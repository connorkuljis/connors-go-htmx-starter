package fragments

import (
	"io/fs"
	"path/filepath"
)

type Templates struct {
	Base       map[string]string
	Components map[string]string
	Views      map[string]string
}

var Tmpls Templates

func InitFragments(filesystem fs.FS, templatesDir string) (Templates, error) {
	var err error

	Tmpls.Base, err = loadfiles(filesystem, templatesDir)
	if err != nil {
		return Tmpls, err
	}

	Tmpls.Components, err = loadfiles(filesystem, filepath.Join(templatesDir, "components"))
	if err != nil {
		return Tmpls, err
	}

	Tmpls.Views, err = loadfiles(filesystem, filepath.Join(templatesDir, "views"))
	if err != nil {
		return Tmpls, err
	}

	return Tmpls, nil
}

func loadfiles(filesystem fs.FS, path string) (map[string]string, error) {
	mm := make(map[string]string)
	files, err := fs.ReadDir(filesystem, path)
	if err != nil {
		return mm, err
	}
	for _, file := range files {
		if !file.IsDir() {
			name := file.Name()
			path := filepath.Join(path, name)
			mm[name] = path
		}
	}
	return mm, nil
}

// getIndexTemplate parses joined base and index view templates.
func IndexTemplate() []string {
	view := []string{
		Tmpls.Base["root.html"],
		Tmpls.Base["layout.html"],
		Tmpls.Base["head.html"],
		Tmpls.Components["footer.html"],
		Tmpls.Components["dev-tool.html"],
		Tmpls.Components["nav.html"],
		Tmpls.Components["header.html"],
		Tmpls.Views["index.html"],
	}

	return view
}
