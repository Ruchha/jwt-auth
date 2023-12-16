package db

import (
	"example/jwt-auth/models"
	"fmt"
	"log"
	"os"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)
var GormDB *gorm.DB
func ConnectDatabase(){
    var err error
    dbuser, dbpass, address, dbname := os.Getenv("DBUSER"), os.Getenv("DBPASS"), os.Getenv("ADDRESS"), os.Getenv("DBNAME")
	dsn := fmt.Sprintf("%v:%v@tcp(%v)/%v?charset=utf8mb4&parseTime=True&loc=Local", dbuser, dbpass, address, dbname)
    GormDB, err = gorm.Open(mysql.Open(dsn), &gorm.Config{})
    if err != nil {
        log.Fatal("Error connecting database")
    }
	if err != nil {
		log.Fatal("error", err)
    }
	GormDB.AutoMigrate(&models.User{}, &models.Token{})
}