package handlers

import (
	// "context"
	"github.com/a-h/templ"
	"github.com/jmarren/deepfried/components"
	// "github.com/jmarren/deepfried/services"
	"net/http"
)

type Hot struct {
	Component *templ.Component
}

func (h *Hot) SetComponent() {
	component := components.Hot()
	h.Component = &component
}

func (h *Hot) GetComponent() *templ.Component {
	return h.Component
}

func NewHot() *Hot {
	h := new(Hot)
	h.SetComponent()
	return h
}

func (h *Hot) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	(*h.Component).Render(r.Context(), w)
}
