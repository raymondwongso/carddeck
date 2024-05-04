// Package middleware contains various middleware needed by HTTP REST API Server
package middleware

import (
	"net/http"
)

// HeaderMiddleware inject various response header needed by the client
// typically for CORS checking and content type
// but for now only set content type to application/json
func HeaderMiddleware(h http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		h.ServeHTTP(w, r)
	})
}
