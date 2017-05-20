package fourChan

import (
	"net/http"
	"proxy-chan/src/services"

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
	Tn_w     int64
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
