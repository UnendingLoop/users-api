package repository

import (
	"context"

	"github.com/UnendingLoop/users-api/cmd/internal/model"
	"gorm.io/gorm"
)

type FriendRepository interface {
	AddFriend(friendship *model.Friendship, ctx context.Context) error
	RemoveFriend(friendship *model.Friendship, ctx context.Context) error
	GetFriends(user int64, ctx context.Context) ([]model.User, error)
}

type GormFriendRepository struct {
	DB *gorm.DB
}

func NewGormFriendRepository(db *gorm.DB) *GormFriendRepository {
	return &GormFriendRepository{DB: db}
}

func (r *GormFriendRepository) AddFriend(friendship *model.Friendship, ctx context.Context) error {
	return r.DB.WithContext(ctx).Create(&friendship).Error
}
func (r *GormFriendRepository) RemoveFriend(friendship *model.Friendship, ctx context.Context) error {
	return r.DB.WithContext(ctx).Delete(&friendship).Error
}
func (r *GormFriendRepository) GetFriends(user int64, ctx context.Context) ([]model.User, error) {
	var friends []model.User

	err := r.DB.WithContext(ctx).
		Joins("JOIN friendships ON users.id = friendships.accepter").
		Where("friendships.requester = ?", user).
		Find(&friends).Error

	return friends, err
}
