package repository

import (
	"context"
	"errors"
	"fmt"

	"github.com/UnendingLoop/users-api/cmd/internal/model"
	"gorm.io/gorm"
)

type UserRepository interface {
	CreateUser(user *model.User, ctx context.Context) error
	GetUserByID(id int64, ctx context.Context) (*model.User, error)
	ListUsers(ctx context.Context) ([]model.User, error)
	DeleteUser(id int64, ctx context.Context) (int64, error)
	UpdateUser(user *model.User, ctx context.Context) error
	CheckUserExists(id int64, ctx context.Context) error
}

type GormUserRepository struct {
	DB *gorm.DB
}

var ErrUserNotFound = errors.New("user not found")
var ErrEmailExists = errors.New("email already exists")
var ErrUserExists = errors.New("user already exists")
var ErrEmptyfields = errors.New("all fields are empty")
var ErrEmptySomeFields = errors.New("some fields are empty")

func NewGormUserRepository(db *gorm.DB) *GormUserRepository {
	return &GormUserRepository{DB: db}
}

// User management methods:
func (r *GormUserRepository) CreateUser(user *model.User, ctx context.Context) error {
	return r.DB.WithContext(ctx).Create(&user).Error
}
func (r *GormUserRepository) GetUserByID(id int64, ctx context.Context) (*model.User, error) {
	var user model.User
	err := r.DB.First(&user, id).Error
	return &user, err
}
func (r *GormUserRepository) ListUsers(ctx context.Context) ([]model.User, error) {
	var users []model.User
	err := r.DB.Find(&users).Error
	return users, err
}
func (r *GormUserRepository) DeleteUser(id int64, ctx context.Context) (int64, error) {
	res := r.DB.Delete(&model.User{}, id)
	return res.RowsAffected, res.Error
}
func (r *GormUserRepository) UpdateUser(user *model.User, ctx context.Context) error {
	//загрузить из базы юзера с этим id - сразу проверить существует ли такой юзер
	var dbUser model.User
	err := r.DB.First(&dbUser, user.ID).Error
	if err != nil {
		return ErrUserNotFound
	}

	//проверить наличие подзаменного имейл в базе
	if user.Email != "" && user.Email != dbUser.Email {
		var tmp model.User
		if err := r.DB.Where("email = ? AND id != ?", user.Email, dbUser.ID).First(&tmp).Error; err == nil {
			return ErrEmailExists
		}
	}

	//скопировать ненулевые поля из user в dbUser
	if user.Name != "" {
		dbUser.Name = user.Name
	}
	if user.Surname != "" {
		dbUser.Surname = user.Surname
	}
	if user.Email != "" {
		dbUser.Email = user.Email
	}

	return r.DB.Save(&dbUser).Error
}

func (r *GormUserRepository) CheckUserExists(id int64, ctx context.Context) error {
	if id < 0 {
		return fmt.Errorf("invalid ID format")
	}
	if _, err := r.GetUserByID(id, ctx); err != nil {
		if errors.Is(err, ErrUserNotFound) {
			return fmt.Errorf("user %d not found: %w", id, err)
		}
		return fmt.Errorf("failed to fetch user: %w", err)
	}
	return nil //если существует, возвращаем nil
}
