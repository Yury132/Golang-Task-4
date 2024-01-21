package handlers

import (
	"context"
	"fmt"
	"html/template"
	"net/http"
	"strconv"
	"strings"

	"github.com/Yury132/Golang-Task-4/internal/models"
	"github.com/Yury132/Golang-Task-4/internal/service"
	"github.com/gorilla/mux"
	"github.com/rs/zerolog"
)

type Service interface {
	// Получение всех пользователей БД
	GetUsersList(ctx context.Context) ([]models.User, error)
	// Получение определенных пользователей по возрасту
	GetUsersListAge(ctx context.Context, userAgeMin int, userAgeMax int) ([]models.User, error)
	// Получение определенных пользователей по полу
	GetUsersListGender(ctx context.Context, gender string) ([]models.User, error)
	// Получение определенных пользователей по национальности
	GetUsersListNation(ctx context.Context, userNation string) ([]models.User, error)
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

	h.log.Log().Msg("Получение всех пользователей")

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
		w.WriteHeader(http.StatusInternalServerError)
		h.log.Error().Err(err).Msg("failed to show start page")
		return
	}
	tmpl.Execute(w, users)
}

// Получение определенных пользователей по возрасту
func (h *Handler) GetUsersListAge(w http.ResponseWriter, r *http.Request) {

	h.log.Log().Msg("Получение определенных пользователей по возрасту")

	// Минимальный возраст из формы POST запрос
	userAgeMin, err := strconv.Atoi(r.FormValue("userAgeMin"))
	if err != nil || userAgeMin < 0 {
		h.log.Error().Err(err).Msg("failed to get min age")
		http.Redirect(w, r, "/users-list", http.StatusSeeOther)
		return
	}

	// Максимальный возраст из формы POST запрос
	userAgeMax, err := strconv.Atoi(r.FormValue("userAgeMax"))
	if err != nil || userAgeMax < 0 {
		h.log.Error().Err(err).Msg("failed to get max age")
		http.Redirect(w, r, "/users-list", http.StatusSeeOther)
		return
	}

	// Получаем пользователей
	users, err := h.service.GetUsersListAge(r.Context(), userAgeMin, userAgeMax)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		h.log.Error().Err(err).Msg("failed to get Users List Age")
		return
	}

	// Отображаем
	tmpl, err := template.ParseFiles("./internal/templates/start.html")
	if err != nil {
		h.log.Error().Err(err).Msg("failed to show start page")
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	tmpl.Execute(w, users)
}

// Получение определенных пользователей по полу
func (h *Handler) GetUsersListGender(w http.ResponseWriter, r *http.Request) {

	h.log.Log().Msg("Получение определенных пользователей по полу")

	vars := mux.Vars(r)
	// Мужчины или женщины
	gender, err := strconv.Atoi(vars["gender"])
	if err != nil {
		h.log.Error().Err(err).Msg("failed to get gender")
		http.Redirect(w, r, "/users-list", http.StatusSeeOther)
		return
	}

	// Определяем необходимый пол
	var getGender string
	if gender == 1 {
		getGender = "м"
	} else {
		getGender = "ж"
	}

	// Получаем пользователей
	users, err := h.service.GetUsersListGender(r.Context(), getGender)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		h.log.Error().Err(err).Msg("failed to get Users List Gender")
		return
	}

	// Отображаем
	tmpl, err := template.ParseFiles("./internal/templates/start.html")
	if err != nil {
		h.log.Error().Err(err).Msg("failed to show start page")
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	tmpl.Execute(w, users)
}

// Получение определенных пользователей по национальности
func (h *Handler) GetUsersListNation(w http.ResponseWriter, r *http.Request) {

	h.log.Log().Msg("Получение определенных пользователей по национальности")

	// Национальность из формы POST запрос
	userNation := r.FormValue("userNation")
	if userNation == "" {
		h.log.Log().Msg("Национальность пустая при фильтрации")
		// Переадресуем пользователя на ту же страницу
		http.Redirect(w, r, "/users-list", http.StatusSeeOther)
		return
	}

	// К верхнему регистру
	userNation = strings.ToUpper(userNation)

	// Получаем пользователей
	users, err := h.service.GetUsersListNation(r.Context(), userNation)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		h.log.Error().Err(err).Msg("failed to get Users List Nation")
		return
	}

	// Отображаем
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
		h.log.Error().Err(err).Msg("failed to get user ID to delete")
		http.Redirect(w, r, "/users-list", http.StatusSeeOther)
		return
	}

	h.log.Log().Msg(fmt.Sprintf("Удаление пользователя с ID=%v", userId))

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

	h.log.Log().Msg("Добавление нового пользователя")

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
		h.log.Error().Err(err).Msg("failed to go user by ID")
		http.Redirect(w, r, "/users-list", http.StatusSeeOther)
		return
	}

	h.log.Log().Msg(fmt.Sprintf("Переход к пользователю с ID=%v", userId))

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
		h.log.Error().Err(err).Msg("failed to get user to edit")
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
