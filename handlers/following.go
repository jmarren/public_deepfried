package handlers

import (
	"fmt"
	"github.com/a-h/templ"
	"github.com/jmarren/deepfried/components"
	"github.com/jmarren/deepfried/services"
	"net/http"
)

type Following struct {
	JHandle
}

func NewFollowingHandler(j JHandle) *Following {
	f := new(Following)
	f.JHandle = j
	return f
}

func (f *Following) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	userService := services.NewUserService(f.Context())

	theirUsername := r.PathValue("username")

	var err error
	var component templ.Component

	fmt.Printf("f.Method: %s\n", f.Method)

	followButton := components.FollowButton(false, theirUsername)
	unfollowButton := components.FollowButton(true, theirUsername)

	switch f.Method {
	case "POST":
		err = userService.FollowUser(theirUsername)
		if err != nil {
			component = followButton
		} else {
			component = unfollowButton
		}
	case "DELETE":
		err = userService.UnFollowUser(theirUsername)
		if err != nil {
			component = unfollowButton
		} else {
			component = followButton
		}
	}

	component.Render(f.Context(), w)

}
