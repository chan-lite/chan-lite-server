package fourChan

import (
	"net/http"
	"chan-lite-server/src/services"

	"github.com/bahlo/goat"
)

type chanLandingPageBoards struct {
	Boards []chanLandingPageBoard
}

type chanLandingPageBoard struct {
	Board string
	Title string
}

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
