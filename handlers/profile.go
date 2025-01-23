package handlers

import (
	"fmt"
	"net/http"

	"github.com/a-h/templ"
	"github.com/jmarren/deepfried/components"
	"github.com/jmarren/deepfried/services"
)

type Profile struct {
	*http.Request
	user *services.User
}

func NewProfile(r *http.Request, user *services.User) *Profile {
	p := new(Profile)
	p.Request = r
	p.user = user
	return p
}

func (p *Profile) GetComponent() *templ.Component {
	username := p.user.Username

	userService := services.NewUserService(p.Context())
	_, auth := userService.GetFromCtx()
	fmt.Printf("profile: authenticated?: %t\n", auth)

	data := userService.GetUserProfile(username)
	fmt.Printf("data.IsMine: %t\n", data.IsMine)

	pins := components.Pins(data.Pins)

	mostPopular := components.CarouselSectionBody(data.Tracks)

	component := components.Profile(auth, data, pins, mostPopular, "somethin")
	return &component
}

func (p *Profile) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("HX-Push-Url", fmt.Sprintf("/%s", p.user.Username))
	(*p.GetComponent()).Render(r.Context(), w)
}
