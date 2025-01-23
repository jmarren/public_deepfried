package handlers

import (
	"github.com/a-h/templ"
	"github.com/jmarren/deepfried/components"
	"github.com/jmarren/deepfried/services"
	"net/http"
)

// Type
type Feed struct {
	JHandle
}

// Constructor
func NewFeed(j JHandle) *Feed {
	e := new(Feed)
	e.JHandle = j
	// e.BaseServices = serv
	// e.Request = r
	return e
}

// Public Methods
func (e *Feed) GetComponent() *templ.Component {
	// get user feed
	user, _ := e.UserService.GetFromCtx()

	userId := user.ID
	playableService := services.NewPlayableService(e.Context())

	data := playableService.GetUserFeed(userId)

	component := components.UserFeed(data, "")
	return &component
}

func (e *Feed) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	component := e.GetComponent()
	(*component).Render(r.Context(), w)
}
