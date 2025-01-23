package handlers

import (
	"context"
	// "fmt"
	"github.com/a-h/templ"
	"github.com/jmarren/deepfried/components"
	"github.com/jmarren/deepfried/services"
	"net/http"
)

// Type
type CarouselSection struct {
	playableService PlayableService
	ctx             context.Context
	title           string
}

// Constructor
func NewCarousel(u PlayableService, ctx context.Context, title string) *CarouselSection {
	cs := new(CarouselSection)
	cs.playableService = services.NewPlayableService(ctx)
	cs.ctx = ctx
	cs.title = title
	return cs
}

// Private Methods

// Public Methods
func (cs *CarouselSection) GetComponent() *templ.Component {
	var data []*services.PlayableElt
	var component templ.Component
	switch cs.title {
	case "Editors Picks":
		data = cs.playableService.GetEditorsPicks()
		component = components.CarouselSectionBody(data)
	case "Just Added":
		data = cs.playableService.GetJustAdded()
		component = components.JustAddedSection(data)
	case "Poppin' Off":
		data = cs.playableService.GetPoppinOff()
		component = components.CarouselSectionBody(data)
	}
	return &component
}

func (e *CarouselSection) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	component := *e.GetComponent()
	component.Render(r.Context(), w)
}
