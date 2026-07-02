package model

import (
	"gorm.io/gorm"
)

type User struct {
	gorm.Model
	UserName string `gorm:"type:varchar(50);uniqueIndex;not null" json:"user_name"`
	Password string `gorm:"not null" json:"-"`
	Email    string `gorm:"type:varchar(50);not null" json:"email"`
}
