package models

type Token struct {
	ID uint
	UserID uint
	User User
	RefreshToken string `gorm:"not null" validate:"format=jwt"`
}