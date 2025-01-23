package services

import (
	"bytes"
	"context"
	"fmt"

	// "log"
	"os"
	"strings"

	// "github.com/jackc/pgx/v5/pgtype"
	"github.com/jmarren/deepfried/db"
	"github.com/jmarren/deepfried/sqlc"
	"github.com/jmarren/deepfried/util"
	// "github.com/jmarren/deepfried/util"
)

type UserService struct {
	ctx context.Context
}

func NewUserService(ctx context.Context) *UserService {
	u := new(UserService)
	u.ctx = ctx
	return u
}

type user_ctx_key_type string

const user_ctx_key user_ctx_key_type = "user"

func (u *UserService) VerifyHeader(authHeader string) (bool, string) {
	if os.Getenv("env") == "dev" {
		if os.Getenv("auth") == "cognito" {
			fmt.Println("dev environment: cognito only")
			return true, "00008"

		} else if os.Getenv("auth") == "true" {
			fmt.Println("dev environment: authenticated")
			return true, "00004"
		}
	}

	if authHeader == "" {
		return false, ""
	}
	amznData := strings.TrimSpace(authHeader)
	payload, err := ParseAndVerifyJWT(amznData)
	if err != nil {
		fmt.Printf("error: invalid jwt: %s\n", err)
		return false, ""
	}
	return true, payload.Sub
}

func (u *UserService) SetUserInCtx(cognitoId string) context.Context {
	res, err := db.Query.GetUserWithCognitoId(u.ctx, cognitoId)
	if err != nil {
		fmt.Printf("error getting user with cogId: %s\n", err)
		return u.ctx
	}
	// fmt.Printf("get user w cogId: %v\n", res)

	user := &User{
		&res,
	}
	newCtx := context.WithValue(u.ctx, user_ctx_key, user)
	u.ctx = newCtx
	return newCtx
}

func (u *UserService) UpdateUsernameInCtx(newUsername string) context.Context {
	current, _ := u.GetFromCtx()
	current.Username = newUsername
	// fmt.Printf("new context user: %v\n", current.Username)
	newCtx := context.WithValue(u.ctx, user_ctx_key, current)
	return newCtx
}

func (u *UserService) GetFromCtx() (*User, bool) {
	user, ok := u.ctx.Value(user_ctx_key).(*User)
	// fmt.Printf("user from ctx: %v\n", user)
	if !ok {
		fmt.Println("user not found in ctx")
		return nil, false
	}
	return user, true
}

func (u *UserService) IsUserAdmin() bool {
	user, auth := u.GetFromCtx()
	if !auth {
		return false
	}

	isAdmin, err := db.Query.IsUserAdmin(u.ctx, user.ID)
	util.EMsg(err, "querying for IsUserAdmin")
	if err != nil {
		return false
	}
	return isAdmin
}

func (u *UserService) GetProfilePhotoFromCtx() string {
	user, _ := u.GetFromCtx()
	if user == nil {
		return ""
	}

	return getProfileSrc(user.ID)
}

func (u *UserService) GetUserByUsername(username string) *User {
	res, err := db.Query.GetUserWithUsername(u.ctx, username)
	if err != nil {
		fmt.Printf("err getting user by username '%s': %s\n", username, err)
		return nil
	}

	return &User{
		&res,
	}
}

func (u *UserService) GetUserProfile(username string) *UserProfile {
	isMine := false
	userInCtx, found := u.GetFromCtx()

	if !found {
		isMine = false
	} else {
		if userInCtx.Username == username {
			isMine = true
		} else {
			isMine = false
		}
	}
	// if userInCtx == nil {
	// 	isMine = false
	// } else if userInCtx.User == nil {
	// 	isMine = false
	// } else {
	// 	if username == userInCtx.Username {
	// 		isMine = true
	// 	}
	// }
	user := u.GetUserByUsername(username)
	profile := NewProfile(u.ctx, user)
	profile.IsMine = isMine
	return profile
}

func (u *UserService) UpdateProfile(newUsername string, newBio string, newProfilePhoto *bytes.Buffer, filetype string) error {
	existingUser := u.GetUserByUsername(newUsername)
	if existingUser != nil {
		return fmt.Errorf("username already exists")
	}
	user, _ := u.GetFromCtx()
	if newUsername != "" {
		u.UpdateUsernameInCtx(newUsername)
	}

	err := user.UpdateProfile(u.ctx, newUsername, newBio, newProfilePhoto, filetype)
	return err
}

