package services

import (
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"chan-lite-server/src/sensitive"
	"sync"
	"time"

	jwt "github.com/dgrijalva/jwt-go"
)

var client = &http.Client{Timeout: 10 * time.Second}
var defaultBitSize = 2048

// Token will expire after one (1) hour.
var expirationTime int64 = 3600000

// SetHeaderAll - TODO
func SetHeaderAll(w http.ResponseWriter) {
	w.Header().Set("Access-Control-Allow-Origin", "*")
}

// CreateUserToken - TODO
func CreateUserToken(data jwt.MapClaims) (string, error) {
	data["granted"] = time.Now().String()
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, data)
	return token.SignedString([]byte(sensitive.TokenSecret))
}

// CheckToken - TODO
func CheckToken(token jwt.MapClaims) error {
	granted := token["granted"].(string)
	layout := "2006-01-02 15:04:05"
	aTime, err := time.Parse(layout, granted[:19])
	if err != nil {
		return err
	}
	milli := aTime.UnixNano() / int64(time.Millisecond)
	currentMilli := time.Now().UnixNano() / int64(time.Millisecond)

	compare := (currentMilli - milli) / 1000
	if compare > expirationTime {
		return errors.New("Token has expired")
	}

	return nil
}

// DecodeToken - TODO
func DecodeToken(tokenString string) (*jwt.Token, error) {
	token, parseError := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		return []byte(sensitive.TokenSecret), nil
	})
	return token, parseError
}

// GetDataFromToken - TODO
func GetDataFromToken(token *jwt.Token) jwt.MapClaims {
	tokenClaims := token.Claims.(jwt.MapClaims)
	return tokenClaims
}

// ErrorMessage - TODO
func ErrorMessage(w http.ResponseWriter, message string) {
	data, toJSONErr := GoroutineToJSON(map[string]interface{}{
		"message": message,
		"success": false,
	})
	if toJSONErr != nil {
		panic(toJSONErr)
	}
	fmt.Fprintf(w, "%s", data)
}

// Error - TODO
func Error(w http.ResponseWriter) {
	fmt.Fprintf(w, "Error")
}

// SuccessMessage - TODO
func SuccessMessage(w http.ResponseWriter, dataString string) {
	data, toJSONErr := GoroutineToJSON(map[string]interface{}{
		"message": dataString,
		"success": true,
	})
	if toJSONErr != nil {
		panic(toJSONErr)
	}
	Success(w, data)
}

// Success - TODO
func Success(w http.ResponseWriter, data []byte) {
	fmt.Fprintf(w, "%s", data)
}

// GoroutineGetRequest - TODO
func GoroutineGetRequest(url string) (*http.Response, error) {
	var wg sync.WaitGroup
	requestOut := make(chan *http.Response)
	requestErr := make(chan error)

	wg.Add(1)
	go func() {
		res, err := http.Get(url)
		requestOut <- res
		requestErr <- err
		wg.Done()
	}()

	go func() {
		wg.Wait()
		close(requestOut)
		close(requestErr)
	}()

	return <-requestOut, <-requestErr
}

// GoroutineRequest - TODO
func GoroutineRequest(url string, data interface{}) error {
	var wg sync.WaitGroup
	out := make(chan error)

	wg.Add(1)
	go func() {
		err := JSONServices(url, data)
		out <- err
		wg.Done()
	}()

	go func() {
		wg.Wait()
		close(out)
	}()

	return <-out
}

// GoroutineToJSON - TODO
func GoroutineToJSON(target interface{}) ([]byte, error) {
	var jsonWait sync.WaitGroup
	jsonOut := make(chan []byte)
	jsonErr := make(chan error)
	jsonWait.Add(1)

	go func() {
		json, err := json.Marshal(target)
		if err != nil {
			jsonErr <- err
		}
		jsonOut <- json
		jsonWait.Done()
	}()

	go func() {
		jsonWait.Wait()
		close(jsonOut)
		close(jsonErr)
	}()

	return <-jsonOut, <-jsonErr
}

// JSONServices - TODO
func JSONServices(url string, target interface{}) error {
	r, err := client.Get(url)
	if err != nil {
		return err
	}
	defer r.Body.Close()
	return json.NewDecoder(r.Body).Decode(target)
}
