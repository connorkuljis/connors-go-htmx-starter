package handlers

import (
	"html/template"
	"net/http"
	"time"

	"github.com/connorkuljis/connors-go-htmx-starter/pkg/webserver"
)

func HandleIndex(s *webserver.WebServer, t *template.Template) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		parcel := map[string]any{"Time": time.Now()}
		htmlBytes, err := webserver.SafeTmplExec(t, "root", parcel)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		webserver.SendHTML(w, htmlBytes)
	}
}
