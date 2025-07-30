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
		return fmt.Errorf("Failed to create a new user: %w", repository.ErrEmptySomeFields)
	}
	if err := US.Repo.CheckIfExistsByEmail(user.Email, ctx); !errors.Is(err, repository.ErrEmailNotFound) {
		return fmt.Errorf("Failed to create a new user: %w", err)
	}
	return US.Repo.CreateUser(user, ctx)
}
func (US *UserServe) GetUserByID(id int64, ctx context.Context) (*model.User, error) {
	if id < 0 {
		return nil, fmt.Errorf("Failed to get user info: invalid ID format")
	}
	if err := US.Repo.CheckIfExistsByID(id, ctx); !errors.Is(err, repository.ErrUserExists) {
		return nil, fmt.Errorf("Failed to get user info: %w", err)
	}

	user, err := US.Repo.GetUserByID(id, ctx)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, fmt.Errorf("Failed to get user info: %w", repository.ErrUserNotFound)
		}
		return nil, fmt.Errorf("Failed to get user info: %w", err)
	}
	return user, nil
}
func (US *UserServe) ListUsers(ctx context.Context) ([]model.User, error) {
	users, err := US.Repo.ListUsers(ctx)
	if err != nil {
		return nil, fmt.Errorf("Failed to fetch users list: %w", err)
	}
	return users, nil
}
func (US *UserServe) DeleteUser(id int64, ctx context.Context) error {
	count, err := US.Repo.DeleteUser(id, ctx)
	if count == 0 && err == nil {
		return fmt.Errorf("Failed to remove user: %w", repository.ErrUserNotFound)
	}
	return err
}
func (US *UserServe) UpdateUser(user *model.User, ctx context.Context) error {
	//проверка на ненулевой input
	if user.Email == "" && user.Name == "" && user.Surname == "" {
		return fmt.Errorf("Failed to update user info: %w", repository.ErrEmptyfields)
	}
	//загрузить из базы юзера с этим id - сразу проверить существует ли такой юзер
	dbUser, err := US.Repo.GetUserByID(user.ID, ctx)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return fmt.Errorf("Failed to update user info: %w", repository.ErrUserNotFound)
		}
		return err
	}

	//скопировать ненулевые поля из user в dbUser,в будущем можно добавить нормализацию регистра
	if user.Name != "" {
		dbUser.Name = user.Name
	}
	if user.Surname != "" {
		dbUser.Surname = user.Surname
	}
	if user.Email != "" {
		if err := US.Repo.CheckIfExistsByEmail(user.Email, ctx); errors.Is(err, repository.ErrEmailNotFound) {
			dbUser.Email = user.Email
		} else {
			return fmt.Errorf("Failed to update user info: %w", err)
		}
	}
	if err := US.Repo.UpdateUser(dbUser, ctx); err != nil {
		return fmt.Errorf("Failed to update user info: %w", err)
	}
	return nil
}
