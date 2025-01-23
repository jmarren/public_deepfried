package handlers

import (
	"bytes"
	"context"

	"github.com/a-h/templ"
	"github.com/jmarren/deepfried/services"
	"github.com/jmarren/deepfried/sqlc"

	"github.com/jackc/pgx/v5/pgtype"
	"net/http"
)

type UserService interface {
	SetUserInCtx(cognitoId string) context.Context
	GetFromCtx() (*services.User, bool)
	GetProfilePhotoFromCtx() string
	GetUserProfile(username string) *services.UserProfile
	GetUserByUsername(username string) *services.User
	UpdateProfile(newUsername string, newBio string, newProfilePhoto *bytes.Buffer, filetype string) error
	UpdateUsernameInCtx(newUsername string) context.Context
	VerifyHeader(authHeader string) (bool, string)
	GetNotifications() []*services.FollowNotification
	MarkFollowNotficationSeen()
	GetFollowing(username string) []*services.UserWithPhoto
	GetFollowers(username string) []*services.UserWithPhoto
	IsUserAdmin() bool
}

type BaseRenderFunc func(*templ.Component) templ.Component

type ComponentHandler interface {
	http.Handler
	GetComponent() *templ.Component
}

type PlayableService interface {
	GetFeatured() *services.PlayableElt
	GetEditorsPicks() []*services.PlayableElt
	GetJustAdded() []*services.PlayableElt
	GetUserDownloads(userId pgtype.UUID, keyword string) []*services.PlayableElt
	GetUserUploads(userId pgtype.UUID) []*services.PlayableElt
	GetUserFeed(userId pgtype.UUID) []*services.UserFeedItem
	GetTrackPage(username string, title string) *services.TrackPage
	GetPoppinOff() []*services.PlayableElt
}

// func (s *SearchService) GetSearchResults(ctx context.Context, keyword string) *AudioSearchResults {
type SearchService interface {
	GetAudioSearchResults(ctx context.Context, filtersStrs *services.SearchAudioFilters) *services.AudioSearchResults
	GetUserSearchResults(ctx context.Context, keyword *string) *services.UserSearchResults
	GetDropdownItems(ctx context.Context, keyword *string) []*sqlc.SearchKeywordForDropdownRow
}

type TagService interface {
	GetMostPopular() []string
}

type UploadService interface {
	NewUpload(filename string, totalSize int, uploadType string)
	Append(filename string, uploadType string, data []byte, index int) error
}

type LogService interface {
	Log(ctx context.Context, msg string)
	WithId(ctx context.Context) context.Context
	GetReqId(ctx context.Context) (int, error)
	LogId(ctx context.Context)
}

type BaseServices struct {
	UserService
	LogService
}

type JHandle struct {
	*http.Request
	UserService
	LogService
}

func (b *BaseServices) LogC() {

}
