package main

import (
	"proxy-chan/src/router"
)

func main() {
	err := router.Router()
	if err != nil {
		panic(err)
	}
}
