package fourChan

import (
	"chan-lite-server/src/database"
	"chan-lite-server/src/services"
	"net/http"
	"strconv"

	"github.com/bahlo/goat"
)

type chanBoardPageBoard struct {
	Threads []chanBoardPageThread
}

type chanBoardPageThread struct {
	Posts []chanBoardPagePost
}

type chanBoardPagePost struct {
	No       int64
	Now      string
	Name     string
	Com      string
	Filename string
	Ext      string
	W        int64
	H        int64
	Tn_w     int64
	Tn_H     int64
	Tim      int64
	Time     int64
}

// GetBoard - TODO
func GetBoard(w http.ResponseWriter, r *http.Request, p goat.Params) {
	services.SetHeaderAll(w)

	board := p["board"]
	page := p["page"]
	url := "https://a.4cdn.org/" + board + "/" + page + ".json"
	data := new(chanBoardPageBoard)

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

// GetSavedBoard - TODO
func GetSavedBoard(w http.ResponseWriter, r *http.Request, p goat.Params, userStringID string) {
	board := p["board"]
	page, pageErr := strconv.Atoi(p["page"])
	perPage, perPageErr := strconv.Atoi(p["perPage"])
	if len(board) < 1 || pageErr != nil || perPageErr != nil {
		services.ErrorMessage(w, "Check parameters")
		return
	}

	data, err := database.GetSavedBoard(board, userStringID, page, perPage)
	if err != nil {
		services.ErrorMessage(w, "Error receiving saved board for user")
		return
	}

	jsonString, jsonError := services.GoroutineToJSON(data)
	if jsonError != nil {
		services.ErrorMessage(w, "Error building data for client")
		return
	}

	services.Success(w, jsonString)
}
