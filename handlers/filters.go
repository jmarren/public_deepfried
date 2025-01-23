package handlers

import (
	"github.com/a-h/templ"
	"github.com/jmarren/deepfried/components"
	"net/http"
)

type Filters struct {
}

func NewFilters() *Filters {
	f := new(Filters)
	return f
}

func (f *Filters) GetComponent() *templ.Component {
	component := components.FiltersModal()
	return &component
}

func (f *Filters) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	component := f.GetComponent()
	(*component).Render(r.Context(), w)
}
