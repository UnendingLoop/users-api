package service

import (
	"context"
	"errors"
	"fmt"

	"github.com/UnendingLoop/users-api/cmd/internal/model"
	"github.com/UnendingLoop/users-api/cmd/internal/repository"
)

type FriendServe struct {
	Repo     repository.FriendRepository
	UserRepo repository.UserRepository
}

type FriendshipService interface {
	AddFriend(user, friend int64, ctx context.Context) error
	RemoveFriend(user, friend int64, ctx context.Context) error
	GetFriends(user int64, ctx context.Context) ([]model.User, error)
}

var ErrUserEqualsFriend = errors.New("user cannot be friend to himself")

func NewFriendService(friendRepo repository.FriendRepository, userRepo repository.UserRepository) FriendServe {
	return FriendServe{Repo: friendRepo, UserRepo: userRepo}
}

func (FS *FriendServe) AddFriend(user, friend int64, ctx context.Context) error {
	if user == friend {
		return ErrUserEqualsFriend
	}
	//Даже если оба юзера существуют, для обраховани дружбы требуется 3 запроса в БД, но зато читаемый код и ошибки.
	// Можно использовать 1 запрос и обрабатывать ошибку БД - если хочется оптимизации.
	if user < 0 || friend < 0 {
		return fmt.Errorf("invalid ID format")
	}
	if err := FS.UserRepo.CheckUserExists(user, ctx); err != nil {
		return err
	}
	if err := FS.UserRepo.CheckUserExists(friend, ctx); err != nil {
		return err
	}

	if err := FS.Repo.AddFriend(user, friend, ctx); err != nil {
		return fmt.Errorf("failed to add a friend: %w", err)
	}
	return nil
}
func (FS *FriendServe) RemoveFriend(user, friend int64, ctx context.Context) error {
	if user == friend {
		return ErrUserEqualsFriend
	}
	if err := FS.UserRepo.CheckUserExists(user, ctx); err != nil {
		return err
	}
	if err := FS.UserRepo.CheckUserExists(friend, ctx); err != nil {
		return err
	}

	if err := FS.Repo.RemoveFriend(user, friend, ctx); err != nil {
		return fmt.Errorf("failed to remove a friend: %w", err)
	}
	return nil
}
func (FS *FriendServe) GetFriends(user int64, ctx context.Context) ([]model.User, error) {
	if err := FS.UserRepo.CheckUserExists(user, ctx); err != nil {
		return nil, err
	}
	return FS.Repo.GetFriends(user, ctx)
}
