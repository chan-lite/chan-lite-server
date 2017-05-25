package main

import (
	"github.com/chan-lite/chan-lite-server/src/router"
)

func main() {
	err := router.Router()
	if err != nil {
		panic(err)
	}
}
