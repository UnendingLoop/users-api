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

	CheckIfExistsByID(id int64, ctx context.Context) error
	CheckIfExistsByEmail(email string, ctx context.Context) error
}

type GormUserRepository struct {
	DB *gorm.DB
}

var ErrUserNotFound = errors.New("user not found")
var ErrUserExists = errors.New("user already exists")

var ErrEmailExists = errors.New("email already exists")
var ErrEmailNotFound = errors.New("email not found")

var ErrEmptyfields = errors.New("all fields are empty")
var ErrEmptySomeFields = errors.New("some fields are empty")

func NewGormUserRepository(db *gorm.DB) *GormUserRepository {
	return &GormUserRepository{DB: db}
}

func (r *GormUserRepository) CreateUser(user *model.User, ctx context.Context) error {
	return r.DB.WithContext(ctx).Create(&user).Error
}
func (r *GormUserRepository) GetUserByID(id int64, ctx context.Context) (*model.User, error) {
	var user model.User
	err := r.DB.WithContext(ctx).First(&user, id).Error
	return &user, err
}
func (r *GormUserRepository) ListUsers(ctx context.Context) ([]model.User, error) {
	var users []model.User
	err := r.DB.WithContext(ctx).Find(&users).Error
	return users, err
}
func (r *GormUserRepository) DeleteUser(id int64, ctx context.Context) (int64, error) {
	res := r.DB.WithContext(ctx).Delete(&model.User{}, id)
	return res.RowsAffected, res.Error
}
func (r *GormUserRepository) UpdateUser(user *model.User, ctx context.Context) error {
	return r.DB.WithContext(ctx).Save(&user).Error
}

func (r *GormUserRepository) CheckIfExistsByID(id int64, ctx context.Context) error {
	if id < 0 {
		return fmt.Errorf("invalid ID format")
	}
	var count int64
	err := r.DB.WithContext(ctx).Model(&model.User{}).Where("id = ?", id).Count(&count).Error
	switch {
	case count == 0 && err == nil:
		return ErrUserNotFound
	case count != 0 && err == nil:
		return ErrUserExists
	default:
		return err
	}
}
func (r *GormUserRepository) CheckIfExistsByEmail(email string, ctx context.Context) error {
	var count int64
	err := r.DB.WithContext(ctx).Model(&model.User{}).Where("email = ?", email).Count(&count).Error
	switch {
	case count == 0 && err == nil:
		return ErrEmailNotFound
	case count != 0 && err == nil:
		return ErrEmailExists
	default:
		return err
	}
}
