package models

import "gorm.io/gorm"

type User struct {
	ID       uint   `json:"id" gorm:"primary_key"`
	Email    string `json:"email" gorm:"unique"`
	Password string `json:"password"`
	gorm.Model
}
