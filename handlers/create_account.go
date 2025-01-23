package handlers

import (
	"bytes"
	"context"
	"fmt"
	"github.com/jmarren/deepfried/awssdk"
	"github.com/jmarren/deepfried/db"
	"github.com/jmarren/deepfried/sqlc"
	"github.com/jmarren/deepfried/util"
	"io"
	"log"
	"net/http"
)

type CreateAccountHandler struct {
	ctx         context.Context
	userService UserService
	cognitoId   string
}

func NewCreateAccountHandler(ctx context.Context, cognitoId string, u UserService) *CreateAccountHandler {
	c := new(CreateAccountHandler)
	c.ctx = ctx
	c.userService = u
	c.cognitoId = cognitoId
	return c
}

func (c *CreateAccountHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	username := r.FormValue("username")

	if c.cognitoId == "" {
		log.Println("error: cognitoId null during new user insertion")
		w.WriteHeader(http.StatusInternalServerError)
	}

	params := sqlc.CreateUserParams{
		CognitoID: c.cognitoId,
		Username:  username,
	}
	userId, err := db.Query.CreateUser(r.Context(), params)
	util.EMsg(err, "creating user")

	f, _, err := r.FormFile("profile-photo-input")
	if err != nil {
		log.Printf("error: %s\n", err)
	}

	fi, err := io.ReadAll(f)
	if err != nil {
		fmt.Println(err)
	}

	fileType := http.DetectContentType(fi)

	outputBuf := new(bytes.Buffer)

	outputBuf.Write(fi)

	destination := fmt.Sprintf("profile_photos/%s", util.UuidString(userId))

	err = awssdk.UploadBufferToS3(outputBuf, destination, fileType)
	if err != nil {
		fmt.Println(err)
	}

	if err == nil {
		log.Printf("username: %s", username)
	} else {
		fmt.Printf("error creating user: %s\n", err)
		code, ok := db.ErrorCode(err)
		if !ok {
			log.Print(err)
			w.WriteHeader(http.StatusInternalServerError)
			return
		}
		if code == "23505" {
			w.WriteHeader(http.StatusConflict)
			return
		} else {
			w.WriteHeader(http.StatusInternalServerError)
			return
		}
	}
	w.Header().Set("HX-Redirect", "/")
}
