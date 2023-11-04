package app

import (
	"log"
	"net/http"
	"time"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"

	"github.com/Niiazgulov/trueconf.git/internal/handlers"
	"github.com/Niiazgulov/trueconf.git/internal/storage"
)

// Запуск сервера
func Start() {
	repo, err := storage.NewFileStorage()
	if err != nil {
		log.Fatal(err)
	}

	r := chi.NewRouter()
	r.Use(middleware.RequestID)
	r.Use(middleware.RealIP)
	r.Use(middleware.Logger)
	r.Use(middleware.Recoverer)
	r.Use(middleware.Timeout(60 * time.Second))
	r.Get("/", func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte(time.Now().String()))
	})
	r.Route("/api/v1/users", func(r chi.Router) {
		r.Get("/", handlers.SearchUsersHandler(repo))
		r.Post("/", handlers.CreateUserHandler(repo))
		r.Get("/{id}", handlers.GetUserHandler(repo))
		r.Patch("/{id}", handlers.UpdateUserHandler(repo))
		r.Delete("/{id}", handlers.DeleteUserHandler(repo))
	})
	http.ListenAndServe(":3333", r)
}
