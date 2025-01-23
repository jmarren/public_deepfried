package handlers

import (
	"github.com/a-h/templ"
	"github.com/jmarren/deepfried/components"
	"github.com/jmarren/deepfried/services"
	"net/http"
)

type MyUploads struct {
	JHandle
}

func NewMyUploads(j JHandle) *MyUploads {
	m := new(MyUploads)
	m.JHandle = j
	return m
}

func (m *MyUploads) GetComponent() *templ.Component {
	user, _ := m.GetFromCtx()
	playableService := services.NewPlayableService(m.Context())
	data := playableService.GetUserUploads(user.ID)
	component := components.TrackCardSearch("Uploads", data)
	return &component
}

func (m *MyUploads) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	component := m.GetComponent()
	(*component).Render(r.Context(), w)
}
