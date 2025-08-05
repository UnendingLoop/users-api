package main

import (
	"fmt"
	"log"
	"net/http"
	"os"

	"github.com/UnendingLoop/users-api/cmd/internal/config"
	"github.com/UnendingLoop/users-api/cmd/internal/handler"
	"github.com/UnendingLoop/users-api/cmd/internal/repository"
	"github.com/UnendingLoop/users-api/cmd/internal/service"
	_ "github.com/UnendingLoop/users-api/docs"
	"github.com/go-chi/chi/v5"
	"github.com/joho/godotenv"
	httpSwagger "github.com/swaggo/http-swagger"
)

// @title Users API
// @version 1.0
// @description REST API для управления пользователями и друзьями
// @host localhost:8080
// @BasePath /
func main() {
	if err := godotenv.Load(); err != nil {
		log.Println("Warning: .env file not found")
	}

	dsn := os.Getenv("DATABASE_URL")
	if dsn == "" {
		log.Fatal("DATABASE_URL is not set in env")
	}
	db := config.ConnectPostgres(dsn)

	userRepo := repository.NewGormUserRepository(db)
	userService := service.NewUserService(userRepo)
	userHandler := handler.UserHandler{Repo: userService}

	friendRepo := repository.NewGormFriendRepository(db)
	friendService := service.NewFriendService(friendRepo, userRepo)
	friendHandler := handler.FriendHandler{Repo: friendService}

	r := chi.NewRouter()

	r.Get("/users", userHandler.ListUsers)
	r.Get("/users/{id}", userHandler.GetUserByID)
	r.Post("/users", userHandler.CreateUser)
	r.Delete("/delete/{id}", userHandler.DeleteUser)
	r.Put("/update/{id}", userHandler.UpdateUser)

	r.Post("/users/{id1}/make_friend/{id2}", friendHandler.MakeFriend)
	r.Get("/users/{id}/friends", friendHandler.GetFriendsList)
	r.Delete("/users/{id1}/remove_friend/{id2}", friendHandler.RemoveFriend)

	r.Get("/swagger/*", httpSwagger.WrapHandler)

	fmt.Println("Server running on http://localhost:8080")
	if err := http.ListenAndServe(":8080", r); err != nil {
		log.Fatal(err)
	}
}
