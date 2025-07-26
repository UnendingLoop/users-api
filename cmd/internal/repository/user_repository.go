package repository

import (
	"context"
	"errors"

	"github.com/UnendingLoop/users-api/cmd/internal/model"
	"gorm.io/gorm"
)

// UserRepository - интерфейс для мокирования в тестах
type UserRepository interface {
	CreateUser(user *model.User, ctx context.Context) error
	GetUserByID(id int64, ctx context.Context) (*model.User, error)
	ListUsers(ctx context.Context) ([]model.User, error)
	DeleteUser(id int64, ctx context.Context) error
	UpdateUser(user *model.User, ctx context.Context) error
}

type FriendshipRepository interface {
	AddFriend(user, friend int64, ctx context.Context) error
	RemoveFriend(user, friend int64, ctx context.Context) error
	GetFriends(user int64, ctx context.Context) ([]model.User, error)
}

type GormUserRepository struct {
	DB *gorm.DB
}

var ErrUserNotFound = errors.New("user not found")
var ErrEmailExists = errors.New("email already exists")
var ErrEmptyfields = errors.New("all fields are empty")

func NewGormUserRepository(db *gorm.DB) *GormUserRepository {
	return &GormUserRepository{DB: db}
}

// User management wethods:
func (r *GormUserRepository) CreateUser(user *model.User, ctx context.Context) error {
	return r.DB.WithContext(ctx).Create(&user).Error
}

func (r *GormUserRepository) GetUserByID(id int64, ctx context.Context) (*model.User, error) {
	var user model.User
	err := r.DB.First(&user, id).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, ErrUserNotFound
		}
		return nil, err
	}

	return &user, nil
}

func (r *GormUserRepository) ListUsers(ctx context.Context) ([]model.User, error) {
	var users []model.User
	err := r.DB.Find(&users).Error
	return users, err
}

func (r *GormUserRepository) DeleteUser(id int64, ctx context.Context) error {
	res := r.DB.Delete(&model.User{}, id)
	if res.RowsAffected == 0 {
		return ErrUserNotFound
	}
	return res.Error
}

func (r *GormUserRepository) UpdateUser(user *model.User, ctx context.Context) error {
	//проверка на ненулевой input
	if user.Email == "" && user.Name == "" && user.Surname == "" {
		return ErrEmptyfields
	}

	//загрузить из базы юзера с этим id - сразу проверить существует ли такой юзер
	var dbUser model.User
	err := r.DB.First(&dbUser, user.ID).Error
	if err != nil {
		return ErrUserNotFound
	}

	//проверить наличие подзаменного имейл в базе
	if user.Email != "" && user.Email != dbUser.Email {
		var tmp model.User
		if err := r.DB.Where("email = ? AND id != ?", user.Email, dbUser.ID).First(&tmp).Error; err != nil {
			return ErrEmailExists
		}
	}

	//скопировать ненулевые поля из user в dbUser
	if user.Name == "" {
		dbUser.Name = user.Name
	}
	if user.Surname == "" {
		dbUser.Surname = user.Surname
	}
	if user.Email == "" {
		dbUser.Email = user.Email
	}

	return r.DB.Save(&dbUser).Error
}

// Friendship management methods:
func (r *GormUserRepository) AddFriend(user, friend int64, ctx context.Context) error {
	friendship := model.Friendship{
		RequesterID: user,
		AccepterID:  friend,
	}
	return r.DB.WithContext(ctx).Create(&friendship).Error
}
func (r *GormUserRepository) RemoveFriend(user, friend int64, ctx context.Context) error {
	friendship := model.Friendship{
		RequesterID: user,
		AccepterID:  friend,
	}
	return r.DB.WithContext(ctx).Delete(&friendship).Error
}
func (r *GormUserRepository) GetFriends(user int64, ctx context.Context) ([]model.User, error) {
	var friends []model.User

	err := r.DB.WithContext(ctx).
		Joins("JOIN friendships ON users.id = friendships.accepter").
		Where("friendships.requester = ?", user).
		Find(&friends).Error

	if err != nil {
		return nil, err
	}

	return friends, nil
}
