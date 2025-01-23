package handlers

import (
	"github.com/a-h/templ"
	"github.com/jmarren/deepfried/components"
	"github.com/jmarren/deepfried/services"
	"net/http"
)

// Type
type Explore struct {
	JHandle
}

// Constructor
func NewExplore(j JHandle) *Explore {
	e := new(Explore)
	e.JHandle = j
	return e
}

// Public Methods
func (e *Explore) GetComponent() *templ.Component {
	// get user
	user, auth := e.UserService.GetFromCtx()

	// get tagbar
	tagService := services.NewTagService(e.Context())
	tagBar := NewTagBar(e.Context(), "", tagService).GetComponent()

	playableService := services.NewPlayableService(e.Context())

	// get the featured section component
	featuredSection := NewFeatured(e.Context(), playableService)
	featuredSectionComponent := featuredSection.GetComponent()

	editorsPickSection := NewCarousel(playableService, e.Context(), "Editors Picks").GetComponent()

	justAddedSection := NewCarousel(playableService, e.Context(), "Just Added").GetComponent()

	poppinOffSection := NewCarousel(playableService, e.Context(), "Poppin' Off").GetComponent()

	// get full component
	component := components.Explore(user, auth, *tagBar, *featuredSectionComponent, *justAddedSection, *editorsPickSection, *poppinOffSection)
	return &component
}

func (e *Explore) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("HX-Trigger", "init-explore-page")
	component := e.GetComponent()
	(*component).Render(r.Context(), w)
}
