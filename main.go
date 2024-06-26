package main

import (
	"embed"
	"fmt"
	"log"
	"log/slog"
	"net/http"
	"os"

	"github.com/connorkuljis/connors-go-htmx-starter/pkg/handlers"
	"github.com/connorkuljis/connors-go-htmx-starter/pkg/middleware"
	"github.com/connorkuljis/connors-go-htmx-starter/pkg/views"
	"github.com/connorkuljis/connors-go-htmx-starter/pkg/webserver"
)

const (
	TemplatesDir = "www/templates"
	StaticDir    = "www/static"
	Port         = "8080"
)

//go:embed www
var embedFS embed.FS

func main() {
	logger := slog.New(slog.NewJSONHandler(os.Stdout, nil))

	s, err := webserver.NewWebServer(embedFS, logger, Port, TemplatesDir, StaticDir)
	if err != nil {
		log.Fatal(err)
	}

	// Routes
	s.MuxRouter.Handle("/static/", http.StripPrefix("/static/", s.StaticFileServer))
	s.MuxRouter.HandleFunc("/", middleware.Logging(s.Logger, handlers.HandleIndex(s, views.IndexView(s))))

	fmt.Println("[ 💿 Spinning up server on http://localhost:" + s.Port + " ]")

	if err := http.ListenAndServe(":"+s.Port, s.MuxRouter); err != nil {
		log.Fatal(fmt.Errorf("Error starting server: %w", err))
	}
}
