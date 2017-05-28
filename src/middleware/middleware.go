package middleware

import (
	"chan-lite-server/src/services"
	"net/http"
	"strconv"

	"github.com/bahlo/goat"
)

// RouteStategy - TODO
type RouteStategy func(w http.ResponseWriter, r *http.Request, p goat.Params)

// StadardAuthRoute - TODO
type StadardAuthRoute func(w http.ResponseWriter, r *http.Request, p goat.Params, userStringID string)

// Auth - TODO
func Auth(routeHandler StadardAuthRoute) RouteStategy {
	return func(w http.ResponseWriter, r *http.Request, p goat.Params) {

		// Set default headers.
		services.SetHeaderAll(w)

		// jwt token auth
		tokenString := r.FormValue("token")
		if len(tokenString) < 1 {
			services.ErrorMessage(w, "No token found")
			return
		}

		// Decode token.
		decodedToken, decodeError := services.DecodeToken(tokenString)
		if decodeError != nil {
			services.ErrorMessage(w, "Invalid token")
			return
		}

		// Receive data from token.
		tokenData := services.GetDataFromToken(decodedToken)
		tokenInvalid := services.CheckToken(tokenData)
		if tokenInvalid != nil {
			services.ErrorMessage(w, "Token has expired")
			return
		}

		// Ensure user ID is present.
		userID := tokenData["ID"].(float64)
		userStringID := strconv.FormatFloat(userID, 'f', -1, 64)
		if len(userStringID) < 1 {
			services.ErrorMessage(w, "No user ID found in token")
			return
		}

		routeHandler(w, r, p, userStringID)
	}
}
