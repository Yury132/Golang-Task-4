package http

import (
	"net/http"

	"github.com/Yury132/Golang-Task-4/internal/transport/http/handlers"
	"github.com/gorilla/mux"
)

func InitRoutes(h *handlers.Handler) *mux.Router {
	r := mux.NewRouter()

	r.HandleFunc("/", h.Home).Methods(http.MethodGet)

	r.HandleFunc("/users-list", h.GetUsersList).Methods(http.MethodGet)

	http.Handle("/", r)

	return r
}
