package models

import (
	"fmt"
	"log"
	"os"

	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/postgres"
	"github.com/joho/godotenv"
)

var db *gorm.DB

func init() {
	err := godotenv.Load()
	if err != nil {
		log.Fatal(err)
		return
	}

	dbName := os.Getenv("db_name")
	dbUser := os.Getenv("db_user")
	dbPassword := os.Getenv("db_password")
	dbType := os.Getenv("db_type")
	dbHost := os.Getenv("db_host")

	conn, err := gorm.Open(dbType, fmt.Sprintf("host=%s user=%s dbname=%s sslmode=disable password=%s", dbHost, dbUser, dbName, dbPassword))
	if err != nil {
		log.Fatal(err)
		return
	}

	db = conn
	db.Debug().AutoMigrate(&User{}, &Quote{})
}

func DB() *gorm.DB {
	return db
}
