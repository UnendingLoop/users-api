package repository

import (
	"github.com/UnendingLoop/users-api/cmd/internal/model"
	"github.com/jmoiron/sqlx"
)

type UserRepository struct {
	DB *sqlx.DB
}

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
	_, err := r.DB.Exec("DELETE FROM users WHERE id = ?", id)
	return err
}

func (r *UserRepository) UpdateUser(user *model.User) error {
	return nil
}
