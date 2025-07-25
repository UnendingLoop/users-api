package handler

import (
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"strconv"

	"github.com/UnendingLoop/users-api/cmd/internal/model"
	"github.com/UnendingLoop/users-api/cmd/internal/repository"
	"github.com/go-chi/chi/v5"
)

type UserHandler struct {
	Repo repository.UserRepository
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

func (UH *UserHandler) DeleteUser(w http.ResponseWriter, r *http.Request) {
	idstr := chi.URLParam(r, "id")
	id, err := strconv.ParseInt(idstr, 10, 64)
	if err != nil {
		http.Error(w, "Failed to parse user id", http.StatusInternalServerError) //HTTP 500 Int server error
		return
	}
	err = UH.Repo.DeleteUser(id)
	if err != nil {
		switch {
		case errors.Is(err, repository.ErrUserNotFound):
			http.Error(w, "User not found", http.StatusNotFound) //HTTP 404 Not found
		default:
			http.Error(w, fmt.Sprintf("Failed to delete user: %v", err), http.StatusBadRequest) //HTTP 400 Bad request
		}
		return
	}
	w.WriteHeader(http.StatusNoContent) //HTTP 204 No content
}
func (UH *UserHandler) UpdateUser(w http.ResponseWriter, r *http.Request) {
	var user model.User
	//распарсить id
	idstr := chi.URLParam(r, "id")
	id, err := strconv.ParseInt(idstr, 10, 64)
	if err != nil {
		http.Error(w, "Failed to parse user id", http.StatusInternalServerError)
		return
	}
	//распарсить тело запроса - достать данные и засунуть в структуру
	if err := json.NewDecoder(r.Body).Decode(&user); err != nil {
		http.Error(w, "Failed to decode user from json", http.StatusBadRequest)
		return
	}
	user.ID = id
	if err := UH.Repo.UpdateUser(&user); err != nil {
		switch {
		case errors.Is(err, repository.ErrEmailExists):
			http.Error(w, "Email already exists", http.StatusConflict) //HTTP 409 Conflict
		case errors.Is(err, repository.ErrEmptyfields):
			http.Error(w, "At least one field must not be empty", http.StatusBadRequest) //HTTP 400 Bad request
		case errors.Is(err, repository.ErrUserNotFound):
			http.Error(w, "User not found", http.StatusNotFound) //HTTP 404 Not found
		}
		return
	}

	//ответ
	w.WriteHeader(http.StatusOK) //HTTP 200 OK
	w.Header().Set("Content-Type", "application/json")
	if err := json.NewEncoder(w).Encode(user); err != nil {
		http.Error(w, "Failed to encode user", http.StatusInternalServerError)
		return
	}
}
