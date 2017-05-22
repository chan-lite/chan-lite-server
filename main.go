package main

import (
	"chan-lite-server/src/router"

	"github.com/joho/godotenv"
)

func main() {
	environmentError := godotenv.Load()
	if environmentError != nil {
		// This is actually not an error
		// if we're on the Heroku server
		// for obvious reasons.
	}
	routerError := router.Router()
	if routerError != nil {
		panic(routerError)
	}
}
