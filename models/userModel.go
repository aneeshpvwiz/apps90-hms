package models

import "gorm.io/gorm"

type User struct {
	ID       uint   `json:"id" gorm:"primary_key"`
	Email    string `json:"email" gorm:"type:varchar(254)"unique"`
	Password string `json:"password gorm:"type:varchar(254)"`
	gorm.Model
}
