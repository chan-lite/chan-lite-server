package fourChan

import (
	"chan-lite-server/src/database"
	"chan-lite-server/src/services"
	"net/http"

	"github.com/bahlo/goat"
)

// GetThread - TODO
func GetThread(w http.ResponseWriter, r *http.Request, p goat.Params) {
	services.SetHeaderAll(w)

	board := p["board"]
	thread := p["thread"]
	url := "https://a.4cdn.org/" + board + "/thread/" + thread + ".json"
	data := new(database.ChanThreadPage)

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
func SaveThread(w http.ResponseWriter, r *http.Request, p goat.Params, userStringID string) {

	// Get data from post and get.
	board := p["board"]
	if len(board) < 1 {
		services.ErrorMessage(w, "No board specified")
	}
	thread := p["thread"]
	if len(thread) < 1 {
		services.ErrorMessage(w, "No thread specified")
	}

	// Get thread
	url := "https://a.4cdn.org/" + board + "/thread/" + thread + ".json"
	threadData := new(database.ChanThreadPage)

	requestError := services.GoroutineRequest(url, threadData)
	if requestError != nil {
		services.ErrorMessage(w, "Error talking to 4chan servers")
		return
	}

	saveThreadError := database.SaveThread(userStringID, board, thread, threadData)
	if saveThreadError != nil {
		services.ErrorMessage(w, "You have already saved this thread")
		return
	}

	services.SuccessMessage(w, "Thread has been saved")
}

// GetSavedThread - TODO
func GetSavedThread(w http.ResponseWriter, r *http.Request, p goat.Params, userStringID string) {
	board := p["board"]
	thread := p["thread"]
	if len(board) < 1 || len(thread) < 1 {
		services.ErrorMessage(w, "Check parameters")
	}

	data, err := database.GetSavedThread(userStringID, board, thread)
	if err != nil {
		services.ErrorMessage(w, "Could not retreive saved thread")
		return
	}

	jsonString, jsonError := services.GoroutineToJSON(data)
	if jsonError != nil {
		services.ErrorMessage(w, "Error building data for client")
		return
	}

	services.Success(w, jsonString)
}
