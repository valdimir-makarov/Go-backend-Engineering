package gateway

import (
	"net/http"

	"github.com/valdimir-makarov/Go-backend-Engineering/chat-service/api-gateway/config"

	"github.com/gorilla/mux"
)

func NewRouter(cfg *config.Config) http.Handler {
	r := mux.NewRouter()

	// Route to auth service
	r.Handle("/register", Proxy(cfg.AuthServiceURL+"/register")).Methods("POST")
	r.Handle("/login", Proxy(cfg.AuthServiceURL+"/login")).Methods("POST")

	// Add more routes here...

	return r
}
