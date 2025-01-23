package handlers

import (
	"github.com/a-h/templ"
	"net/http"
)

type AuthRedirect struct {
	JHandle
}

func NewAuthRedirect(j JHandle) *AuthRedirect {
	a := new(AuthRedirect)
	a.JHandle = j
	return a
}

func (a *AuthRedirect) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("HX-Push-Url", "/")
	w.Header().Set("HX-Redirect", "/login")
	w.WriteHeader(http.StatusTemporaryRedirect)
}

func (a *AuthRedirect) GetComponent() *templ.Component {
	component := NewExplore(a.JHandle).GetComponent()
	return component
}
