package handlers

import (
	"fmt"
	"github.com/a-h/templ"
	"github.com/google/uuid"
	"net/http"
	// "github.com/jmarren/deepfried/cache"
	"github.com/jmarren/deepfried/components"
	"github.com/jmarren/deepfried/util"
)

type Upload struct{}

func NewUpload() *Upload {
	u := new(Upload)
	return u
}

func (u *Upload) GetComponent() *templ.Component {
	uploadId := uuid.New().String()
	fmt.Println("NEW UPLOAD")
	fmt.Printf("uploadId: %s\n", uploadId)
	component := components.Upload(uploadId)
	return &component
}

func (u *Upload) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	currentUrl := r.Header.Get("HX-Current-Url")
	newUrl, err := util.UpdateQueryParam(currentUrl, "modal", "upload")
	util.EMsg(err, "adding ?modal=upload to url")
	w.Header().Set("HX-Push-Url", newUrl)

	component := u.GetComponent()
	(*component).Render(r.Context(), w)
}
