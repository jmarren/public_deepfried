package handlers

import (
	"context"
	"fmt"
	// "github.com/jmarren/deepfried/services"
	"net/http"
)

type PostHandler struct {
	jwtVerified bool
	cognitoId   string
	BaseServices
	ctx context.Context
}

func NewPostHandler(ctx context.Context, jwtVerified bool, cognitoId string, b BaseServices) *PostHandler {
	p := new(PostHandler)
	p.ctx = ctx
	p.jwtVerified = jwtVerified
	p.BaseServices = b
	p.cognitoId = cognitoId
	return p

}

func (p *PostHandler) Restricted(auth bool, h http.Handler) http.Handler {
	if auth {
		return h
	} else {
		fmt.Println("user not authenticated")
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			w.WriteHeader(http.StatusUnauthorized)
		})
	}
}

func (p *PostHandler) ValidJwtOnly(h http.Handler) http.Handler {
	if !p.jwtVerified {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			w.WriteHeader(http.StatusUnauthorized)
		})
	}
	return h
}

func (p *PostHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	mux := http.NewServeMux()

	jHandle := JHandle{
		r,
		p.UserService,
		p.LogService,
	}

	mux.Handle("POST /downloads/{id}", NewDownloadHandler(jHandle))

	/* One route for valid jwt holders to create an account */
	mux.Handle("POST /users", p.ValidJwtOnly(NewCreateAccountHandler(p.ctx, p.cognitoId, p.UserService)))

	/* All others must exists in db */
	routeMap := make(map[string]http.Handler)
	routeMap["POST /upload"] = NewPostUpload(p.ctx, p.UserService)
	routeMap["POST /audio"] = NewPostAudioForm(p.ctx, p.UserService)
	routeMap["PATCH /users"] = NewProfileEdit(jHandle)
	routeMap["/following/{username}"] = NewFollowingHandler(jHandle)
	routeMap["PATCH /notifications"] = NewNotifications(jHandle)

	_, auth := p.UserService.GetFromCtx()
	for route, handler := range routeMap {
		mux.Handle(route, p.Restricted(auth, handler))
	}
	mux.ServeHTTP(w, r)
}
