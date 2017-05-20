package user

import (
	"net/http"
	"proxy-chan/src/database"
	"proxy-chan/src/services"

	"golang.org/x/crypto/bcrypt"

	"github.com/bahlo/goat"
	jwt "github.com/dgrijalva/jwt-go"
)

// CreateUser - TODO
func CreateUser(w http.ResponseWriter, r *http.Request, p goat.Params) {
	services.SetHeaderAll(w)
	// get data
	name := r.FormValue("name")
	email := r.FormValue("email")
	device := r.FormValue("device")
	password := r.FormValue("password")
	// check data
	if len(name) < 1 || len(email) < 1 || len(device) < 1 || len(password) < 1 {
		services.ErrorMessage(w, "Check parameters")
		return
	}
	// Hashing the password with the default cost of 10
	hashedPassword, passwordHashError := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if passwordHashError != nil {
		services.ErrorMessage(w, "Error hashing password")
		return
	}
	// attempt to create user
	user, createUserErr := database.CreateUser(name, email, device, hashedPassword)
	if createUserErr != nil {
		services.ErrorMessage(w, "User already exists")
		return
	}
	// get user token
	token, tokenErr := services.CreateUserToken(jwt.MapClaims{"ID": user.ID})
	if tokenErr != nil {
		services.ErrorMessage(w, "Error creating token")
		return
	}
	// build json to send to client
	data, toJSONErr := services.GoroutineToJSON(map[string]interface{}{
		"token":   token,
		"success": true,
	})
	if toJSONErr != nil {
		services.ErrorMessage(w, "An unexpected error has occurred")
		return
	}
	// send results to client
	services.Success(w, data)
}
