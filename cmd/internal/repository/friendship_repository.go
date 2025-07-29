package repository

import (
	"context"

	"github.com/UnendingLoop/users-api/cmd/internal/model"
	"gorm.io/gorm"
)

type FriendRepository interface {
	AddFriend(user, friend int64, ctx context.Context) error
	RemoveFriend(user, friend int64, ctx context.Context) error
	GetFriends(user int64, ctx context.Context) ([]model.User, error)
}

type GormFriendRepository struct {
	DB *gorm.DB
}

func NewGormFriendRepository(db *gorm.DB) *GormFriendRepository {
	return &GormFriendRepository{DB: db}
}

// Friendship management methods:
func (r *GormFriendRepository) AddFriend(user, friend int64, ctx context.Context) error {
	friendship := model.Friendship{
		RequesterID: user,
		AccepterID:  friend,
	}
	return r.DB.WithContext(ctx).Create(&friendship).Error
}
func (r *GormFriendRepository) RemoveFriend(user, friend int64, ctx context.Context) error {
	friendship := model.Friendship{
		RequesterID: user,
		AccepterID:  friend,
	}
	return r.DB.WithContext(ctx).Delete(&friendship).Error
}
func (r *GormFriendRepository) GetFriends(user int64, ctx context.Context) ([]model.User, error) {
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
