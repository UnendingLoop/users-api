package main

import (
	"fmt"
	"log"
	"net/http"

	"github.com/UnendingLoop/users-api/cmd/internal/config"
	"github.com/UnendingLoop/users-api/cmd/internal/handler"
	"github.com/UnendingLoop/users-api/cmd/internal/repository"
	"github.com/go-chi/chi/v5"
)

func main() {
	db := config.ConnectDB("app.db")
	userRepo := repository.NewUserRepository(db)
	userHandler := handler.UserHandler{Repo: userRepo}

	r := chi.NewRouter()

	r.Get("/users", userHandler.ListUsers)
	r.Get("/users/{id}", userHandler.GetUserByID)
	r.Post("/users", userHandler.CreateUser)

	r.Delete("/delete/{id}", func(w http.ResponseWriter, r *http.Request) {})
	r.Put("/update/{id}", func(w http.ResponseWriter, r *http.Request) {})
	//3. 🛡 Валидация: защита от пустых полей и дублирующихся email

	fmt.Println("Server running on http://localhost:8080")
	if err := http.ListenAndServe(":8080", r); err != nil {
		log.Fatal(err)
	}
}
