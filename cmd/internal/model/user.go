package model

type User struct {
	ID      int64  `db:"id"`
	Name    string `db:"name"`
	Surname string `db:"surname"`
	Email   string `db:"email"`
}
