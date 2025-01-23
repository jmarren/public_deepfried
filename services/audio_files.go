package services

import (
	"context"
	"github.com/jackc/pgx/v5/pgtype"
	"github.com/jmarren/deepfried/consts"
	"github.com/jmarren/deepfried/db"
	"github.com/jmarren/deepfried/sqlc"
	"github.com/jmarren/deepfried/util"
)

type AudioService struct {
	ctx context.Context
}

func NewAudioService(ctx context.Context) *AudioService {
	a := new(AudioService)
	a.ctx = ctx
	return a
}

// type AddAudioFileParams struct {
// 	UserID              pgtype.UUID
// 	Title               string
// 	AudioSrc            string
// 	Bpm                 int32
// 	MusicalKey          consts.MusicalKey
// 	MusicalKeySignature consts.MusicalKeySignature
// 	MajorMinor          consts.MajorMinor
// 	PlaybackSeconds     int32
// 	FileSize            int32
// 	UsageRights         pgtype.Text
// 	ArtworkSrc          string
// 	VisArr              []int32
// }

func (a *AudioService) WriteAudioFileToDb(userId pgtype.UUID, title string, audioSrc string, bpm int32, musicalKey consts.MusicalKey, keySig consts.MusicalKeySignature, majorMinor consts.MajorMinor, usageRights string, playbackSeconds int, visArr []int32, artworkSrc string, filesize int32, tags []string, stemFileNames []string) {
	audioFileId, err := db.Query.AddAudioFile(a.ctx, sqlc.AddAudioFileParams{
		UserID:              userId,
		Title:               title,
		Bpm:                 bpm,
		MusicalKey:          musicalKey,
		MusicalKeySignature: keySig,
		MajorMinor:          majorMinor,
		UsageRights: pgtype.Text{
			String: usageRights,
			Valid:  true,
		},
		PlaybackSeconds: int32(playbackSeconds),
		VisArr:          visArr,
		ArtworkSrc:      artworkSrc,
		AudioSrc:        audioSrc,
	})
	util.EMsg(err, "adding audio file to db")

	tagService := NewTagService(a.ctx)
	tagService.AddTags(tags, audioFileId)

	stemService := NewStemService(a.ctx)
	stemService.AddStems(stemFileNames, audioFileId)

}
