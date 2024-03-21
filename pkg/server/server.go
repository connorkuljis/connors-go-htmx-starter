package server

import (
	"io/fs"
	"net/http"
	"path/filepath"
)

// Server encapsulates all dependencies for the web Server.
// HTTP handlers access information via receiver types.
type Server struct {
	FileSystem fs.FS // in-memory or disk
	Router     *http.ServeMux
	Templates  Templates

	Port         string
	StaticDir    string // location of static assets
	TemplatesDir string // location of html templates, makes template parsing less verbose.
}

type Templates struct {
	BaseLayout BaseLayout
	Components Components
	Views      Views
}

type BaseLayout struct {
	Root   string
	Head   string
	Layout string
}

type Components struct {
	DevTool string
	Header  string
	Nav     string
	Footer  string
}

type Views struct {
	Index    string
	Projects string
}

type AppData struct {
	Title   string
	DevMode bool
}

const (
	StaticDirName    = "web/static"
	TemplatesDirName = "web/templates"
)

func NewServer(fileSystem fs.FS, port string) *Server {
	return &Server{
		FileSystem:   fileSystem,
		Router:       http.NewServeMux(),
		Port:         port,
		TemplatesDir: TemplatesDirName,
		StaticDir:    StaticDirName,
		Templates:    loadTemplates(TemplatesDirName),
	}
}

// Routes instatiates http Handlers and associated patterns on the server.
func (s *Server) Routes() {
	s.Router.Handle("/static/", http.FileServer(http.FS(s.FileSystem)))
	s.Router.HandleFunc("/", s.HandleIndex())
}

func loadTemplates(dir string) Templates {
	return Templates{
		BaseLayout: BaseLayout{
			Root:   filepath.Join(dir, "root.html"),
			Head:   filepath.Join(dir, "head.html"),
			Layout: filepath.Join(dir, "layout.html"),
		},

		Components: Components{
			DevTool: filepath.Join(dir, "components", "dev-tool.html"),
			Header:  filepath.Join(dir, "components", "header.html"),
			Nav:     filepath.Join(dir, "components", "nav.html"),
			Footer:  filepath.Join(dir, "components", "footer.html"),
		},
		Views: Views{
			Index: filepath.Join(dir, "views", "index.html"),
		},
	}
}
