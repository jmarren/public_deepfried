package handlers

import (
	"fmt"
	"strings"
	// "github.com/a-h/templ"
	// "github.com/jmarren/deepfried/components"
	"net/http"

	// "github.com/a-h/templ"
	"github.com/jmarren/deepfried/services"
)

type ProfileRoute struct {
	JHandle
}

func NewProfileRoute(j JHandle) *ProfileRoute {
	p := new(ProfileRoute)
	p.JHandle = j
	return p
}

func (p *ProfileRoute) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	ctx := p.Context()
	path := p.URL.Path
	path = strings.Trim(path, "/")
	pathVals := strings.Split(path, "/")
	username := pathVals[0]

	var userProfileRequested *services.User
	var auth bool

	if username == "me" {
		userProfileRequested, auth = p.GetFromCtx()
		if !auth {
			w.WriteHeader(http.StatusUnauthorized)
		}
	} else {
		userProfileRequested = p.GetUserByUsername(username)
	}

	if userProfileRequested == nil {
		fmt.Printf("profile requested is nil")
		w.Header().Set("HX-Retarget", "#page-content")
		w.Header().Set("HX-Reswap", "innerHTML")
		NewBase(p.JHandle, NewExplore(p.JHandle)).ServeHTTP(w, p.Request)
	} else {
		// track := pathVals[1]
		track := ""
		if len(pathVals) > 1 {
			track = pathVals[1]
		}
		// track := r.PathValue("track")
		if track == "" {
			profile := NewProfile(r, userProfileRequested)
			NewBase(p.JHandle, profile).ServeHTTP(w, r)
		} else {
			playableService := services.NewPlayableService(ctx)
			trackPage := NewTrackPage(ctx, playableService, userProfileRequested.Username, track) // Need to handle user exists but track doesnt
			NewBase(p.JHandle, trackPage).ServeHTTP(w, r)
		}
	}
}
