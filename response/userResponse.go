package response

import "gorm.io/gorm"

type UserResponse struct {
	gorm.Model
	ID    uint
	Email string
}