package model

import "gorm.io/gorm"

type UserCredential struct {
	gorm.Model
	Id       uint32 `json:"id" gorm:"primaryKey"`
	Username string `json:"username" gorm:"unique;not null"`
	Password string `json : "password"`
	Role     string `json:"role" gorm:"not null"`
	Balance  int64  `json:"balance"`
}