func (u *UserService) FollowUser(username string) error {
	user, auth := u.GetFromCtx()
	if !auth {
		return fmt.Errorf("user not authentciated")
	}
	myId := user.ID
	err := db.Query.FollowUsername(u.ctx, sqlc.FollowUsernameParams{
		TheirUsername: username,
		MyID:          myId,
	})

	util.EMsg(err, "following user")
	if err == nil {
		err = db.Query.AddFollowNotification(u.ctx, sqlc.AddFollowNotificationParams{
			MyID:          myId,
			TheirUsername: username,
		})
		util.EMsg(err, "adding follow notification")
	}

	return err
}

func (u *UserService) MarkFollowNotficationSeen() {
	user, auth := u.GetFromCtx()
	if !auth {
		return
	}
	myId := user.ID
	err := db.Query.MarkAllFollowNotificationSeen(u.ctx, myId)
	util.EMsg(err, "marking notification as seen")

}

// type UnFollowUsernameParams struct {
// 	TheirUsername string
// 	MyID          pgtype.UUID
// }

func (u *UserService) UnFollowUser(username string) error {
	user, auth := u.GetFromCtx()
	if !auth {
		return fmt.Errorf("user not authenticated")
	}
	myId := user.ID
	err := db.Query.UnFollowUsername(u.ctx, sqlc.UnFollowUsernameParams{
		TheirUsername: username,
		MyID:          myId,
	})
	util.EMsg(err, "following user")
	return err
}

type FollowNotification struct {
	*sqlc.User
	ProfilePhotoSrc string
}

func (u *UserService) GetNotifications() []*FollowNotification {
	user, auth := u.GetFromCtx()
	if !auth {
		return nil
	}
	myId := user.ID
	res, err := db.Query.GetFollowNotifications(u.ctx, myId)
	var followNotifications []*FollowNotification

	for _, item := range res {
		followNotifications = append(followNotifications, &FollowNotification{
			&item,
			item.GetProfilePhotoSrc(),
		})
	}

	util.EMsg(err, "getting follow notifications")
	return followNotifications
}

type UserWithPhoto struct {
	*sqlc.User
	ProfilePhotoSrc string
	IAmFollowing    bool
}

// func (u *User) AmIFollowing(ctx context.Context, myId pgtype.UUID) bool {

func (u *UserService) GetFollowers(username string) []*UserWithPhoto {
	res, err := db.Query.GetFollowers(u.ctx, username)

	user, auth := u.GetFromCtx()
	if !auth {
		return nil
	}
	myId := user.ID

	var users []*UserWithPhoto
	for _, item := range res {
		user := &User{&item}
		users = append(users, &UserWithPhoto{
			&item,
			item.GetProfilePhotoSrc(),
			user.AmIFollowing(u.ctx, myId),
		})
	}
	util.EMsg(err, "getting followers")
	return users
}

func (u *UserService) GetFollowing(username string) []*UserWithPhoto {
	res, err := db.Query.GetFollowing(u.ctx, username)
	var users []*UserWithPhoto
	user, auth := u.GetFromCtx()
	if !auth {
		return nil
	}
	myId := user.ID
	for _, item := range res {
		user := &User{&item}
		users = append(users, &UserWithPhoto{
			&item,
			item.GetProfilePhotoSrc(),
			user.AmIFollowing(u.ctx, myId),
		})
	}
	util.EMsg(err, "getting following")
	return users
}

type Pin struct {
	*PlayableElt
	Tags []string
}

// User Profile
type UserProfile struct {
	*User
	IAmFollowing    bool
	IsMine          bool
	ProfilePhotoSrc string
	*sqlc.GetUserInfoRow
	Pins   []*Pin
	Tracks []*PlayableElt
}

func NewProfile(ctx context.Context, u *User) *UserProfile {
	user, auth := NewUserService(ctx).GetFromCtx()

	p := new(UserProfile)
	p.User = u
	p.Pins = u.GetPins(ctx)
	p.Tracks = u.GetMostPopularAudio(ctx)
	p.GetUserInfoRow = u.GetUserInfo(ctx)
	p.ProfilePhotoSrc = u.GetProfilePhoto()
	if !auth {
		p.IAmFollowing = false
	} else {
		p.IAmFollowing = u.AmIFollowing(ctx, user.ID)
	}
	return p
}
