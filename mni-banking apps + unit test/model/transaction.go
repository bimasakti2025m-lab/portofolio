package model

import "gorm.io/gorm"

type Transaction struct {
	gorm.Model
	ID         uint `gorm:"primaryKey"`
	FromUserID uint
	ToUserID   uint
	Amount     int64
	Type       string
	Status     string
}