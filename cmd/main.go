package main

import (
	"fmt"
	"log"
	"net/http"

	"github.com/UnendingLoop/users-api/cmd/internal/config"
	"github.com/UnendingLoop/users-api/cmd/internal/handler"
	"github.com/UnendingLoop/users-api/cmd/internal/repository"
	"github.com/UnendingLoop/users-api/cmd/internal/service"
	"github.com/go-chi/chi/v5"
)

func main() {
	//SQLite
	//db := config.ConnectDB("app.db")

	//PostgresQL
	dsn := "host=localhost user=fanil password=0123 dbname=sandbox port=5432 sslmode=disable"
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

	fmt.Println("Server running on http://localhost:8080")
	if err := http.ListenAndServe(":8080", r); err != nil {
		log.Fatal(err)
	}
}
