package server

// handlePartialProjects sends a template partial for HTMX requests, or sends a complete or "root" template containing the partial.
// Assume the HTMX caller uses `hx-push-url`.
// func (s *server) handlePartialProjects(tmpls map[string]*template.Template) http.HandlerFunc {
// 	data := map[string]interface{}{
// 		"SiteData":  s.GetSiteData(),
// 		"Username":  "connorkuljis",
// 		"PageTitle": "Projects Page!",
// 		"Projects":  []string{"block-cli", "green-tiles"},
// 	}

// 	return func(w http.ResponseWriter, r *http.Request) {
// 		isHTMXReq := r.Header.Get("HX-Request") != ""

// 		if isHTMXReq {
// 			buf, err := safeTmplParse(tmpls["projects"], "projects", data)
// 			if err != nil {
// 				http.Error(w, err.Error(), http.StatusInternalServerError)
// 				return
// 			}
// 			sendHTML(w, buf)
// 			return
// 		}

// 		buf, err := safeTmplParse(tmpls["full"], "root", data)
// 		if err != nil {
// 			http.Error(w, err.Error(), http.StatusInternalServerError)
// 			return
// 		}
// 		sendHTML(w, buf)
// 	}
// }
