package middleware

import (
	"log/slog"
	"net/http"
)

func Logging(logger *slog.Logger, next http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		logger.Info(
			"incoming_request",
			"method", r.Method,
			"path", r.URL.Path,
		)
		next(w, r)
	}
}
