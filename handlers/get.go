package handlers

import (
	// "context"
	"fmt"
	"github.com/jmarren/deepfried/services"
	"net/http"
)

type GetHandler struct {
	BaseServices
	jwtVerified bool
}

func NewGetHandler(jwtVerified bool, serv BaseServices) *GetHandler {
	g := new(GetHandler)
	g.BaseServices = serv
	g.jwtVerified = jwtVerified
	return g
}

// If the user exists in context, allow the handler
// If the user does not exist but their jwt has been verified, return the create account modal
// If neither is the case, redirect to the sign up page at /login (cognito path to require authentication)
func (g *GetHandler) Restricted(user *services.User, j JHandle, h ComponentHandler) ComponentHandler {
	if user != nil {
		return h
	} else if g.jwtVerified {
		return NewCreateAccountModal()
	} else {
		return NewAuthRedirect(j)
	}
}

// Each route has a handler, may or may not be a page,  and may or may not be restricted
type Route struct {
	isPage       bool
	handler      ComponentHandler
	isRestricted bool
}

// if the route is a page it's handler will be passed as a parameter to NewBase, which will
// render it is the request has the header "HX-Request" == "true", otherwise Base will fetch
// the component with GetComponent() and render the component inside of Base
func (g *GetHandler) applyRoutes(mux *http.ServeMux, r *http.Request) *http.ServeMux {

	jHandle := JHandle{
		r,
		g.UserService,
		g.BaseServices,
	}

	user, _ := g.GetFromCtx()

	RouteMap := make(map[string]Route)
	// RouteMap["GET /account-dropdown"] = Route{
	// 	false,
	// 	new(AccountDropdown),
	// 	true,
	// }
	RouteMap["GET /modal/upload"] = Route{
		false,
		NewUpload(),
		true,
	}
	RouteMap["GET /modal/create-account"] = Route{
		false,
		NewCreateAccountModal(),
		true,
	}
	RouteMap["GET /my-downloads"] = Route{
		true,
		NewMyDownloads(jHandle),
		true,
	}
	RouteMap["GET /my-uploads"] = Route{
		true,
		NewMyUploads(jHandle),
		true,
	}
	RouteMap["GET /feed"] = Route{
		true,
		NewFeed(jHandle),
		true,
	}
	RouteMap["GET /modal/edit-profile"] = Route{
		false,
		NewEditProfileModal(jHandle),
		true,
	}
	RouteMap["GET /modal/upload-form"] = Route{
		false,
		NewUploadForm(jHandle),
		true,
	}
	RouteMap["GET /"] = Route{
		true,
		NewExplore(jHandle),
		false,
	}
	RouteMap["GET /explore"] = Route{
		true,
		NewExplore(jHandle),
		false,
	}
	RouteMap["GET /search"] = Route{
		true,
		NewSearch(jHandle),
		false,
	}
	RouteMap["GET /search-bar-dropdown"] = Route{
		false,
		NewSearchDropdown(jHandle),
		false,
	}
	RouteMap["GET /modal/filters"] = Route{
		false,
		NewFilters(),
		false,
	}
	RouteMap["GET /modal/users/{username}/{relation}"] = Route{
		false,
		NewUsersModal(jHandle),
		false,
	}
	RouteMap["GET /player"] = Route{
		false,
		NewPlayer(jHandle),
		false,
	}

	for route, config := range RouteMap {
		handler := config.handler
		if config.isRestricted {
			handler = g.Restricted(user, jHandle, handler)
		}
		if config.isPage {
			mux.Handle(route, NewBase(jHandle, handler))
		} else {
			mux.Handle(route, handler)
		}
	}
	mux.Handle("GET /{user}", NewProfileRoute(jHandle))
	mux.Handle("GET /{user}/{track}", NewProfileRoute(jHandle))
	return mux
}

func GetRequestMiddleware(r *http.Request, w *http.ResponseWriter, h http.HandlerFunc) http.HandlerFunc {
	isPreload := r.Header.Get("HX-Preloaded")

	if isPreload == "true" {
		fmt.Println("preload request")
	}

	audioQueue := r.Header.Get("HX-Audio-Queue")
	fmt.Printf("audioQueue:  %v\n", audioQueue)

	fmt.Printf("r.URL.Path: %s\n", r.URL.Path)

	return func(w http.ResponseWriter, r *http.Request) {
		h.ServeHTTP(w, r)
	}
}

func (g *GetHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {

	// create new servemux
	mux := http.NewServeMux()

	isPreload := r.Header.Get("HX-Preloaded")

	if isPreload == "true" {
		fmt.Println("preload request")
	}

	audioQueue := r.Header.Get("HX-Audio-Queue")
	fmt.Printf("audioQueue:  %v\n", audioQueue)

	fmt.Printf("r.URL.Path: %s\n", r.URL.Path)

	// apply the component handler function
	g.applyRoutes(mux, r)

	// serve
	mux.ServeHTTP(w, r)
}
