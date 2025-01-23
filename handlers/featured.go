package handlers

import (
	"context"
	"github.com/a-h/templ"
	"github.com/jmarren/deepfried/components"
	"net/http"
)

type Featured struct {
	PlayableService PlayableService
	ctx             context.Context
}

func NewFeatured(ctx context.Context, p PlayableService) *Featured {
	f := new(Featured)
	f.ctx = ctx
	f.PlayableService = p
	return f
}

func (f *Featured) GetComponent() *templ.Component {
	playable := f.PlayableService.GetFeatured()
	component := components.FeaturedSectionBody(playable)
	return &component
}

func (f *Featured) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	component := f.GetComponent()
	(*component).Render(r.Context(), w)
}
