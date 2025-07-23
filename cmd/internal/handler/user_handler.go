package handler

import (
	"encoding/json"
	"net/http"
	"strconv"

	"github.com/UnendingLoop/users-api/cmd/internal/model"
	"github.com/UnendingLoop/users-api/cmd/internal/repository"
	"github.com/go-chi/chi/v5"
)

type UserHandler struct {
	Repo *repository.UserRepository
}

func (UH *UserHandler) ListUsers(w http.ResponseWriter, r *http.Request) {
	users, err := UH.Repo.ListUsers()
	if err != nil {
		http.Error(w, "Failed to fetch users", http.StatusInternalServerError)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	if err := json.NewEncoder(w).Encode(users); err != nil {
		http.Error(w, "Failed to encode users", http.StatusInternalServerError)
	}
}

func (UH *UserHandler) GetUserByID(w http.ResponseWriter, r *http.Request) {
	idstr := chi.URLParam(r, "id")
	id, err := strconv.ParseInt(idstr, 10, 64)
	if err != nil {
		http.Error(w, "Failed to parse user id", http.StatusInternalServerError)
		return
	}

	user, err := UH.Repo.GetUserByID(id)
	if err != nil {
		http.Error(w, "Failed to find user ID", http.StatusNotFound)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	if err := json.NewEncoder(w).Encode(user); err != nil {
		http.Error(w, "Failed to encode user", http.StatusInternalServerError)
		return
	}
}

func (UH *UserHandler) CreateUser(w http.ResponseWriter, r *http.Request) {
	var newUser model.User

	if err := json.NewDecoder(r.Body).Decode(&newUser); err != nil {
		http.Error(w, "Invalid JSON", http.StatusBadRequest)
		return
	}

	if newUser.Name == "" || newUser.Surname == "" || newUser.Email == "" {
		http.Error(w, "Missing name or email", http.StatusBadRequest)
		return
	}

	if err := UH.Repo.CreateUser(&newUser); err != nil {
		http.Error(w, "Failed to create user", http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusCreated)
	w.Header().Set("Content-Type", "application/json")
	if err := json.NewEncoder(w).Encode(newUser); err != nil {
		http.Error(w, "Failed do encode user", http.StatusInternalServerError)
		return
	}
}
