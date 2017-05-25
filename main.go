package main

import (
	"github.com/chan-lite/chan-lite-server/src/router"
)

func main() {
	// We no longer need this as we're using
	// Heroku locally.
	// environmentError := godotenv.Load()
	// if environmentError != nil {
		// This is actually not an error
		// if we're on the Heroku server
		// for obvious reasons.
	// }
	routerError := router.Router()
	if routerError != nil {
		panic(routerError)
	}
}
