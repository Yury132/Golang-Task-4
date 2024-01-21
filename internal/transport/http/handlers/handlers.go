package handlers

import (
	"context"
	"fmt"
	"html/template"
	"net/http"
	"strconv"

	"github.com/Yury132/Golang-Task-4/internal/models"
	"github.com/Yury132/Golang-Task-4/internal/service"
	"github.com/gorilla/mux"
	"github.com/rs/zerolog"
)

type Service interface {
	// Получение всех пользователей БД
	GetUsersList(ctx context.Context) ([]models.User, error)
	// Добавление нового пользователя, если точно такой же уже не существует в БД
	HandleUser(ctx context.Context, name string, surname string, patronymic string) error
	// Удаление пользователя
	DeleteUser(ctx context.Context, id int) error
	// Получение конкретного пользователя по ID
	GetUser(ctx context.Context, id int) (models.User, error)
	// Обновление данных конкретного пользователя по ID
	EditUser(ctx context.Context, id int, getUserName string, getUserSurname string, getUserPatronymic string) error
}

type Handler struct {
	log     zerolog.Logger
	service Service
}

// Все пользователи в БД
func (h *Handler) GetUsersList(w http.ResponseWriter, r *http.Request) {
	//w.Header().Set("Content-Type", "application/json")

	users, err := h.service.GetUsersList(r.Context())
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		h.log.Error().Err(err).Msg("failed to get users list")
		return
	}

	// data, err := json.Marshal(users)
	// if err != nil {
	// 	w.WriteHeader(http.StatusInternalServerError)
	// 	h.log.Error().Err(err).Msg("failed to marshal users list")
	// 	return
	// }

	// w.Write(data)

	tmpl, err := template.ParseFiles("./internal/templates/start.html")
	if err != nil {
		h.log.Error().Err(err).Msg("failed to show start page")
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	tmpl.Execute(w, users)
}

// Удаление пользователя по ID
func (h *Handler) DeleteUser(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)

	// ID пользователя
	userId, err := strconv.Atoi(vars["userId"])
	if err != nil || userId < 0 {
		http.NotFound(w, r)
		http.Redirect(w, r, "/users-list", http.StatusSeeOther)
		return
	}

	h.log.Log().Msg(fmt.Sprintf("Delete user with ID=%v", userId))

	// Удаляем из таблицы в БД
	err = h.service.DeleteUser(r.Context(), userId)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		h.log.Error().Err(err).Msg("failed to delete user")
	}

	http.Redirect(w, r, "/users-list", http.StatusSeeOther)
}

// Добавление нового пользователя, если точно такой же уже не существует в БД
func (h *Handler) CreateUser(w http.ResponseWriter, r *http.Request) {

	// Имя пользователя из формы POST запрос
	getUserName := r.FormValue("userName")
	if getUserName == "" {
		h.log.Log().Msg("Имя пользователя пустое при добавлении")
		// Переадресуем пользователя на ту же страницу
		http.Redirect(w, r, "/users-list", http.StatusSeeOther)
		return
	}

	// Фамилия пользователя из формы POST запрос
	getUserSurname := r.FormValue("userSurname")
	if getUserSurname == "" {
		h.log.Log().Msg("Фамилия пользователя пустая при добавлении")
		// Переадресуем пользователя на ту же страницу
		http.Redirect(w, r, "/users-list", http.StatusSeeOther)
		return
	}

	// Отчество пользователя из формы POST запрос
	getUserPatronymic := r.FormValue("userPatronymic")
	if getUserPatronymic == "" {
		h.log.Log().Msg("Отчество пользователя пустое при добавлении")
		// Переадресуем пользователя на ту же страницу
		http.Redirect(w, r, "/users-list", http.StatusSeeOther)
		return
	}

	// Добавляем нового пользователя, проверяя при этом его существование в БД
	err := h.service.HandleUser(r.Context(), getUserName, getUserSurname, getUserPatronymic)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		h.log.Error().Err(err).Msg("failed to Create User")
	}

	// Переадресуем пользователя на ту же страницу
	http.Redirect(w, r, "/users-list", http.StatusSeeOther)
}

// Переход к конкретному пользователю по ID
func (h *Handler) GoUser(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)

	// ID пользователя
	userId, err := strconv.Atoi(vars["userId"])
	if err != nil || userId < 0 {
		http.NotFound(w, r)
		http.Redirect(w, r, "/users-list", http.StatusSeeOther)
		return
	}

	h.log.Log().Msg(fmt.Sprintf("Go to user with ID=%v", userId))

	// Получаем конкретного пользователя по ID
	user, err := h.service.GetUser(r.Context(), userId)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		h.log.Error().Err(err).Msg("failed to get user by ID")
		http.Redirect(w, r, "/users-list", http.StatusSeeOther)
		return
	}

	// Переходим на страницу
	tmpl, err := template.ParseFiles("./internal/templates/user.html")
	if err != nil {
		h.log.Error().Err(err).Msg("failed to show user page")
		w.WriteHeader(http.StatusInternalServerError)
		http.Redirect(w, r, "/users-list", http.StatusSeeOther)
		return
	}
	// Передаем данные
	tmpl.Execute(w, user)
}

// Обновление данных конкретного пользователя по ID
func (h *Handler) EditUser(w http.ResponseWriter, r *http.Request) {

	// ID пользователя
	userId, err := strconv.Atoi(r.FormValue("userID"))
	if err != nil || userId < 0 {
		http.NotFound(w, r)
		http.Redirect(w, r, "/users-list", http.StatusSeeOther)
		return
	}

	// Имя пользователя из формы POST запрос
	getUserName := r.FormValue("userName")
	if getUserName == "" {
		h.log.Log().Msg("Имя пользователя пустое при редактировании")
		// Переадресуем пользователя на ту же страницу
		http.Redirect(w, r, "/go-user/"+r.FormValue("userID"), http.StatusSeeOther)
		return
	}

	// Фамилия пользователя из формы POST запрос
	getUserSurname := r.FormValue("userSurname")
	if getUserSurname == "" {
		h.log.Log().Msg("Фамилия пользователя пустая при редактировании")
		// Переадресуем пользователя на ту же страницу
		http.Redirect(w, r, "/go-user/"+r.FormValue("userID"), http.StatusSeeOther)
		return
	}

	// Отчество пользователя из формы POST запрос
	getUserPatronymic := r.FormValue("userPatronymic")
	if getUserPatronymic == "" {
		h.log.Log().Msg("Отчество пользователя пустое при редактировании")
		// Переадресуем пользователя на ту же страницу
		http.Redirect(w, r, "/go-user/"+r.FormValue("userID"), http.StatusSeeOther)
		return
	}

	h.log.Log().Msg(fmt.Sprintf("Edit user with ID=%v", userId))

	// Обновляем данные
	err = h.service.EditUser(r.Context(), userId, getUserName, getUserSurname, getUserPatronymic)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		h.log.Error().Err(err).Msg("failed to Edit User")
	}

	http.Redirect(w, r, "/go-user/"+r.FormValue("userID"), http.StatusSeeOther)

}

func New(log zerolog.Logger, service service.Service) *Handler {
	return &Handler{
		log:     log,
		service: service,
	}
}
