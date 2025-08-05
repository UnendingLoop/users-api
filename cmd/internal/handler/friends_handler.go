package handler

import (
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"strconv"

	"github.com/UnendingLoop/users-api/cmd/internal/repository"
	"github.com/UnendingLoop/users-api/cmd/internal/service"
	"github.com/go-chi/chi/v5"
)

// FriendHandler handles HTTP requests related to user friendships.
type FriendHandler struct {
	Repo service.FriendServe
}

// MakeFriend - хендлер для создания связи между 2мя существующими в базе пользователями
// @Summary      Хендлер для создания новой связи - дружбы
// @Description  Создаёт новую связь между 2мя существующими пользователями
// @Tags         friendship
// @Produce      plain
// @Param        id1	path	int	true	"User id 1 - friendship requester"
// @Param        id2  	path	int	true	"User id 2 - friendship acceptor"
// @Success      201   {string}  string  "Successfully added a new friend"
// @Failure      400   {string}  string  "Invalid data"
// @Failure      404   {string}  string  "One or both users don't exist"
// @Failure      500   {string}  string  "Internal server error"
// @Router       /users/{id1}/make_friend/{id2} [post]
func (FH *FriendHandler) MakeFriend(w http.ResponseWriter, r *http.Request) {
	requester, err := strconv.ParseInt(chi.URLParam(r, "id1"), 10, 64)
	if err != nil {
		http.Error(w, "Invalid user id", http.StatusBadRequest)
		return
	}
	acceptor, err := strconv.ParseInt(chi.URLParam(r, "id2"), 10, 64)
	if err != nil {
		http.Error(w, "Invalid user id", http.StatusBadRequest)
		return
	}

	if err := FH.Repo.AddFriend(requester, acceptor, r.Context()); err != nil {
		switch {
		case errors.Is(err, repository.ErrUserNotFound):
			http.Error(w, fmt.Sprintf("Failed to make friendship: %v", err), http.StatusNotFound)
			return
		default:
			http.Error(w, fmt.Sprintf("Failed to make friendship: %v", err), http.StatusInternalServerError)
			return
		}
	}
	w.WriteHeader(http.StatusCreated)
}

// RemoveFriend - хендлер для удаления связи между 2мя пользователями
// @Summary      Удаление существующей связи - дружбы
// @Description  Удаляет существующую связь между 2мя пользователями, id обоих берутся из URL
// @Tags         friendship
// @Produce      plain
// @Param        id1	path	int	true	"User id 1 - friendship requester"
// @Param        id2  	path	int	true	"User id 2 - friendship acceptor"
// @Success      204   {object}  model.User
// @Failure      400   {string}  string  "Invalid data input"
// @Failure      404   {string}  string  "One or both users don't exist"
// @Failure      500   {string}  string  "Internal server error"
// @Router       /users/{id1}/remove_friend/{id2} [delete]
func (FH *FriendHandler) RemoveFriend(w http.ResponseWriter, r *http.Request) {
	requester, err := strconv.ParseInt(chi.URLParam(r, "id1"), 10, 64)
	if err != nil {
		http.Error(w, "Invalid user id", http.StatusBadRequest)
		return
	}
	acceptor, err := strconv.ParseInt(chi.URLParam(r, "id2"), 10, 64)
	if err != nil {
		http.Error(w, "Invalid user id", http.StatusBadRequest)
		return
	}

	if err := FH.Repo.RemoveFriend(requester, acceptor, r.Context()); err != nil {
		switch {
		case errors.Is(err, repository.ErrUserNotFound):
			http.Error(w, fmt.Sprintf("Failed to remove friend: %v", err), http.StatusNotFound)
			return
		case errors.Is(err, repository.ErrUserEqualsFriend):
			http.Error(w, fmt.Sprintf("Failed to remove friend: %v", err), http.StatusBadRequest)
			return
		default:
			http.Error(w, fmt.Sprintf("Failed to remove friend: %v", err), http.StatusInternalServerError)
			return
		}
	}

	w.WriteHeader(http.StatusNoContent)
}

// GetFriendsList - хендлер для получения списка друзей пользователя
// @Summary      Получение списка друзей пользователя
// @Description  Возвращает массив JSON из пользователей, которые состоят в связи с указанным в запросе пользователем
// @Tags         friendship
// @Produce      json
// @Param        id		path	int	true	"User id - friendship requester"
// @Success      200   {array}  model.User  "Successful load of friends list"
// @Failure      400   {string}  string  "Invalid data"
// @Failure      404   {string}  string  "User not found"
// @Failure      500   {string}  string  "Internal server error"
// @Router       /users/{id}/friends [get]
func (FH *FriendHandler) GetFriendsList(w http.ResponseWriter, r *http.Request) {
	requester, err := strconv.ParseInt(chi.URLParam(r, "id"), 10, 64)
	if err != nil {
		http.Error(w, "Failed to parse user id", http.StatusBadRequest)
		return
	}
	friends, err := FH.Repo.GetFriends(requester, r.Context())
	if err != nil {
		switch {
		case errors.Is(err, repository.ErrUserNotFound):
			http.Error(w, fmt.Sprintf("Failed to get list of friends: %v", err), http.StatusNotFound)
			return
		default:
			http.Error(w, fmt.Sprintf("Failed to get list of friends: %v", err), http.StatusInternalServerError)
			return
		}
	}

	w.Header().Set("Content-Type", "application/json")
	if err := json.NewEncoder(w).Encode(friends); err != nil {
		http.Error(w, "Failed to encode users", http.StatusInternalServerError)
		return
	}
	w.WriteHeader(http.StatusOK)
}
