package controllers

import (
	"example/jwt-auth/db"
	"example/jwt-auth/models"
	"example/jwt-auth/response"
	"example/jwt-auth/validate"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"
)

func Registration(c *gin.Context) {
	var user models.User

	if err := c.BindJSON(&user); err != nil {
		verr :=  err.(validator.ValidationErrors)
		c.JSON(http.StatusBadRequest, validate.ErrorParse(verr))
		return
	}
	
	var existingUser models.User
	if err := db.GormDB.Where("email = ?", user.Email).First(&existingUser).Error; err != gorm.ErrRecordNotFound {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Email already exists"})
		return
	}
	
	result := db.GormDB.Create(&user)
	if result.Error != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": result.Error.Error()})
		return
	}

	userResponse := response.UserResponse{
		Model: user.Model,
		ID: user.ID,
		Email: user.Email,
	}
	accessToken, refreshToken := GenerateToken(&userResponse)
	SaveToken(user.ID, refreshToken)
	c.SetCookie("refreshToken", refreshToken, 2592000,"/","localhost",false,true);

	c.JSON(http.StatusOK, gin.H{"accessToken": accessToken, "refreshToken": refreshToken, "user": &userResponse})
}
func Login(c *gin.Context) {
	var user models.User

	if err := c.BindJSON(&user); err != nil {
		verr :=  err.(validator.ValidationErrors)
		c.JSON(http.StatusBadRequest, validate.ErrorParse(verr))
		return
	}

	var existingUser models.User
	if err := db.GormDB.Where("email = ?", user.Email).First(&existingUser).Error; err == gorm.ErrRecordNotFound {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Wrong email/password"})
		return
	}
	if err := bcrypt.CompareHashAndPassword([]byte(existingUser.Password), []byte(user.Password)); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Wrong email/password"})
		return
	}

	userResponse := response.UserResponse{
		Model: existingUser.Model,
		ID: existingUser.ID,
		Email: existingUser.Email,
	}
	accessToken, refreshToken := GenerateToken(&userResponse)
	SaveToken(existingUser.ID, refreshToken)


	c.SetCookie("refreshToken", refreshToken, 2592000,"/","localhost",false,true);
	c.JSON(http.StatusOK, gin.H{"accessToken": accessToken, "refreshToken": refreshToken, "user": &userResponse})
}
func Logout(c *gin.Context){
	refreshToken,_ := c.Cookie("refreshToken")
	c.SetCookie("refreshToken", "", -1, "/", "localhost", false, true)
	RemoveToken(refreshToken)
	c.JSON(http.StatusOK, gin.H{"status": "Logged out"})
}
func Refresh(c *gin.Context){
	refreshToken, err := c.Cookie("refreshToken")
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"status": "Access denied"})
		return
	}
	RefreshToken(refreshToken)
}