package model

import "time"

type User struct {
	ID      int64  `gorm:"primaryKey" json:"id"`
	Name    string `gorm:"not null" json:"name"`
	Surname string `gorm:"not null" json:"surname"`
	Email   string `gorm:"uniqueIndex;not null" json:"email"`

	Friends  []*Friendship `gorm:"foreignKey:RequesterID"`
	FriendOf []*Friendship `gorm:"foreignKey:AccepterID"`
}

type Friendship struct {
	RequesterID int64     `gorm:"primaryKey;column:requester"`
	AccepterID  int64     `gorm:"primaryKey;column:accepter"`
	CreatedAt   time.Time `gorm:"column:created_at"`

	Requester *User `gorm:"foreignKey:RequesterID;references:ID"`
	Accepter  *User `gorm:"foreignKey:AccepterID;references:ID"`
}
