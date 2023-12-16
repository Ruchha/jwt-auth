package main

import (
	"example/jwt-auth/controllers"
	"example/jwt-auth/db"
	"example/jwt-auth/initializers"
	"example/jwt-auth/validate"
	"github.com/gin-gonic/gin"
)

func main() {
	initializers.LoadEnvVariables()
	db.ConnectDatabase()
	validate.LoadValidator()
	r := gin.New()
	r.POST("/registration", controllers.Registration)
	r.POST("/login", controllers.Login)
	r.POST("/logout", controllers.Logout)
	r.GET("/refresh", controllers.Refresh)
	// r.GET("/users")
	r.Run()
}