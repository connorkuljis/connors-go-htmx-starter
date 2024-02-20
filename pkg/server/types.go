package server

import (
	"io/fs"
	"net/http"
)

const (
	headHTML   = "templates/head.html"
	layoutHTML = "templates/layout.html"
	rootHTML   = "templates/root.html"
	heroHTML   = "templates/components/hero.html"
	footerHTML = "templates/components/footer.html"
	navHTML    = "templates/components/nav.html"
	indexHTML  = "templates/views/index.html"
)

type htmlFile string

type view []htmlFile

// Server encapsulates all dependencies for the web server.
// HTTP handlers access information via receiver types.
type server struct {
	Port         string
	Router       *http.ServeMux
	TemplatesDir string // location of html templates, makes template parsing less verbose.
	StaticDir    string // location of static assets
	FileSystem   fs.FS  // in-memory or disk
}

// Represents site information
type siteData struct {
	Title   string
	DevMode bool
}
