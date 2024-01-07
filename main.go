package main

import (
	"embed"
	"io/fs"
	"log"
	"net/http"
	"path/filepath"
	"text/template"

	"github.com/go-chi/chi/v5"
	"github.com/gorilla/sessions"
)

const SessionName = "session"

// Server encapsulates all dependencies for the web server.
// HTTP handlers access information via receiver types.
type Server struct {
	Port         string
	Router       *chi.Mux
	TemplatesDir string // location of html templates, makes template parsing less verbose.
	StaticDir    string // location of static assets
	FileSystem   fs.FS  // in-memory or disk
	Sessions     *sessions.CookieStore
}

//go:embed templates/* static/*
var inMemoryFS embed.FS

type HTMLFile string

const (
	RootHTML   HTMLFile = "root.html"
	HeadHTML   HTMLFile = "head.html"
	LayoutHTML HTMLFile = "layout.html"
	HeroHTML   HTMLFile = "components/hero.html"
	FooterHTML HTMLFile = "components/footer.html"
)

func main() {
	port := "8080"
	// router := http.NewServeMux()
	router := chi.NewMux()
	store := sessions.NewCookieStore([]byte("special_key"))
	templateDir := "templates"
	staticDir := "static"

	log.Println("[ 💿 Spinning up server on http://localhost:" + port + " ]")

	s := Server{
		Router:       router,
		Port:         port,
		TemplatesDir: templateDir,
		StaticDir:    staticDir,
		FileSystem:   inMemoryFS,
		Sessions:     store,
	}

	s.routes()

	err := http.ListenAndServe(":"+s.Port, s.Router)
	if err != nil {
		panic(err)
	}
}

func compileTemplates(s *Server, files []HTMLFile) *template.Template {
	var filenames []string
	for i := range files {
		currentFilename := string(files[i])
		filenames = append(filenames, filepath.Join(s.TemplatesDir, currentFilename))
	}

	tmpl, err := template.ParseFS(s.FileSystem, filenames...)
	if err != nil {
		panic(err)
	}

	return tmpl
}

func (s *Server) routes() {
	s.Router.Handle("/static/*", http.FileServer(http.FS(s.FileSystem)))
	s.Router.HandleFunc("/", s.handleIndex())
}

func (s *Server) handleIndex() http.HandlerFunc {
	type PageData struct {
		Username string
		Option   string
		Offset   string
	}

	var indexHTML = []HTMLFile{
		RootHTML,
		HeadHTML,
		LayoutHTML,
		HeroHTML,
		FooterHTML,
	}

	tmpl := compileTemplates(s, indexHTML)

	return func(w http.ResponseWriter, r *http.Request) {
		tmpl.ExecuteTemplate(w, "root", nil)
	}
}
