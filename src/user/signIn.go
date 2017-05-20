package user

import (
	"net/http"
	"proxy-chan/src/database"
	"proxy-chan/src/services"

	"golang.org/x/crypto/bcrypt"

	"github.com/bahlo/goat"
	jwt "github.com/dgrijalva/jwt-go"
)

// SignIn - TODO
func SignIn(w http.ResponseWriter, r *http.Request, p goat.Params) {
	services.SetHeaderAll(w)
	// get post params
	email := r.FormValue("email")
	password := r.FormValue("password")
	// check post params
	if len(email) < 1 || len(password) < 1 {
		services.ErrorMessage(w, "Check parameters")
		return
	}
	// get databse user based on params
	databaseUser, databaseError := database.GetUserByEmail(email)
	if databaseError != nil {
		services.ErrorMessage(w, "User does not exist")
		return
	}
	// compare passwords
	hashedPassword := []byte(databaseUser.Password)
	bytePassword := []byte(password)
	passwordCompareError := bcrypt.CompareHashAndPassword(hashedPassword, bytePassword)
	if passwordCompareError != nil {
		services.ErrorMessage(w, "Incorrect password")
		return
	}
	// get token if passwords match
	token, tokenErr := services.CreateUserToken(jwt.MapClaims{"ID": databaseUser.ID})
	if tokenErr != nil {
		services.ErrorMessage(w, "Error creating token")
		return
	}
	// prepare json to return to client on success
	data, toJSONErr := services.GoroutineToJSON(map[string]interface{}{
		"token":   token,
		"success": true,
	})
	if toJSONErr != nil {
		services.ErrorMessage(w, "An unexpected error has occurred")
		return
	}
	services.Success(w, data)
}
