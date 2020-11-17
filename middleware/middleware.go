package middleware

import (
	"log"
	"net/http"
	"time"
)

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

	return
}

// LoggingMiddleware logs the incoming HTTP request & its duration.
func LoggingMiddleware(logger *log.Logger) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		fn := func(w http.ResponseWriter, r *http.Request) {
			// Manages any panic that may be raised inside the handler next.
			defer func() {
				if err := recover(); err != nil {
					w.WriteHeader(http.StatusInternalServerError)
					logger.Fatal(err)
				}
			}()

			// Stores start time in order to get duration after the handler returns.
			wrapped := wrapResponseWriter(w)
			start := time.Now()
			logger.Printf("Starting handler for endpoint [%v %v]\n", r.Method, r.URL.EscapedPath())
			next.ServeHTTP(wrapped, r)
			logger.Printf("Handler for endpoint [%v %v] returned [%v]. Duration: %v\n", r.Method, r.URL.EscapedPath(), wrapped.status, time.Since(start))
		}

		return http.HandlerFunc(fn)
	}
}
