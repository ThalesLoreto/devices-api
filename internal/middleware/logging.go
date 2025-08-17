package middleware

import (
	"log"
	"net/http"
	"time"
)

// responseWriter wraps http.ResponseWriter to capture status code
type responseWriter struct {
	http.ResponseWriter
	statusCode int
}

// WriteHeader captures the status code
func (rw *responseWriter) WriteHeader(code int) {
	rw.statusCode = code
	rw.ResponseWriter.WriteHeader(code)
}

// LoggingMiddleware logs details about each HTTP request
func LoggingMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		start := time.Now()

		// Wrap the response writer to capture status code
		wrapped := &responseWriter{
			ResponseWriter: w,
			statusCode:     http.StatusOK, // Default status code
		}

		// Call the next handler
		next.ServeHTTP(wrapped, r)

		// Log the request details
		duration := time.Since(start)
		log.Printf(
			"[%s] %s %s - Status: %d - Duration: %v - User-Agent: %s - Remote: %s",
			time.Now().Format("2006-01-02 15:04:05"),
			r.Method,
			r.RequestURI,
			wrapped.statusCode,
			duration,
			r.UserAgent(),
			r.RemoteAddr,
		)
	})
}
