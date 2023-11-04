package handlers

import (
	"encoding/json"
	"log"
	"net/http"
	"os"
	"strconv"
	"time"

	"github.com/Niiazgulov/trueconf.git/internal/storage"
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/render"
)

// SearchUsersHandler - обработчик эндпоинта GET "/api/v1/users/" - просмотр всех записей о пользователях.
func SearchUsersHandler(repo storage.UsersFile) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		render.JSON(w, r, repo.UsersStorage.List)
	}
}

// CreateUserHandler - обработчик эндпоинта POST "/api/v1/users/" - добавление в файл нового пользователя.
func CreateUserHandler(repo storage.UsersFile) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		request := UserRequest{}
		if err := render.Bind(r, &request); err != nil {
			_ = render.Render(w, r, ErrInvalidRequest(err))
			log.Printf("CreateUserHandler: unable to render: %v", err)
			http.Error(w, "CreateUserHandler: unable to render", http.StatusInternalServerError)
			return
		}

		repo.UsersStorage.Increment++
		user := storage.User{
			CreatedAt:   time.Now(),
			DisplayName: request.DisplayName,
			Email:       request.Email,
		}

		id := strconv.Itoa(repo.UsersStorage.Increment)
		repo.UsersStorage.List[id] = user

		b, err := json.Marshal(&repo.UsersStorage)
		if err != nil {
			log.Printf("CreateUserHandler: unable to marshal: %v", err)
			http.Error(w, "CreateUserHandler: unable to marshal", http.StatusInternalServerError)
			return
		}
		if err := os.Truncate(storage.Store, 0); err != nil {
			log.Printf("CreateUserHandler: unable to Truncate file: %v", err)
			http.Error(w, "CreateUserHandler: unable to Truncate file", http.StatusInternalServerError)
			return
		}
		repo.FileJSON.Write(b)

		render.Status(r, http.StatusCreated)
		render.JSON(w, r, map[string]interface{}{
			"user_id": id,
		})
	}
}

// GetUserHandler - обработчик эндпоинта GET "/api/v1/users/{id}" - просмотр записи конкретного пользователя (по id).
func GetUserHandler(repo storage.UsersFile) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		id := chi.URLParam(r, "id")
		render.JSON(w, r, repo.UsersStorage.List[id])
	}
}

// UpdateUserHandler - обработчик эндпоинта PATCH "/api/v1/users/{id}" - изменение информации конкретного пользователя (по id).
func UpdateUserHandler(repo storage.UsersFile) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		request := UserRequest{}

		if err := render.Bind(r, &request); err != nil {
			_ = render.Render(w, r, ErrInvalidRequest(err))
			log.Printf("UpdateUserHandler: unable to render: %v", err)
			return
		}

		id := chi.URLParam(r, "id")

		if _, ok := repo.UsersStorage.List[id]; !ok {
			_ = render.Render(w, r, ErrInvalidRequest(ErrUserNotFound))
			log.Println("UpdateUserHandler: unable to render")
		}

		u := repo.UsersStorage.List[id]
		u.DisplayName = request.DisplayName
		u.Email = request.Email
		repo.UsersStorage.List[id] = u

		b, err := json.Marshal(&repo.UsersStorage)
		if err != nil {
			log.Printf("UpdateUserHandler: unable to marshal internal file storage map: %v", err)
			http.Error(w, "UpdateUserHandler: unable to marshal", http.StatusInternalServerError)
			return
		}
		if err := os.Truncate(storage.Store, 0); err != nil {
			log.Printf("UpdateUserHandler: unable to Truncate file: %v", err)
			http.Error(w, "UpdateUserHandler: unable to Truncate file", http.StatusInternalServerError)
			return
		}
		repo.FileJSON.Write(b)
		render.Status(r, http.StatusNoContent)
	}
}

// DeleteUserHandler - обработчик эндпоинта DELETE "/api/v1/users/{id}" - удаление конкретного пользователя (по id).
func DeleteUserHandler(repo storage.UsersFile) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		id := chi.URLParam(r, "id")

		if _, ok := repo.UsersStorage.List[id]; !ok {
			_ = render.Render(w, r, ErrInvalidRequest(ErrUserNotFound))
			log.Fatal("DeleteUser: unable to render")
		}

		delete(repo.UsersStorage.List, id)

		b, err := json.Marshal(&repo.UsersStorage)
		if err != nil {
			log.Printf("DeleteUserHandler: unable to marshal internal file storage map: %v", err)
			http.Error(w, "DeleteUserHandler: unable to marshal", http.StatusInternalServerError)
			return
		}
		if err := os.Truncate(storage.Store, 0); err != nil {
			log.Printf("DeleteUserHandler: unable to Truncate file: %v", err)
			http.Error(w, "DeleteUserHandler: unable to Truncate file", http.StatusInternalServerError)
			return
		}
		repo.FileJSON.Write(b)
		render.Status(r, http.StatusNoContent)

	}
}
