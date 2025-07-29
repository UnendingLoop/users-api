package handler

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"

	"github.com/UnendingLoop/users-api/cmd/internal/service"
	"github.com/go-chi/chi/v5"
)

type FriendHandler struct {
	Repo service.FriendServe
}

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
		http.Error(w, fmt.Sprintf("Failed to make friendship: %v", err), http.StatusInternalServerError)
	}

	w.WriteHeader(http.StatusCreated)
}
func (FH *FriendHandler) RemoveFriend(w http.ResponseWriter, r *http.Request) {
	requester, err := strconv.ParseInt(chi.URLParam(r, "id1"), 10, 64)
	if err != nil {
		http.Error(w, "Failed to parse user id", http.StatusBadRequest)
		return
	}
	acceptor, err := strconv.ParseInt(chi.URLParam(r, "id2"), 10, 64)
	if err != nil {
		http.Error(w, "Invalid user id", http.StatusBadRequest)
		return
	}

	if err := FH.Repo.RemoveFriend(requester, acceptor, r.Context()); err != nil {
		http.Error(w, fmt.Sprintf("Failed to remove friend: %v", err), http.StatusBadRequest)
		return
	}

	w.WriteHeader(http.StatusNoContent)
}
func (FH *FriendHandler) GetFriendsList(w http.ResponseWriter, r *http.Request) {
	requester, err := strconv.ParseInt(chi.URLParam(r, "id"), 10, 64)
	if err != nil {
		http.Error(w, "Failed to parse user id", http.StatusInternalServerError) //HTTP 500 Int server error
		return
	}
	friends, err := FH.Repo.GetFriends(requester, r.Context())
	if err != nil {
		http.Error(w, fmt.Sprintf("Failed to get list of friends: %v", err), http.StatusBadRequest)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	if err := json.NewEncoder(w).Encode(friends); err != nil {
		http.Error(w, "Failed to encode users", http.StatusInternalServerError)
		return
	}
}
