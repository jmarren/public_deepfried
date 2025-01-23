package handlers

import (
	"bytes"
	// "context"
	"fmt"
	"github.com/jmarren/deepfried/components"
	// "github.com/jmarren/deepfried/services"
	"io"
	"log"
	"net/http"
)

type EditProfile struct {
	JHandle
}

func NewProfileEdit(j JHandle) *EditProfile {
	e := new(EditProfile)
	e.JHandle = j
	return e
}

func (e *EditProfile) ServeError(w http.ResponseWriter, r *http.Request, msg string) {
	w.Header().Set("HX-Retarget", "#edit-profile-error")
	w.Header().Set("HX-Reswap", "outerHTML")
	components.EditProfileError(msg).Render(r.Context(), w)
}

func (e *EditProfile) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	newUsername := e.FormValue("username")
	newBio := e.FormValue("bio")
	filetype := ""
	f, _, err := e.FormFile("profile-photo-input")

	var newProfilePhoto = new(bytes.Buffer)

	if f != nil {
		fmt.Println("profile photo received")
		if err != nil {
			log.Printf("error reading uploaded file: %s\n", err)
			e.ServeError(w, e.Request, "please try a different image file")
			return
		}
		fi, err := io.ReadAll(f)
		if err != nil {
			e.ServeError(w, e.Request, "please try a different image file")
			return
		}

		filetype = http.DetectContentType(fi)

		newProfilePhoto.Write(fi)
	}

	err = e.UpdateProfile(newUsername, newBio, newProfilePhoto, filetype)

	if err != nil {
		e.ServeError(w, e.Request, fmt.Sprintf("%s", err))
		return
	}

	user, auth := e.GetFromCtx()
	if !auth {
		fmt.Println("user is nil in EditProfile")
	}
	fmt.Printf("======= user: %s\n", user.Username)

	if newUsername != "" {
		user.Username = newUsername
		newCtx := e.UpdateUsernameInCtx(newUsername)
		e.Request = e.WithContext(newCtx)
	}
	w.Header().Set("HX-Trigger", "exit-modal")

	NewProfile(e.Request, user).ServeHTTP(w, e.Request)
}
