package controllers

import (
	"example/jwt-auth/db"
	"example/jwt-auth/models"
	"example/jwt-auth/response"
	"fmt"
	"log"
	"os"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"gorm.io/gorm"
)

func GenerateToken(payload *response.UserResponse) (string, string){
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"sub": payload,
		"exp": time.Now().Add(time.Minute * 60).Unix(),
	})
	accessToken, err := token.SignedString([]byte(os.Getenv("JWT_ACCESS_SECRET")))
	if err != nil {
		log.Fatal("Failed generating access token")
	}
	token = jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"sub": payload,
		"exp": time.Now().Add(time.Hour * 24 * 30).Unix(),
	})
	if err != nil {
		log.Fatal("Failed generating refresh token")
	}
	refreshToken, err := token.SignedString([]byte(os.Getenv("JWT_REFRESH_SECRET")))
	if err != nil {
		log.Fatal("Failed generating refresh token")
	}
	return accessToken, refreshToken
}

func SaveToken(userId uint, refreshToken string) models.Token{
	var token models.Token
	if err := db.GormDB.Where("user_id = ?", userId).First(&token).Error; err != gorm.ErrRecordNotFound {
		token.RefreshToken = refreshToken
		db.GormDB.Save(&token)
		//db.GormDB.Model(&token).Update("refresh_token", refreshToken)
	} else{
	token.UserID, token.RefreshToken = userId, refreshToken
	db.GormDB.Create(&token)
	}
	return token
}

func RemoveToken(refreshToken string) {
	db.GormDB.Where("refresh_token = ?", refreshToken).Delete(&models.Token{})
}

func RefreshToken(refreshToken string){
	mySigningKey := []byte(os.Getenv("JWT_REFRESH_SECRET"))

	tokenString := refreshToken

	token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
		}
		return mySigningKey, nil
	})

	if claims, ok := token.Claims.(jwt.MapClaims); ok && token.Valid {
		fmt.Println(claims)
	} else {
		fmt.Println(err)
	}
}

