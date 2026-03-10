package middleware

import (
	"net/http"
	"os"
)

type CacheControlMiddleware struct {
	development bool
}

func NewCacheControlMiddleware() *CacheControlMiddleware {
	devMode := os.Getenv("DEV") != "false" && os.Getenv("DEV") != "0"
	return &CacheControlMiddleware{development: devMode}
}

func (m *CacheControlMiddleware) Handler(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if m.development {
			w.Header().Set("Cache-Control", "no-cache, no-store, must-revalidate")
			w.Header().Set("Pragma", "no-cache")
			w.Header().Set("Expires", "0")
		}
		next.ServeHTTP(w, r)
	})
}
