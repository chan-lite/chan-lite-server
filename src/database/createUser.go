package database

import (
	"errors"

	"github.com/jinzhu/gorm"
)

// UserTest - TODO
type UserTest struct {
	gorm.Model
	Name     string `gorm:"size:255"`
	Email    string `gorm:"size:255"`
	Device   string `gorm:"size:255"`
	Password []byte `gorm:"size:255"`
}

// CreateUser - TODO
func CreateUser(name string, email string, device string, hashedPassword []byte) (*UserTest, error) {
	// build user data
	user := UserTest{Name: name, Email: email, Device: device, Password: hashedPassword}
	// connect to database
	db, err := Connect()
	if err != nil {
		return &user, err
	}
	// check if user with email already exists
	count := 0
	db.Model(&UserTest{}).Where("email = ?", user.Email).Count(&count)
	if count != 0 {
		return &user, errors.New("User already exists")
	}
	// create user if they do not exist
	databaseUser := db.Create(&user)
	// close database
	defer db.Close()
	// return user
	return databaseUser.Value.(*UserTest), err
}
