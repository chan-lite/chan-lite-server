package router

import (
	"fmt"
	"os"
	"chan-lite-server/src/fourChan"
	"chan-lite-server/src/user"

	"github.com/bahlo/goat"
)

// Router - TODO
func Router() error {
	// create goat router
	r := goat.New()
	// 4chan routes
	r.Get("/chan/landing", "chan_landing", fourChan.GetLanding)
	r.Get("/chan/board/:board/:page", "chan_board", fourChan.GetBoard)
	r.Get("/chan/thread/:board/:thread", "chan_thread", fourChan.GetThread)
	r.Get("/chan/image", "chan_image", fourChan.GetImage)
	// user routes
	r.Post("/chan/user/create", "chan_user_create", user.CreateUser)
	r.Post("/chan/user/get", "chan_user_get", user.GetUser)
	r.Post("/chan/user/signin", "chan_user_signin", user.SignIn)
	// finished
	fmt.Println("listening on http://localhost:" + os.Getenv("PORT"))
	return r.Run(":" + os.Getenv("PORT"))
}
