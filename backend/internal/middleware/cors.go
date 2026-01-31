package middleware

import (
	"net/http"
	"os"
	"strings"
)

// getAllowedOrigin determines which origin to allow for the current request.
// It reads a comma-separated list of allowed origins from the CORS_ALLOWED_ORIGINS
// environment variable. If CORS_ALLOWED_ORIGINS is empty, it falls back to
// reflecting the request's Origin header (permissive, suitable for development).
func getAllowedOrigin(r *http.Request) string {
	origin := r.Header.Get("Origin")
	if origin == "" {
		return ""
	}

	allowed := os.Getenv("CORS_ALLOWED_ORIGINS")
	if allowed == "" {
		// Permissive mode: reflect the request origin instead of using "*".
		return origin
	}

	for _, o := range strings.Split(allowed, ",") {
		if strings.TrimSpace(o) == origin {
			return origin
		}
	}

	return ""
}

// CORSMiddleware adds CORS headers to the response.
func CORSMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if allowedOrigin := getAllowedOrigin(r); allowedOrigin != "" {
			w.Header().Set("Access-Control-Allow-Origin", allowedOrigin)
			// Inform caches that the response varies based on Origin.
			w.Header().Add("Vary", "Origin")
		}
		w.Header().Set("Access-Control-Allow-Methods", "GET, POST, DELETE, OPTIONS")
		w.Header().Set("Access-Control-Allow-Headers", "Content-Type, Authorization")

		if r.Method == http.MethodOptions {
			w.WriteHeader(http.StatusOK)
			return
		}

		next.ServeHTTP(w, r)
	})
}
