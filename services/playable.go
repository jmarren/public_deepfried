package services

import (
	"context"
	"fmt"
	"github.com/google/uuid"

	"github.com/jackc/pgx/v5/pgtype"
	"github.com/jmarren/deepfried/db"
	"github.com/jmarren/deepfried/sqlc"
	"github.com/jmarren/deepfried/util"
)

// type user_ctx_key_type string
//
// const user_ctx_key user_ctx_key_type = "user"
//

type playing_ctx_type string

const playing_ctx_key playing_ctx_type = "currently_playing"

type PlayableService struct {
	ctx context.Context
}

func NewPlayableService(ctx context.Context) *PlayableService {
	p := new(PlayableService)
	p.ctx = ctx
	return p
}

func (p *PlayableService) SetPlayingInCtx(playing string) context.Context {
	p.ctx = context.WithValue(p.ctx, playing_ctx_key, playing)
	return p.ctx
}

func (p *PlayableService) GetPlayingFromCtx() string {
	res := p.ctx.Value(playing_ctx_key)
	playing, ok := res.(string)
	if !ok {
		return ""
	}
	return playing
}

func (p *PlayableService) GetFeatured() *PlayableElt {
	dbData, err := db.Query.GetFeaturedTrack(p.ctx)
	util.EMsg(err, "getting featured track")

	playable := NewPlayableElt(&dbData, p.ctx)

	tags, err := db.Query.GetAudioFileTags(p.ctx, playable.ID)
	util.EMsg(err, "getting tags for featured")
	playable.Tags = tags

	return playable
}

func (p *PlayableService) GetJustAdded() []*PlayableElt {
	dbData, err := db.Query.GetJustAdded(p.ctx)
	util.EMsg(err, "Getting just added")
	var playables []*PlayableElt
	for _, row := range dbData {
		playable := NewPlayableElt(&row, p.ctx)
		playables = append(playables, playable)
	}
	return playables

}

func (p *PlayableService) GetPoppinOff() []*PlayableElt {
	dbData, err := db.Query.GetMostPopularAudioFiles(p.ctx)
	util.EMsg(err, "Getting editors picks")
	var playables []*PlayableElt
	for _, row := range dbData {
		playable := NewPlayableElt(&row, p.ctx)
		playables = append(playables, playable)
	}
	return playables

}

func (p *PlayableService) GetEditorsPicks() []*PlayableElt {
	dbData, err := db.Query.GetEditorsPicks(p.ctx)
	util.EMsg(err, "Getting editors picks")
	var playables []*PlayableElt
	for _, row := range dbData {
		playable := NewPlayableElt(&row, p.ctx)
		playables = append(playables, playable)
	}
	return playables
}

func (p *PlayableService) GetUserDownloads(userId pgtype.UUID, keyword string) []*PlayableElt {
	var res []sqlc.Playable
	var err error
	if keyword == "" {
		res, err = db.Query.GetUserDownloads(p.ctx, userId)
	} else {
		fmt.Printf("searching with keyword: %s\n", keyword)
		res, err = db.Query.GetUserDownloadsWithKeyword(p.ctx, sqlc.GetUserDownloadsWithKeywordParams{
			UserID:  userId,
			Keyword: keyword,
		})
	}

	util.EMsg(err, "getting user downloads")
	var playables []*PlayableElt
	for _, row := range res {
		playable := NewPlayableElt(&row, p.ctx)
		playables = append(playables, playable)
	}
	return playables
}

// This is essentially repeated in the user service
// Remove one ?
func (p *PlayableService) GetUserUploads(userId pgtype.UUID) []*PlayableElt {
	res, err := db.Query.GetUserAudioFiles(p.ctx, userId)
	util.EMsg(err, "getting user audio files")

	return p.playablesToElts(res)
}

func (p *PlayableService) playablesToElts(in []sqlc.Playable) []*PlayableElt {
	var playables []*PlayableElt
	for _, row := range in {
		playable := NewPlayableElt(&row, p.ctx)
		playables = append(playables, playable)
	}
	return playables
}

type UserFeedItem struct {
	*PlayableElt
	VisArr      []int32
	UsageRights pgtype.Text
	TagArray    []string
}

