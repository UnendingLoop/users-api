package handler

import (
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"strconv"

	"github.com/UnendingLoop/users-api/cmd/internal/model"
	"github.com/UnendingLoop/users-api/cmd/internal/repository"
	"github.com/UnendingLoop/users-api/cmd/internal/service"
	"github.com/go-chi/chi/v5"
)

// UserHandler handles HTTP-requests related to user management.
type UserHandler struct {
	Repo service.UserServe
}

// CreateUser - хендлер для создания нового пользователя в базе
// @Summary      Хендлер для создания нового пользователя
// @Description  Создаёт нового пользователя из данных в теле запроса
// @Tags         users
// @Accept       json
// @Produce      json
// @Param        user  body      model.User  true  "User info"
// @Success      201   {object}  model.User
// @Failure      400   {string}  string  "Incomplete data input"
// @Failure      500   {string}  string  "Internal server error"
// @Failure      409   {string}  string  "Email conflict: already in use"
// @Router       /users [post]
func (UH *UserHandler) CreateUser(w http.ResponseWriter, r *http.Request) {
	var newUser model.User

	if err := json.NewDecoder(r.Body).Decode(&newUser); err != nil {
		http.Error(w, "Invalid JSON", http.StatusBadRequest)
		return
	}

	if err := UH.Repo.CreateUser(&newUser, r.Context()); err != nil {
		switch {
		case errors.Is(err, repository.ErrEmptySomeFields):
			http.Error(w, fmt.Sprintf("Validation error: %v", err), http.StatusBadRequest)
			return
		case errors.Is(err, repository.ErrUserExists):
			http.Error(w, fmt.Sprintf("Conflict: %v", err), http.StatusConflict)
			return
		default:
			http.Error(w, fmt.Sprintf("Internal error: %v", err), http.StatusInternalServerError)
			return
		}
	}

	w.Header().Set("Content-Type", "application/json")
	if err := json.NewEncoder(w).Encode(newUser); err != nil {
		http.Error(w, "Failed do encode user", http.StatusInternalServerError)
		return
	}
	w.WriteHeader(http.StatusCreated)
}

// ListUsers - хендлер для получения списка всех юзеров из базы
// @Summary      Хендлер для получения списка всех юзеров из базы
// @Description  Отдает массив из всех пользователей базы
// @Tags         users
// @Produce      json
// @Success      200   {array}  model.User
// @Failure      500   {string}  string  "Internal server error"
// @Router       /users [get]
func (UH *UserHandler) ListUsers(w http.ResponseWriter, r *http.Request) {
	users, err := UH.Repo.ListUsers(r.Context())
	if err != nil {
		http.Error(w, "Failed to fetch users", http.StatusInternalServerError)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	if err := json.NewEncoder(w).Encode(users); err != nil {
		http.Error(w, "Failed to encode users", http.StatusInternalServerError)
	}
	w.WriteHeader(http.StatusOK)
}

// GetUserByID - хендлер для получения пользователя по ID
// @Summary      Получение пользователя по ID
// @Description  Возвращает пользователя в формате JSON по ID из URL
// @Tags         users
// @Produce      json
// @Param        id   path      int  true  "ID пользователя"
// @Success      200  {object}  model.User
// @Failure      404  {string}  string  "User not found"
// @Failure      500  {string}  string  "Internal server error"
// @Router       /users/{id} [get]
func (UH *UserHandler) GetUserByID(w http.ResponseWriter, r *http.Request) {
	idstr := chi.URLParam(r, "id")
	id, err := strconv.ParseInt(idstr, 10, 64)
	if err != nil {
		http.Error(w, "Failed to parse user id", http.StatusInternalServerError)
		return
	}
	user, err := UH.Repo.GetUserByID(id, r.Context())
	if err != nil {
		http.Error(w, "Failed to find user ID", http.StatusNotFound)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	if err := json.NewEncoder(w).Encode(user); err != nil {
		http.Error(w, "Failed to encode user", http.StatusInternalServerError)
		return
	}
	w.WriteHeader(http.StatusOK)
}

// DeleteUser - хендлер для удаления пользователя по ID
// @Summary      Удаление пользователя по ID
// @Description  Удаляет пользователя по ID из URL
// @Tags         users
// @Param        id   path      int  true  "ID пользователя"
// @Success      204  {string}  string  "No Content"
// @Failure      404  {string}  string  "User not found"
// @Failure      400  {string}  string  "Bad request"
// @Failure      500  {string}  string  "Internal server error"
// @Router       /delete/{id}	[delete]
func (UH *UserHandler) DeleteUser(w http.ResponseWriter, r *http.Request) {
	idstr := chi.URLParam(r, "id")
	id, err := strconv.ParseInt(idstr, 10, 64)
	if err != nil {
		http.Error(w, "Failed to parse user id", http.StatusInternalServerError)
		return
	}
	err = UH.Repo.DeleteUser(id, r.Context())
	if err != nil {
		switch {
		case errors.Is(err, repository.ErrUserNotFound):
			http.Error(w, "User not found", http.StatusNotFound)
		default:
			http.Error(w, fmt.Sprintf("Failed to delete user: %v", err), http.StatusBadRequest)
		}
		return
	}
	w.WriteHeader(http.StatusNoContent) //HTTP 204 No content
}

// UpdateUser - хендлер для обновления данных пользователя по ID
// @Summary      Обновление пользователя по ID
// @Description  Обновляет пользователя по ID из URL, новые данные берутся из тела запроса
// @Tags         users
// @Accept       json
// @Produce      json
// @Param        id   path      int  true  "ID пользователя"
// @Success      200  {object}  model.User	"User updated successfully"
// @Failure      404  {string}  string  "User not found"
// @Failure      400  {string}  string  "Bad request"
// @Failure      409  {string}  string  "Conflict: new email already in use"
// @Failure      500  {string}  string  "Internal server error"
// @Router       /update/{id}	[put]
func (UH *UserHandler) UpdateUser(w http.ResponseWriter, r *http.Request) {
	var user model.User
	//распарсить id
	idStr := chi.URLParam(r, "id")
	id, err := strconv.ParseInt(idStr, 10, 64)
	if err != nil {
		http.Error(w, "Failed to parse user id", http.StatusInternalServerError)
		return
	}
	//распарсить тело запроса - достать данные и засунуть в структуру
	if err := json.NewDecoder(r.Body).Decode(&user); err != nil {
		http.Error(w, "Failed to decode user from json", http.StatusBadRequest) //HTTP 400 Bad request
		return
	}
	user.ID = id
	if err := UH.Repo.UpdateUser(&user, r.Context()); err != nil {
		switch {
		case errors.Is(err, repository.ErrEmailExists):
			http.Error(w, "Email already exists", http.StatusConflict) //HTTP 409 Conflict
		case errors.Is(err, repository.ErrEmptyFields):
			http.Error(w, "At least one field must not be empty", http.StatusBadRequest) //HTTP 400 Bad request
		case errors.Is(err, repository.ErrUserNotFound):
			http.Error(w, "User not found", http.StatusNotFound) //HTTP 404 Not found
		}
		return
	}

	//ответ
	w.Header().Set("Content-Type", "application/json")
	if err := json.NewEncoder(w).Encode(user); err != nil {
		http.Error(w, "Failed to encode user", http.StatusInternalServerError)
		return
	}
	w.WriteHeader(http.StatusOK) //HTTP 200 OK
}
