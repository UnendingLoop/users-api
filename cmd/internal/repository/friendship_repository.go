package repository

import (
	"context"

	"github.com/UnendingLoop/users-api/cmd/internal/model"
	"gorm.io/gorm"
)

// FriendRepository определяет контракт для работы с дружескими связями между пользователями.
// Интерфейс может быть реализован с использованием любой СУБД или подхода к хранению данных.
type FriendRepository interface {
	// AddFriend создает новую связь дружбы между двумя пользователями.
	AddFriend(ctx context.Context, friendship *model.Friendship) error

	// RemoveFriend удаляет существующую связь дружбы между двумя пользователями.
	RemoveFriend(ctx context.Context, friendship *model.Friendship) error

	// GetFriends возвращает список пользователей, являющихся друзьями указанного пользователя.
	GetFriends(ctx context.Context, user int64) ([]model.User, error)
}

// GormFriendRepository — реализация FriendRepository на базе GORM ORM.
type GormFriendRepository struct {
	DB *gorm.DB
}

// NewGormFriendRepository создает новый экземпляр GormFriendRepository с переданной GORM-базой данных.
func NewGormFriendRepository(db *gorm.DB) *GormFriendRepository {
	return &GormFriendRepository{DB: db}
}

func (r *GormFriendRepository) AddFriend(ctx context.Context, friendship *model.Friendship) error {
	return r.DB.WithContext(ctx).Create(&friendship).Error
}
func (r *GormFriendRepository) RemoveFriend(ctx context.Context, friendship *model.Friendship) error {
	return r.DB.WithContext(ctx).Delete(&friendship).Error
}
func (r *GormFriendRepository) GetFriends(ctx context.Context, user int64) ([]model.User, error) {
	var friends []model.User

	err := r.DB.WithContext(ctx).
		Joins("JOIN friendships ON users.id = friendships.accepter").
		Where("friendships.requester = ?", user).
		Find(&friends).Error

	return friends, err
}
