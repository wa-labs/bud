package service

import (
	"net/http"
	"strings"
)

// BasicAuth ...
func BasicAuth(auth string) func(http.Handler) http.Handler {
	auth = strings.TrimSpace(auth)
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			token := r.Header.Get("Basic-Auth-Token")
			if token == "" || token != auth {
				http.Error(w, "Unauthorized", 401)
				return
			}
			next.ServeHTTP(w, r)
		})
	}
}
