package user

import (
	"chan-lite-server/src/database"
	"chan-lite-server/src/services"
	"net/http"
	"strconv"

	"github.com/bahlo/goat"
)

// GetUser - get basic user data with a token
func GetUser(w http.ResponseWriter, r *http.Request, p goat.Params) {
	services.SetHeaderAll(w)

	tokenString := r.FormValue("token")
	if len(tokenString) < 1 {
		services.ErrorMessage(w, "No token found")
	}

	decodedToken, decodeError := services.DecodeToken(tokenString)
	if decodeError != nil {
		services.ErrorMessage(w, "Invalid token")
		return
	}

	tokenData := services.GetDataFromToken(decodedToken)
	tokenInvalid := services.CheckToken(tokenData)
	if tokenInvalid != nil {
		services.ErrorMessage(w, "Token has expired")
		return
	}

	userID := tokenData["ID"].(float64)
	userStringID := strconv.FormatFloat(userID, 'f', -1, 64)

	userData, getUserErr := database.GetUser(userStringID)
	if getUserErr != nil {
		services.ErrorMessage(w, "Error getting user from database")
		return
	}

	services.SuccessMessage(w, "User's name is "+userData.Name)
}
