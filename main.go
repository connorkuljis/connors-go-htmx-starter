package main

import (
	"embed"
	"io/fs"
	"log"
	"net/http"
	"text/template"

	"github.com/go-chi/chi/v5"
	"github.com/gorilla/sessions"
)

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
var embedFS embed.FS

type HTMLFile string

const (
	TemplatesDirName = "templates"
	StaticDirName    = "static"
	Port             = "8080"

	RootHTML   HTMLFile = "templates/root.html"
	HeadHTML   HTMLFile = "templates/head.html"
	LayoutHTML HTMLFile = "templates/layout.html"
	HeroHTML   HTMLFile = "templates/components/hero.html"
	FooterHTML HTMLFile = "templates/components/footer.html"
)

func main() {
	router := chi.NewMux()
	store := sessions.NewCookieStore([]byte("special_key"))

	s := Server{
		Router:       router,
		Port:         Port,
		TemplatesDir: TemplatesDirName,
		StaticDir:    StaticDirName,
		FileSystem:   embedFS,
		Sessions:     store,
	}

	s.Router.Handle("/static/*", http.FileServer(http.FS(s.FileSystem)))
	s.Router.HandleFunc("/", s.handleIndex())

	log.Println("[ 💿 Spinning up server on http://localhost:" + s.Port + " ]")

	err := http.ListenAndServe(":"+s.Port, s.Router)
	if err != nil {
		panic(err)
	}
}

func (s *Server) handleIndex() http.HandlerFunc {
	type ViewData struct {
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

	tmpl := s.compileTemplates("index.html", indexHTML, nil)

	return func(w http.ResponseWriter, r *http.Request) {
		tmpl.ExecuteTemplate(w, "root", nil)
	}
}

func (s *Server) compileTemplates(name string, targets []HTMLFile, funcs template.FuncMap) *template.Template {
	tmpl := template.New(name)

	if funcs != nil {
		tmpl.Funcs(funcs)
	}

	var matchedTargets []string
	for _, file := range targets {
		matchedTargets = append(matchedTargets, string(file))
	}

	tmpl, err := tmpl.ParseFS(s.FileSystem, matchedTargets...)
	if err != nil {
		log.Fatal(err)
	}

	return tmpl
}
