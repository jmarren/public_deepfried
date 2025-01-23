package handlers

import (
	"fmt"
	// "context"
	"github.com/a-h/templ"
	"github.com/jmarren/deepfried/components"
	"net/http"
)

type UploadForm struct {
	JHandle
}

func NewUploadForm(j JHandle) *UploadForm {
	u := new(UploadForm)
	u.JHandle = j
	// u.ctx = ctx
	// u.title = title
	return u
}

func (u *UploadForm) GetComponent() *templ.Component {
	title := u.Header.Get("X-file-name")
	uploadId := u.Header.Get("X-upload-id")
	fmt.Printf("UPLOAD_ID: %s\n", uploadId)
	component := components.UploadForm(title, uploadId)
	return &component
}

func (u *UploadForm) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	component := u.GetComponent()
	(*component).Render(r.Context(), w)
}
