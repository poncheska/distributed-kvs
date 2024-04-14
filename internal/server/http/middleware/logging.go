package middleware

import (
	"log/slog"
	"net/http"
	"time"
)

func WithLoggingMiddleware(h http.Handler, log *slog.Logger) http.Handler {
	return loggingMiddleware(h.ServeHTTP, log)
}

func loggingMiddleware(next http.HandlerFunc, log *slog.Logger) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		wrapped := wrapResponseWriter(w)

		start := time.Now()

		next.ServeHTTP(wrapped, r)

		log.Info(
			"incoming http request",
			"method", r.Method,
			"path", r.URL.Path,
			"duration", time.Since(start).String(),
			"status", changeDefaultStatus(wrapped.Status()),
		)
	}
}

// responseWriter is a minimal wrapper for http.ResponseWriter that allows the
// written HTTP status code to be captured for logging.
type responseWriter struct {
	http.ResponseWriter
	status      int
	wroteHeader bool
}

func wrapResponseWriter(w http.ResponseWriter) *responseWriter {
	return &responseWriter{ResponseWriter: w}
}

func (rw *responseWriter) Status() int {
	return rw.status
}

func (rw *responseWriter) WriteHeader(code int) {
	if rw.wroteHeader {
		return
	}

	rw.status = code
	rw.ResponseWriter.WriteHeader(code)
	rw.wroteHeader = true
}

func changeDefaultStatus(statusCode int) int {
	if statusCode == 0 {
		return http.StatusOK
	}

	return statusCode
}
