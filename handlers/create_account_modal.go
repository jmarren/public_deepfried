package handlers

import (
	"github.com/a-h/templ"
	"github.com/jmarren/deepfried/components"
	"net/http"
	// // "github.com/jmarren/deepfried/services"
)

type CreateAccountModal struct{}

func NewCreateAccountModal() *CreateAccountModal {
	return &CreateAccountModal{}
}

func (c *CreateAccountModal) GetComponent() *templ.Component {
	component := components.CreateAccount()
	return &component
}

func (c *CreateAccountModal) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("HX-Retarget", "#modal")
	w.Header().Set("HX-Reswap", "innerHTML")
	(*c.GetComponent()).Render(r.Context(), w)
}
