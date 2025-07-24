package repository

import (
	"errors"
	"fmt"

	"github.com/UnendingLoop/users-api/cmd/internal/model"
	"github.com/jmoiron/sqlx"
)

type UserRepository struct {
	DB *sqlx.DB
}

var ErrUserNotFound = errors.New("user not found")

func NewUserRepository(db *sqlx.DB) *UserRepository {
	return &UserRepository{DB: db}
}

func (r *UserRepository) CreateUser(user *model.User) error {
	res, err := r.DB.Exec("INSERT INTO users(name, surname, email) VALUES (?,?,?)", user.Name, user.Surname, user.Email)
	if err != nil {
		return err
	}
	user.ID, _ = res.LastInsertId()
	return nil
}

func (r *UserRepository) GetUserByID(id int64) (*model.User, error) {
	var user model.User
	err := r.DB.Get(&user, "SELECT * FROM users WHERE ID = ?", id)
	if err != nil {
		return nil, err
	}
	return &user, nil
}

func (r *UserRepository) ListUsers() ([]model.User, error) {
	var users []model.User
	err := r.DB.Select(&users, "SELECT * FROM users")
	return users, err
}

func (r *UserRepository) DeleteUser(id int64) error {
	res, err := r.DB.Exec("DELETE FROM users WHERE id = ?", id)
	if err != nil {
		return err
	}
	if rows, _ := res.RowsAffected(); rows == 0 {
		return ErrUserNotFound
	}
	return nil
}

func (r *UserRepository) UpdateUser(user *model.User) error {
	//проверка на ненулевой input
	if user.Email == "" && user.Name == "" && user.Surname == "" {
		return fmt.Errorf("all fields are empty")
	}

	//загрузить из базы юзера с этим id - сразу проверить существует ли такой юзер
	var dbUser model.User
	err := r.DB.Get(&dbUser, "SELECT * FROM users WHERE id = ?", user.ID)
	if err != nil {
		return err
	}

	//проверить наличие нового имейл в базе
	var existingID int64
	err = r.DB.Get(&existingID, "SELECT id FROM users WHERE email = ? AND id != ?", user.Email, dbUser.ID)
	if err == nil || existingID != 0 {
		return fmt.Errorf("email already exists")
	}

	//скопировать в нулевые поля user поля из dbUser
	if user.Name == "" {
		user.Name = dbUser.Name
	}
	if user.Surname == "" {
		user.Surname = dbUser.Surname
	}
	if user.Email == "" {
		user.Email = dbUser.Email
	}

	//обновить юзера в БД
	_, err = r.DB.Exec("UPDATE users SET name = ?, surname = ?, email = ?  WHERE id = ?", user.Name, user.Surname, user.Email, user.ID)
	if err != nil {
		return fmt.Errorf("failed to update user: %v", err)
	}
	return nil
}
