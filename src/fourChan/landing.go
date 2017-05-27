package fourChan

import (
	"chan-lite-server/src/database"
	"chan-lite-server/src/services"
	"net/http"
	"strconv"

	"github.com/bahlo/goat"
)

type chanLandingPageBoards struct {
	Boards []chanLandingPageBoard
}

type chanLandingPageBoard struct {
	Board string
	Title string
}

var (
	url  = "https://a.4cdn.org/boards.json"
	data = new(chanLandingPageBoards)
)

// GetLanding - TODO
func GetLanding(w http.ResponseWriter, r *http.Request, p goat.Params) {
	services.SetHeaderAll(w)

	url := "https://a.4cdn.org/boards.json"
	data := new(chanLandingPageBoards)

	requestError := services.GoroutineRequest(url, data)

	if requestError != nil {
		services.Error(w)
		return
	}

	jsonString, jsonError := services.GoroutineToJSON(data)

	if jsonError != nil {
		services.Error(w)
		return
	}

	services.Success(w, jsonString)
}

// GetSavedLanding - TODO
func GetSavedLanding(w http.ResponseWriter, r *http.Request, p goat.Params) {
	services.SetHeaderAll(w)

	url := "https://a.4cdn.org/boards.json"
	data := new(chanLandingPageBoards)

	//Token check
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
	}

	requestError := services.GoroutineRequest(url, data)

	if requestError != nil {
		services.Error(w)
		return
	}

	// We now have the user token
	// request default boards
	// and match along with boards
	// the current user has saved.
	// Return common boards.

	usersSavedBoards, errorGettingUsersSavedBoards := database.GetSavedThreads(userStringID)
	if errorGettingUsersSavedBoards != nil {
		services.ErrorMessage(w, "Error receiving boards saved by user")
	}

	var merged []chanLandingPageBoard

	for e := 0; e < len(usersSavedBoards); e++ {
		for i := 0; i < len(data.Boards); i++ {
			currentBoard := data.Boards[i]
			currentSaved := usersSavedBoards[e]
			if currentSaved == currentBoard.Board {
				merged = append(merged, currentBoard)
			}
		}
	}

	jsonString, jsonError := services.GoroutineToJSON(chanLandingPageBoards{Boards: merged})

	if jsonError != nil {
		services.Error(w)
		return
	}

	services.Success(w, jsonString)
}