func (p *PlayableService) GetUserFeed(userId pgtype.UUID) []*UserFeedItem {
	res, err := db.Query.GetUserFeed(p.ctx, userId)
	util.EMsg(err, "getting user feed items")

	var userFeedItems []*UserFeedItem
	for _, row := range res {
		playable := NewPlayableElt(&row.Playable, p.ctx)
		userFeedItems = append(userFeedItems, &UserFeedItem{
			playable,
			row.VisArr,
			row.UsageRights,
			row.TagArray,
		})
	}
	return userFeedItems
}

type PlayerData struct {
	Current *PlayableElt
	Queue   []*PlayableElt
}

func (p *PlayableService) GetPlayerData(playing string, audioQueue []string) *PlayerData {
	player := new(PlayerData)
	player.Current = p.GetAudioFileById(playing)
	for _, queueItem := range audioQueue {
		player.Queue = append(player.Queue, p.GetAudioFileById(queueItem))
	}

	return player
}

func (p *PlayableService) GetAudioFileById(id string) *PlayableElt {
	parsedId, err := uuid.Parse(id)
	util.EMsg(err, "parsing audio file id to uuid")
	pgId := pgtype.UUID{
		Bytes: parsedId,
		Valid: true,
	}
	res, err := db.Query.GetAudioFileById(p.ctx, pgId)
	util.EMsg(err, "getting audio file by Id")

	playableElt := NewPlayableElt(&res, p.ctx)

	return playableElt

}

// type XtraAudioInfo struct {
// 	Bpm          int32
// 	Key          enums.MusicalKey
// 	KeySignature enums.MusicalKeySignature
// 	IsMinor      bool
// 	UploadedTime string
// 	Tags         []string
// 	UsageRights  string
// }
//
// type CreatorInfo struct {
// 	Username        string
// 	Following       int64
// 	Followers       int64
// 	ProfilePhotoSrc string
// 	Bio             string
// }
//
// type Playable struct {
// 	ID              int32
// 	UserID          pgtype.UUID
// 	AudioSrc        string
// 	Username        string
// 	Title           string
// 	Bpm             int32
// 	PlaybackSeconds int32
// 	Created         pgtype.Timestamp
// 	ArtworkSrc      string
// }

type TrackPage struct {
	*PlayableElt
	Pins []*Pin
	*sqlc.GetUserInfoRow
	VisArr []int32
}

func (p *PlayableService) GetTrackPage(username string, title string) *TrackPage {
	res, err := db.Query.GetPlayableByTitleAndUsername(p.ctx, sqlc.GetPlayableByTitleAndUsernameParams{
		Title:    title,
		Username: username,
	})

	userService := NewUserService(p.ctx)

	user := userService.GetUserByUsername(username)
	pins := user.GetPins(p.ctx)

	numPins := len(pins)

	numOther := 4 - numPins

	if numOther > 0 {
		others, err := db.Query.GetUserAudioFilesWithLimit(p.ctx, sqlc.GetUserAudioFilesWithLimitParams{
			UserID:          user.ID,
			NumberOfResults: int32(numOther),
		})
		fmt.Printf("others: %v\n", others)
		util.EMsg(err, "getting other audio files (less than 4 pins found)")
		for _, other := range others {
			playable := NewPlayableElt(&other.Playable, p.ctx)
			pins = append(pins,
				&Pin{
					playable,
					other.TagArray,
				})
		}
	}

	util.EMsg(err, "getting playable elt")

	playableElt := NewPlayableElt(&res.Playable, p.ctx)

	userInfo, err := db.Query.GetUserInfo(p.ctx, username)
	util.EMsg(err, "getting user info")

	return &TrackPage{
		playableElt,
		pins,
		&userInfo,
		res.VisArr,
	}
}

// type PlayableEltData struct {
// 	*PlayableEltConstructor
// 	EltId string
// 	*user.User
// 	AudioSrc   string
// 	ArtworkSrc string
// 	Tags       *[]string
// 	IsPlaying  bool
// 	Bpm        int32
// }
//
//
// type User struct {
// 	*sqlc.User
// 	ProfilePhotoSrc string
// 	Followers       int32
// 	Following       int32
// 	Bio             string
// }
//
// type PlayableInfo struct {
// 	*MusicalKeyInfo
// 	Bpm          int32
// 	Tags         []string
// 	VisArr       *[]int32
// 	UsageRights  string
// 	UploadedTime string
// }
//
