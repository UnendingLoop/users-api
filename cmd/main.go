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
	userRepo := repository.NewGormUserRepository(db)
	userHandler := handler.UserHandler{Repo: userRepo}

	r := chi.NewRouter()

	r.Get("/users", userHandler.ListUsers)
	r.Get("/users/{id}", userHandler.GetUserByID)
	r.Post("/users", userHandler.CreateUser)
	r.Delete("/delete/{id}", userHandler.DeleteUser)
	r.Put("/update/{id}", userHandler.UpdateUser)

	fmt.Println("Server running on http://localhost:8080")
	if err := http.ListenAndServe(":8080", r); err != nil {
		log.Fatal(err)
	}
}
