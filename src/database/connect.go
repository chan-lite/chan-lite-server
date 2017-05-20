package database

import (
	"proxy-chan/src/sensitive"

	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/mysql"
)

// Connect - TODO
func Connect() (*gorm.DB, error) {
	connectionInfo := sensitive.DatabaseUser + ":" + sensitive.DatabasePassword + "@tcp(" + sensitive.DatabaseHost + ":" + sensitive.DatabasePort + ")/" + sensitive.DatabaseName
	connectionSettings := "?charset=utf8&parseTime=True&loc=Local"
	db, err := gorm.Open("mysql", connectionInfo+connectionSettings)
	return db, err
}
