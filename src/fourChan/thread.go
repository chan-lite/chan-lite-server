package fourChan

import (
	"chan-lite-server/src/database"
	"chan-lite-server/src/services"
	"net/http"
	"strconv"

	"github.com/bahlo/goat"
)

type chanThreadPage struct {
	Posts []chanThreadPost
}

type chanThreadPost struct {
	No       int64
	Now      string
	Name     string
	Com      string
	Filename string
	Ext      string
	W        int64
	H        int64
	Tn_W     int64
	Tn_H     int64
	Tim      int64
	Time     int64
}

// GetThread - TODO
func GetThread(w http.ResponseWriter, r *http.Request, p goat.Params) {
	services.SetHeaderAll(w)

	board := p["board"]
	thread := p["thread"]
	url := "https://a.4cdn.org/" + board + "/thread/" + thread + ".json"
	data := new(chanThreadPage)

	requestError := services.GoroutineRequest(url, data)
	if requestError != nil {
		services.ErrorMessage(w, "Error talking to 4chan servers")
		return
	}

	jsonString, jsonError := services.GoroutineToJSON(data)
	if jsonError != nil {
		services.ErrorMessage(w, "Error creating JSON")
		return
	}

	services.Success(w, jsonString)
}

// SaveThread - TODO
func SaveThread(w http.ResponseWriter, r *http.Request, p goat.Params) {
	// Set default headers.
	services.SetHeaderAll(w)
	// Get data from post and get.
	board := p["board"]
	if len(board) < 1 {
		services.ErrorMessage(w, "No board specified")
	}
	thread := p["thread"]
	if len(thread) < 1 {
		services.ErrorMessage(w, "No thread specified")
	}
	tokenString := r.FormValue("token")
	if len(tokenString) < 1 {
		services.ErrorMessage(w, "No token found")
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
		services.ErrorMessage(w, "No user ID found")
	}

	// Get thread
	url := "https://a.4cdn.org/" + board + "/thread/" + thread + ".json"
	data := new(chanThreadPage)

	requestError := services.GoroutineRequest(url, data)

	if requestError != nil {
		services.ErrorMessage(w, "Error talking to 4chan servers")
		return
	}

	threadDataJSON, threadDataError := services.GoroutineToJSON(data)

	if threadDataError != nil {
		services.ErrorMessage(w, "Error creating JSON")
		return
	}

	database.SaveThread(userStringID, threadDataJSON)

	services.SuccessMessage(w, userStringID)

}
