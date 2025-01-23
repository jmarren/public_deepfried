package handlers

import (
	"github.com/a-h/templ"
	"github.com/jmarren/deepfried/components"
	"net/http"
)

type AccountDropdown struct{}

func (a *AccountDropdown) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	component := a.GetComponent()
	(*component).Render(r.Context(), w)
}

func (a *AccountDropdown) GetComponent() *templ.Component {
	component := components.AccountDropdown()
	return &component
}
