package main

import (
	"chan-lite-server/src/router"
)

func main() {
	err := router.Router()
	if err != nil {
		panic(err)
	}
}
