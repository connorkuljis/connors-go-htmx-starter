package webserver

import (
	"bytes"
	"fmt"
	"html/template"
	"io/fs"
	"net/http"
)

func SafeTmplExec(t *template.Template, name string, data any) ([]byte, error) {
	var b bytes.Buffer
	t.Option("missingkey=error")
	err := t.ExecuteTemplate(&b, name, data)
	if err != nil {
		return b.Bytes(), err
	}
	return b.Bytes(), nil
}

func SendHTML(w http.ResponseWriter, bufBytes []byte) error {
	w.Header().Set("Content-Type", "text/html; charset=UTF-8")
	_, err := w.Write(bufBytes)
	if err != nil {
		return err
	}
	return nil
}

// MustParseTemplates constructs a new html/template template
func ParseTemplates(name string, funcMap template.FuncMap, fs fs.FS, templates ...string) (*template.Template, error) {
	t, err := template.New(name).Funcs(funcMap).ParseFS(fs, templates...)
	if err != nil {
		return nil, fmt.Errorf("Error building template name='%s': %w", name, err)
	}
	return t, nil
}
