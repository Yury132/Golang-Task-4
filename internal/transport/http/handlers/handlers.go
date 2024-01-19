package handlers

import (
	"context"
	"encoding/json"
	"html/template"
	"net/http"

	"github.com/Yury132/Golang-Task-4/internal/models"
	"github.com/Yury132/Golang-Task-4/internal/service"
	"github.com/rs/zerolog"
)

type Service interface {
	GetUsersList(ctx context.Context) ([]models.User, error)
	HandleUser(ctx context.Context, name string, email string) error
}

type Handler struct {
	log     zerolog.Logger
	service Service
}

// Стартовая страница
func (h *Handler) Home(w http.ResponseWriter, r *http.Request) {

	tmpl, err := template.ParseFiles("./internal/templates/start.html")
	if err != nil {
		h.log.Error().Err(err).Msg("filed to show home page")
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	tmpl.Execute(w, nil)
}

// Все пользователи в БД
func (h *Handler) GetUsersList(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	users, err := h.service.GetUsersList(r.Context())
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		h.log.Error().Err(err).Msg("failed to get users list")
		return
	}

	data, err := json.Marshal(users)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		h.log.Error().Err(err).Msg("failed to marshal users list")
		return
	}

	w.Write(data)
}

func New(log zerolog.Logger, service service.Service) *Handler {
	return &Handler{
		log:     log,
		service: service,
	}
}
