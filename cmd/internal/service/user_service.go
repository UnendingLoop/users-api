package service

import (
	"context"
	"errors"
	"fmt"

	"github.com/UnendingLoop/users-api/cmd/internal/model"
	"github.com/UnendingLoop/users-api/cmd/internal/repository"
	"gorm.io/gorm"
)

type UserServe struct {
	Repo repository.UserRepository
}

type UserService interface {
	CreateUser(user *model.User, ctx context.Context) error
	GetUserByID(id int64, ctx context.Context) (*model.User, error)
	ListUsers(ctx context.Context) ([]model.User, error)
	DeleteUser(id int64, ctx context.Context) error
	UpdateUser(user *model.User, ctx context.Context) error
}

func NewUserService(userRepo repository.UserRepository) UserServe {
	return UserServe{Repo: userRepo}
}

func (US *UserServe) CreateUser(user *model.User, ctx context.Context) error {
	//валидация входных данных - перенесено из хендлера
	if user.Name == "" || user.Surname == "" || user.Email == "" {
		return repository.ErrEmptySomeFields
	}
	if err := US.Repo.CheckUserExists(user.ID, ctx); err == nil {
		return repository.ErrUserExists
	}
	return US.CreateUser(user, ctx)
}
func (US *UserServe) GetUserByID(id int64, ctx context.Context) (*model.User, error) {
	if err := US.Repo.CheckUserExists(id, ctx); err != nil {
		return nil, err
	}
	user, err := US.GetUserByID(id, ctx)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, repository.ErrUserNotFound
		}
		return nil, err
	}
	return user, nil
}
func (US *UserServe) ListUsers(ctx context.Context) ([]model.User, error) {
	users, err := US.ListUsers(ctx)
	if err != nil {
		return nil, fmt.Errorf("Failed to fetch users list: %w", err)
	}
	return users, nil
}
func (US *UserServe) DeleteUser(id int64, ctx context.Context) error {
	count, err := US.Repo.DeleteUser(id, ctx)
	if count == 0 && err == nil {
		return fmt.Errorf("Failed to remove user: %v", repository.ErrUserNotFound)
	}
	return err
}
func (US *UserServe) UpdateUser(user *model.User, ctx context.Context) error {
	//проверка на ненулевой input
	if user.Email == "" && user.Name == "" && user.Surname == "" {
		return fmt.Errorf("Failed to update user info: %v", repository.ErrEmptyfields)
	}

	return US.UpdateUser(user, ctx)
}
