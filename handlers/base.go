package handlers

import (
	// "context"
	"fmt"
	"net/http"

	"github.com/a-h/templ"
	"github.com/jmarren/deepfried/components"
	"github.com/jmarren/deepfried/services"
	"github.com/jmarren/deepfried/util"
	// "github.com/jmarren/deepfried/services"
)

type Base struct {
	JHandle
	componentHandler ComponentHandler
}

func NewBase(j JHandle, ch ComponentHandler) *Base {
	b := new(Base)
	b.JHandle = j
	b.componentHandler = ch
	return b
}

func (b *Base) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	b.Log(b.Context(), "base")
	modal := r.URL.Query().Get("modal")
	ctx := r.Context()
	switch modal {
	case "upload":
		ctx = templ.WithChildren(ctx, *NewUpload().GetComponent())
	case "create-account":
		ctx = templ.WithChildren(ctx, *NewCreateAccountModal().GetComponent())
	case "filters":
		ctx = templ.WithChildren(ctx, *NewFilters().GetComponent())
	}

	if r.Header.Get("HX-Request") == "true" {
		b.componentHandler.ServeHTTP(w, r)
	} else {
		page := b.componentHandler.GetComponent()
		user, auth := b.GetFromCtx()
		profilePhotoSrc := ""
		var notifications []*services.FollowNotification
		if !auth {
			profilePhotoSrc = ""
			notifications = []*services.FollowNotification{}
		} else {
			profilePhotoSrc = user.GetProfilePhoto()
			notifications = b.GetNotifications()
		}
		if page == nil {
			fmt.Println("nil page")
		}
		component := components.Base(profilePhotoSrc, notifications, *page)
		if component == nil {
			fmt.Println("base is nil")
		}
		err := component.Render(ctx, w)
		util.EMsg(err, "rendering base")
	}

}
