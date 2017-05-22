package database

import (
	"os"

	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/mysql"
	// _ "github.com/jinzhu/gorm/dialects/postgres"
)

// Connect - TODO
func Connect() (*gorm.DB, error) {

	// mysql
	connectionInfo := os.Getenv("DB_USER") + ":" + os.Getenv("DB_PASSWORD") + "@tcp(" + os.Getenv("DB_HOST") + ":" + os.Getenv("DB_PORT") + ")/" + os.Getenv("DB_NAME")
	connectionSettings := "?charset=utf8&parseTime=True&loc=Local"
	db, err := gorm.Open("mysql", connectionInfo+connectionSettings)

	// postgres
	// user := "user=" + os.Getenv("DB_USER")
	// password := "password=" + os.Getenv("DB_PASSWORD")
	// host := "host=" + os.Getenv("DB_HOST")
	// name := "user=" + os.Getenv("DB_NAME")
	// db, err := gorm.Open("postgres", host+" "+user+" "+name+" sslmode=disable "+password)

	return db, err
}
