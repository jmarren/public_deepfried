package handlers

import (
	"net/http"

	"github.com/a-h/templ"
	"github.com/jmarren/deepfried/components"
)

type AdminPage struct {
	JHandle
}

func NewAdminPage(j JHandle) *AdminPage {
	a := new(AdminPage)
	a.JHandle = j
	return a
}

func (a *AdminPage) GetComponent() *templ.Component {
	component := components.Admin()
	return &component
}

func (a *AdminPage) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	component := a.GetComponent()
	(*component).Render(r.Context(), w)
}
