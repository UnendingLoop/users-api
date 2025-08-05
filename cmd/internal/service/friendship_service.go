package service

import (
	"context"
	"errors"
	"fmt"

	"github.com/UnendingLoop/users-api/cmd/internal/model"
	"github.com/UnendingLoop/users-api/cmd/internal/repository"
)

// FriendServe
type FriendServe struct {
	Repo     repository.FriendRepository
	UserRepo repository.UserRepository
}

type FriendshipService interface {
	AddFriend(user, friend int64, ctx context.Context) error
	RemoveFriend(user, friend int64, ctx context.Context) error
	GetFriends(user int64, ctx context.Context) ([]model.User, error)
}

func NewFriendService(friendRepo repository.FriendRepository, userRepo repository.UserRepository) FriendServe {
	return FriendServe{Repo: friendRepo, UserRepo: userRepo}
}

func (FS *FriendServe) AddFriend(user, friend int64, ctx context.Context) error {
	if user == friend {
		return fmt.Errorf("Failed to make a friendship: %w", repository.ErrUserEqualsFriend)
	}
	//Даже если оба юзера существуют, для образования дружбы требуется 3 запроса в БД, но зато читаемый код и ошибки.
	// Можно использовать 1 запрос и обрабатывать ошибку БД - если хочется оптимизации.
	// нужно добавить проверку существования дружбы
	if user < 0 || friend < 0 {
		return fmt.Errorf("invalid ID format")
	}
	if err := FS.UserRepo.CheckIfExistsByID(user, ctx); errors.Is(err, repository.ErrUserNotFound) {
		return fmt.Errorf("Failed to make a friendship: %w", err)
	}
	if err := FS.UserRepo.CheckIfExistsByID(friend, ctx); errors.Is(err, repository.ErrUserNotFound) {
		return fmt.Errorf("Failed to make a friendship: %w", err)
	}
	friendship := &model.Friendship{
		RequesterID: user,
		AccepterID:  friend,
	}
	if err := FS.Repo.AddFriend(ctx, friendship); err != nil {
		return fmt.Errorf("Failed to make a friendship: %w", err)
	}
	return nil
}
func (FS *FriendServe) RemoveFriend(user, friend int64, ctx context.Context) error {
	if user == friend {
		return fmt.Errorf("Failed to remove a friend: %w", repository.ErrUserEqualsFriend)
	}
	if err := FS.UserRepo.CheckIfExistsByID(user, ctx); errors.Is(err, repository.ErrUserNotFound) {
		return fmt.Errorf("Failed to remove a friend: %w", err)
	}
	if err := FS.UserRepo.CheckIfExistsByID(friend, ctx); !errors.Is(err, repository.ErrUserNotFound) {
		return fmt.Errorf("Failed to remove a friend: %w", err)
	}

	friendship := &model.Friendship{
		RequesterID: user,
		AccepterID:  friend,
	}
	if err := FS.Repo.RemoveFriend(ctx, friendship); err != nil {
		return fmt.Errorf("Failed to remove a friend: %w", err)
	}
	return nil
}
func (FS *FriendServe) GetFriends(user int64, ctx context.Context) ([]model.User, error) {
	if err := FS.UserRepo.CheckIfExistsByID(user, ctx); errors.Is(err, repository.ErrUserNotFound) {
		return nil, fmt.Errorf("Failed to fetch list of friends: %w", err)
	}

	res, err := FS.Repo.GetFriends(ctx, user)
	if err != nil {
		return nil, fmt.Errorf("Failed to fetch list of friends: %w", err)
	}
	return res, err
}
