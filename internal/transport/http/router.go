package http

import (
	"net/http"

	"github.com/Yury132/Golang-Task-4/internal/transport/http/handlers"
	"github.com/gorilla/mux"
)

func InitRoutes(h *handlers.Handler) *mux.Router {
	r := mux.NewRouter()

	// Все пользователи в БД
	r.HandleFunc("/users-list", h.GetUsersList)
	//
	//r.HandleFunc("/users-list-age", h.GetUsersListAge)
	//
	//r.HandleFunc("/users-list-gender", h.GetUsersListGender)
	//
	//r.HandleFunc("/users-list-nation", h.GetUsersListNation)
	// Удаление пользователя по ID
	r.HandleFunc("/delete-user/{userId:[0-9]+}", h.DeleteUser).Methods(http.MethodGet)
	// Добавление нового пользователя, если точно такой же уже не существует в БД
	r.HandleFunc("/create-user", h.CreateUser).Methods(http.MethodPost)
	// Переход к конкретному пользователю по ID
	r.HandleFunc("/go-user/{userId:[0-9]+}", h.GoUser).Methods(http.MethodGet)
	// Обновление данных конкретного пользователя по ID
	r.HandleFunc("/edit-user", h.EditUser).Methods(http.MethodPost)

	http.Handle("/", r)

	return r
}
