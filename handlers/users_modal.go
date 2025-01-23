package handlers

import (
	// "context"
	// "fmt"
	"github.com/a-h/templ"
	"github.com/jmarren/deepfried/components"
	"github.com/jmarren/deepfried/services"

	// "github.com/jmarren/deepfried/services"
	"net/http"
)

type UsersModal struct {
	JHandle
}

func NewUsersModal(j JHandle) *UsersModal {
	u := new(UsersModal)
	u.JHandle = j
	return u
}

func (u *UsersModal) GetComponent() *templ.Component {
	relation := u.PathValue("relation")
	username := u.PathValue("username")

	if username == "me" {
		user, _ := u.GetFromCtx()
		username = user.Username
	}

	users := []*services.UserWithPhoto{}

	if relation == "following" {
		users = u.GetFollowing(username)
	}

	if relation == "followers" {
		users = u.GetFollowers(username)
	}

	component := components.UsersModal(users)
	return &component
}

func (u *UsersModal) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	component := u.GetComponent()
	(*component).Render(r.Context(), w)
}
