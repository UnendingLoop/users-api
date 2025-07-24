package model

type User struct {
	ID      int64  `gorm:"primaryKey" json:"id"`
	Name    string `gorm:"not null" json:"name"`
	Surname string `gorm:"not null" json:"surname"`
	Email   string `gorm:"uniqueIndex;not null" json:"email"`
}
