package handlers

import (
	"net/http"

	"github.com/a-h/templ"
	"github.com/jmarren/deepfried/components"
	"github.com/jmarren/deepfried/util"
	// // "github.com/jmarren/deepfried/services"
)

type EditProfileModal struct {
	JHandle
}

func NewEditProfileModal(j JHandle) *EditProfileModal {
	e := new(EditProfileModal)
	e.JHandle = j
	return e
}

func (e *EditProfileModal) GetComponent() *templ.Component {
	user, auth := e.GetFromCtx()
	if !auth {
		return nil
	}
	profilePhotoSrc := user.GetProfilePhoto()
	bio := user.GetBio(e.Context())
	component := components.EditProfile(user, profilePhotoSrc, bio)
	return &component
}

func (e *EditProfileModal) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	currentUrl := r.Header.Get("HX-Current-Url")
	newUrl, err := util.UpdateQueryParam(currentUrl, "modal", "edit-profile")
	util.EMsg(err, "updating query params for edit-profile modal")
	w.Header().Set("HX-Push-Url", newUrl)
	(*e.GetComponent()).Render(r.Context(), w)
}
