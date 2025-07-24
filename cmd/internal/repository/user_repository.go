package repository

import (
	"errors"

	"github.com/UnendingLoop/users-api/cmd/internal/model"
	"gorm.io/gorm"
)

type GormUserRepository struct {
	DB *gorm.DB
}

var ErrUserNotFound = errors.New("user not found")
var ErrEmailExists = errors.New("email already exists")
var ErrEmptyfields = errors.New("all fields are empty")

func NewGormUserRepository(db *gorm.DB) *GormUserRepository {
	return &GormUserRepository{DB: db}
}

func (r *GormUserRepository) CreateUser(user *model.User) error {
	return r.DB.Create(user).Error
}

func (r *GormUserRepository) GetUserByID(id int64) (*model.User, error) {
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

func (r *GormUserRepository) ListUsers() ([]model.User, error) {
	var users []model.User
	err := r.DB.Find(&users).Error
	return users, err
}

func (r *GormUserRepository) DeleteUser(id int64) error {
	res := r.DB.Delete(&model.User{}, id)
	if res.RowsAffected == 0 {
		return ErrUserNotFound
	}
	return res.Error
}

func (r *GormUserRepository) UpdateUser(user *model.User) error {
	//проверка на ненулевой input
	if user.Email == "" && user.Name == "" && user.Surname == "" {
		return ErrEmptyfields
	}

	//загрузить из базы юзера с этим id - сразу проверить существует ли такой юзер
	var dbUser model.User
	err := r.DB.First(&dbUser, user.ID).Error
	if err != nil {
		return err
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
