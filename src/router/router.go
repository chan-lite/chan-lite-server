package router

import (
	"chan-lite-server/src/fourChan"
	"chan-lite-server/src/middleware"
	"chan-lite-server/src/user"
	"fmt"
	"net/http"
	"os"
	"sync"

	"github.com/bahlo/goat"
)

type singleRoute struct {
	route   string
	name    string
	handler func(w http.ResponseWriter, r *http.Request, p goat.Params)
}

type multipleRoutes struct {
	routes    []singleRoute
	routeType string
}

type routeChannelHold struct {
	complete bool
}

// Router - TODO
func Router() error {
	// Create router channel.
	router := make(chan *goat.Router)
	// Spin up Goroutine to init channel
	go func() {
		router <- goat.New()
	}()
	// Define get routes.
	getRoutes := []singleRoute{
		singleRoute{
			route:   "/chan/landing",
			name:    "chan_landing",
			handler: fourChan.GetLanding},
		singleRoute{
			route:   "/chan/board/:board/:page",
			name:    "chan_board",
			handler: fourChan.GetBoard},
		singleRoute{
			route:   "/chan/thread/:board/:thread",
			name:    "chan_thread",
			handler: fourChan.GetThread},
		singleRoute{
			route:   "/chan/image",
			name:    "chan_image",
			handler: fourChan.GetImage}}
	// Define post routes.
	postRoutes := []singleRoute{
		singleRoute{
			route:   "/chan/user/create",
			name:    "chan_user_create",
			handler: user.CreateUser},
		singleRoute{
			route:   "/chan/user/get",
			name:    "chan_user_get",
			handler: user.GetUser},
		singleRoute{
			route:   "/chan/user/signin",
			name:    "chan_user_signin",
			handler: user.SignIn},
		singleRoute{
			route:   "/chan/user/save/get/landing",
			name:    "chan_user_thread_save_get",
			handler: middleware.Auth(fourChan.GetSavedLanding)},
		singleRoute{
			route:   "/chan/user/save/get/board/:board/:page/:perPage",
			name:    "chan_user_thread_save_get_board",
			handler: middleware.Auth(fourChan.GetSavedBoard)},
		singleRoute{
			route:   "/chan/user/save/thread/:board/:thread",
			name:    "chan_user_thread_save",
			handler: middleware.Auth(fourChan.SaveThread)}}
	// Receive router from channel.
	r := <-router
	close(router)
	// Create wait group.
	var wg sync.WaitGroup
	wg.Add(len(getRoutes))
	wg.Add(len(postRoutes))
	// Iterator through routes and set
	// with goroutines
	go func(routers []multipleRoutes) {
		for e := 0; e < len(routers); e++ {
			go func(routes []singleRoute, handlerType string) {
				for i := 0; i < len(routes); i++ {
					route := routes[i]
					switch handlerType {
					case "GET":
						r.Get(route.route, route.name, route.handler)
					case "POST":
						r.Post(route.route, route.name, route.handler)
					default:
						// Silence is bliss
					}
					wg.Done()
				}
			}(routers[e].routes, routers[e].routeType)
		}
	}([]multipleRoutes{
		multipleRoutes{routes: getRoutes, routeType: "GET"},
		multipleRoutes{routes: postRoutes, routeType: "POST"},
	})
	// Create goroutine to return value
	returnError := make(chan error)
	go func() {
		wg.Wait()
		// Eagerly tell user we're listening
		fmt.Println("listening on http://localhost:" + os.Getenv("PORT"))
		// // Begin listening on port.
		returnError <- r.Run(":" + os.Getenv("PORT"))
	}()
	// Get value to return.
	returnValue := <-returnError
	close(returnError)
	return returnValue
}
