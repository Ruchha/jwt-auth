package models

import (
	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"
)

type User struct {
	gorm.Model
	Email string `gorm:"unique;not null" binding:"required,email"`
	Password string `gorm:"not null" binding:"required,min=8"`
}

 func (u *User) BeforeCreate(tx *gorm.DB) (err error) {
 	hash, err := bcrypt.GenerateFromPassword([]byte(u.Password), bcrypt.DefaultCost)
 	u.Password = string(hash)
 	return
}
