package handlers

import (
	// "context"
	"fmt"
	"net/http"
	// "net/http/httputil"

	"github.com/jmarren/deepfried/services"
	// "github.com/jmarren/deepfried/util"
)

type DefaultHandler struct{}

func NewDefaultHandler() *DefaultHandler {
	d := new(DefaultHandler)
	return d
}

//	func Restricted(user *services.User, j JHandle, h ComponentHandler) ComponentHandler {
//		if user != nil {
//			return h
//		} else if g.jwtVerified {
//			return NewCreateAccountModal()
//		} else {
//			return NewAuthRedirect(j)
//		}
//	}
func (d *DefaultHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	fmt.Printf("\n\nNEW REQUEST: %s %s\n\n", r.Method, r.URL.Path)

	// create log service and add the generated request id to ctx
	logService := services.NewLogService()
	r = r.WithContext(logService.WithId(r.Context()))

	logService.Log(r.Context(), fmt.Sprintf("%s %s", r.Method, r.URL.Path))
	logService.LogId(r.Context())

	// create user service and verify jwt
	userService := services.NewUserService(r.Context())
	authHeader := r.Header.Get("X-Amzn-Oidc-Data")
	jwtVerified, cognitoId := userService.VerifyHeader(authHeader)

	// if verified, add user to ctx
	if jwtVerified {
		ctx := userService.SetUserInCtx(cognitoId)
		r = r.WithContext(ctx)
	}

	// package services
	baseServices := BaseServices{
		userService,
		logService,
	}

	// create server
	mux := http.NewServeMux()

	// get any currently playing song from the "X-Playing" header and add it to context
	currentlyPlaying := r.Header.Get("HX-Audio-Playing")

	playableService := services.NewPlayableService(r.Context())
	ctx := playableService.SetPlayingInCtx(currentlyPlaying)
	r = r.WithContext(ctx)

	jHandle := JHandle{
		r,
		userService,
		logService,
	}

	_, auth := userService.GetFromCtx()
	fmt.Printf("authenticated?: %t\n", auth)

	fmt.Printf("jwt verified: %t\n", jwtVerified)

	// handler routes by method
	mux.Handle("GET /", NewGetRouter(jwtVerified, jHandle))
	mux.Handle(" /", NewPostHandler(r.Context(), jwtVerified, cognitoId, baseServices))
	// mux.Handle("PATCH /", NewPostHandler(r.Context(), jwtVerified, cognitoId, baseServices))
	mux.ServeHTTP(w, r)
}
