package gateway

import (
	"api-gateway/config"
	"net/http"

	"github.com/gorilla/mux"
)

func NewRouter(cfg *config.Config) http.Handler {
	r := mux.NewRouter()

	// Middleware example: simple auth check
	r.Use(authMiddleware)

	// Proxy routes to user service
	r.PathPrefix("/user").HandlerFunc(proxyToUserService(cfg.UserServiceURL))

	return r
}

func authMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		token := r.Header.Get("Authorization")
		if token == "" {
			http.Error(w, "Unauthorized", http.StatusUnauthorized)
			return
		}
		// TODO: Validate token properly
		next.ServeHTTP(w, r)
	})
}
