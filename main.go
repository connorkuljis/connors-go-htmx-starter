package main

import (
	"embed"
	"io/fs"
	"log"
	"net/http"

	"github.com/gorilla/sessions"
)

// Server encapsulates all dependencies for the web server.
// HTTP handlers access information via receiver types.
type Server struct {
	Port         string
	Router       *http.ServeMux
	TemplatesDir string // location of html templates, makes template parsing less verbose.
	StaticDir    string // location of static assets
	FileSystem   fs.FS  // in-memory or disk
	Sessions     *sessions.CookieStore
}

//go:embed templates/* static/*
var embedFS embed.FS

const (
	port             = "8080"
	staticDirName    = "static"
	templatesDirName = "templates"
)

func main() {
	router := http.NewServeMux()
	store := sessions.NewCookieStore([]byte("special_key"))

	s := Server{
		Router:       router,
		Port:         port,
		TemplatesDir: templatesDirName,
		StaticDir:    staticDirName,
		FileSystem:   embedFS,
		Sessions:     store,
	}

	s.routes()

	log.Println("[ 💿 Spinning up server on http://localhost:" + s.Port + " ]")

	if err := http.ListenAndServe(":"+s.Port, s.Router); err != nil {
		log.Fatal("Error starting server", err)
	}
}
