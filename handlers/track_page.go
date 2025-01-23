package handlers

import (
	"context"
	"github.com/a-h/templ"
	"github.com/jmarren/deepfried/components"
	"net/http"
)

type TrackPage struct {
	ctx             context.Context
	playableService PlayableService
	username        string
	title           string
}

func NewTrackPage(ctx context.Context, playableService PlayableService, username string, title string) *TrackPage {
	t := new(TrackPage)
	t.ctx = ctx
	t.username = username
	t.title = title
	t.playableService = playableService
	return t
}

func (t *TrackPage) GetComponent() *templ.Component {
	data := t.playableService.GetTrackPage(t.username, t.title)
	component := components.TrackPage(data)
	return &component
}

func (t *TrackPage) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	component := t.GetComponent()
	(*component).Render(r.Context(), w)
}
