package fourChan

import (
	"chan-lite-server/src/services"
	"io/ioutil"
	"net/http"

	"github.com/bahlo/goat"
)

// GetImage - TODO
func GetImage(w http.ResponseWriter, r *http.Request, p goat.Params) {
	services.SetHeaderAll(w)

	url := r.URL.Query().Get("image")

	imageResponse, imageError := services.GoroutineGetRequest(url)

	if imageError != nil {
		services.Error(w)
		return
	}

	defer imageResponse.Body.Close()

	bytesData, bytesError := ioutil.ReadAll(imageResponse.Body)

	if bytesError != nil {
		services.Error(w)
		return
	}

	services.Success(w, bytesData)
}
