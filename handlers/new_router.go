package handlers

import (
	"fmt"
	"net/http"
	// "strings"
	// "net/url"
	// "github.com/jmarren/deepfried/services"
)

type GetRouter struct {
	jwtVerified bool
	JHandle
}

func NewGetRouter(jwtVerified bool, j JHandle) *GetRouter {
	g := new(GetRouter)
	g.jwtVerified = jwtVerified
	g.JHandle = j
	return g
}

func (g *GetRouter) Restricted(j JHandle, h ComponentHandler) ComponentHandler {
	_, auth := j.GetFromCtx()
	fmt.Printf("getRouter: authenticated?: %t\n", auth)
	if auth {
		return h
	} else if g.jwtVerified {
		return NewCreateAccountModal()
	} else {
		return NewAuthRedirect(j)
	}
}

func (g *GetRouter) AdminOnly(j JHandle, h ComponentHandler) ComponentHandler {
	if j.IsUserAdmin() {
		return h
	} else {
		return NewExplore(j)
	}
}

func (g *GetRouter) ServeHTTP(w http.ResponseWriter, r *http.Request) {

	jHandle := g.JHandle
	isPage := false
	isRestricted := true
	// adminOny := false
	var ch ComponentHandler

	switch r.URL.Path {
	case "/player":
		isPage = false
		isRestricted = false
		ch = NewPlayer(jHandle)
	case "/modal/upload":
		isRestricted = true
		isPage = false
		ch = NewUpload()
	case "/modal/create-account":
		isPage = false
		if g.jwtVerified {
			ch = NewCreateAccountModal()
		}
	case "/my-downloads":
		isRestricted = true
		isPage = true
		ch = NewMyDownloads(jHandle)
	case "/my-uploads":
		isRestricted = true
		isPage = true
		ch = NewMyUploads(jHandle)
	case "/feed":
		isRestricted = true
		isPage = true
		ch = NewFeed(jHandle)
	case "/modal/edit-profile":
		isRestricted = true
		isPage = false
		ch = NewEditProfileModal(jHandle)
	case "/modal/upload-form":
		isRestricted = true
		isPage = false
		ch = NewUploadForm(jHandle)
	case "/":
		isRestricted = false
		isPage = true
		ch = NewExplore(jHandle)
	case "/explore":
		isRestricted = false
		isPage = true
		ch = NewExplore(jHandle)
	case "/search":
		isRestricted = false
		isPage = true
		ch = NewSearch(jHandle)
	case "/search-bar-dropdown":
		isRestricted = false
		isPage = false
		ch = NewSearchDropdown(jHandle)
	case "/modal/filters":
		isRestricted = false
		isPage = false
		ch = NewFilters()
	case "/x":
		isPage = true
		isRestricted = true
		ch = g.AdminOnly(jHandle, NewAdminPage(jHandle))
	}

	if ch == nil {
		NewProfileRoute(jHandle).ServeHTTP(w, r)
		return
	}

	if isRestricted {
		ch = g.Restricted(jHandle, ch)
	}
	if isPage {
		NewBase(jHandle, ch).ServeHTTP(w, r)
		return
	}

	if ch != nil {
		ch.ServeHTTP(w, r)
		return
	}

	w.WriteHeader(http.StatusNotFound)
}
