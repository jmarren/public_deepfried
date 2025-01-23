package handlers

import (
	// "fmt"
	"github.com/a-h/templ"
	"github.com/jmarren/deepfried/components"
	"github.com/jmarren/deepfried/services"
	"net/http"
)

type MyDownloads struct {
	JHandle
}

func NewMyDownloads(j JHandle) *MyDownloads {
	m := new(MyDownloads)
	m.JHandle = j
	return m
}

func (m *MyDownloads) GetComponent() *templ.Component {
	user, auth := m.GetFromCtx()

	if !auth {
		return nil
	}

	keyword := m.FormValue("downloads_keyword")

	playableService := services.NewPlayableService(m.Context())
	data := playableService.GetUserDownloads(user.ID, keyword)
	component := components.TrackCardSearch("Downloads", data)
	return &component
}

func (m *MyDownloads) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	component := m.GetComponent()
	(*component).Render(r.Context(), w)
}
