package database

// GetUser - TODO
func GetUser(ID string) (*UserTest, error) {

	user := new(UserTest)

	db, err := Connect()
	if err != nil {
		return user, err
	}

	db.Where("ID = ?", ID).First(&user)

	defer db.Close()

	return user, err
}

// GetUserByEmail - TODO
func GetUserByEmail(email string) (*UserTest, error) {
	user := new(UserTest)

	db, err := Connect()
	if err != nil {
		return user, err
	}

	db.Where("email = ?", email).First(&user)

	defer db.Close()

	return user, err
}
