package services

import (
	"bytes"
	"context"
	"fmt"
	"log"
	// "os"
	// "strings"

	"github.com/jackc/pgx/v5/pgtype"
	"github.com/jmarren/deepfried/awssdk"
	"github.com/jmarren/deepfried/db"
	"github.com/jmarren/deepfried/sqlc"
	"github.com/jmarren/deepfried/util"
)

type User struct {
	*sqlc.User
}

func NewUser() *User {
	u := new(User)
	return u
}

func (u *User) GetPins(ctx context.Context) []*Pin {
	fmt.Printf("pins -- username: %s\n\tuserId: %s\n", u.Username, util.UuidString(u.ID))
	res, err := db.Query.GetUserPins(ctx, u.ID)
	if err != nil {
		fmt.Printf("err getting user pins: %s\n", err)
		return nil
	}
	var pins []*Pin
	for _, pin := range res {
		playable := NewPlayableElt(&pin.Playable, ctx)
		pins = append(pins,
			&Pin{
				playable,
				pin.TagArray,
			})
	}
	return pins
}

func (u *User) GetMostPopularAudio(ctx context.Context) []*PlayableElt {
	res, err := db.Query.GetUserAudioFiles(ctx, u.ID)
	util.EMsg(err, "getting user audio files")

	var playables []*PlayableElt
	for _, row := range res {
		playable := NewPlayableElt(&row, ctx)
		playables = append(playables, playable)
	}
	return playables
}

func (u *User) GetUserInfo(ctx context.Context) *sqlc.GetUserInfoRow {
	res, err := db.Query.GetUserInfo(ctx, u.Username)
	util.EMsg(err, "getting user info")
	return &res
}

func (u *User) GetProfilePhoto() string {
	if u.ID.Valid != true {
		return ""
	}
	return getProfileSrc(u.ID)
}

func (u *User) GetBio(ctx context.Context) string {
	bio, err := db.Query.GetUserBio(ctx, u.ID)
	util.EMsg(err, "getting user bio")
	return bio.String
}

func (u *User) UpdateProfile(ctx context.Context, newUsername string, newBio string, newProfilePhoto *bytes.Buffer, filetype string) error {
	fmt.Printf("newUsername: %s\nnewBio:%s\nnewProfilePhoto.Len(): %d\n", newUsername, newBio, newProfilePhoto.Len())
	if newUsername != "" {
		err := db.Query.UpdateUserUsername(ctx, sqlc.UpdateUserUsernameParams{
			ID:       u.ID,
			Username: newUsername,
		})
		if err != nil {
			return fmt.Errorf("an error occurred. Please try again later.")
		}
	}

	if newBio != "" {
		err := db.Query.UpdateUserBio(ctx, sqlc.UpdateUserBioParams{
			UserID: u.ID,
			Bio: pgtype.Text{
				String: newBio,
				Valid:  true,
			},
		})
		if err != nil {
			return fmt.Errorf("an error occurred. Please try again later.")
		}
	}

	if newProfilePhoto.Len() != 0 {
		destination := fmt.Sprintf("profile_photos/%s", util.UuidString(u.ID))

		err := awssdk.UploadBufferToS3(newProfilePhoto, destination, filetype)
		if err != nil {
			log.Printf("error uploading new profile photo file: %s\n", err)
			return fmt.Errorf("Please try a different image file")
		}
	}
	return nil
}

func (u *User) AmIFollowing(ctx context.Context, myId pgtype.UUID) bool {
	IAmFollowing, err := db.Query.GetAmIFollowing(ctx, sqlc.GetAmIFollowingParams{
		MyID:    myId,
		TheirID: u.ID,
	})
	util.EMsg(err, "getting 'am I following' from db")
	return IAmFollowing
}

// func (u *User) FollowUser(ctx context.Context,
