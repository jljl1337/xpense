package middleware

import (
	"fmt"
	"log/slog"
	"net/http"
	"time"

	"github.com/jljl1337/xpense/internal/env"
)

// responseWriter wraps http.ResponseWriter to capture status code
type responseWriter struct {
	http.ResponseWriter
	statusCode int
}

func (rw *responseWriter) WriteHeader(code int) {
	rw.statusCode = code
	rw.ResponseWriter.WriteHeader(code)
}

// Logging middleware logs all requests
func (m *MiddlewareProvider) Logging() Middleware {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			if r.URL.Path == "/health" && !env.LogHealthCheck {
				next.ServeHTTP(w, r)
				return
			}

			start := time.Now()

			// Create a response writer wrapper to capture status code
			wrapped := &responseWriter{ResponseWriter: w, statusCode: http.StatusOK}

			next.ServeHTTP(wrapped, r)

			duration := time.Since(start)
			slog.Info(
				fmt.Sprintf(
					"%-6s %d %s %s %s",
					r.Method,
					wrapped.statusCode,
					r.RequestURI,
					duration,
					r.RemoteAddr,
				),
			)
		})
	}
}
