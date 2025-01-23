package services

import (
	"context"
	"fmt"

	"github.com/jmarren/deepfried/sqlc"
	"github.com/jmarren/deepfried/util"
)

type PlayableElt struct {
	*sqlc.Playable
	EltId           string
	ProfilePhotoSrc string
	DisplayTime     string
	Tags            []string
	IsPlaying       bool
}

func (p *PlayableElt) GetPlayingFromCtx(ctx context.Context) string {
	res := ctx.Value("playing")
	playing, ok := res.(string)
	if !ok {
		return ""
	}
	return playing
}

func NewPlayableElt(data *sqlc.Playable, ctx context.Context) *PlayableElt {
	p := new(PlayableElt)
	p.Playable = data
	p.AdjustSrcs()
	p.EltId = modifyInvalidChars(fmt.Sprintf("%s-%s", p.Username, p.Title))
	playableService := NewPlayableService(ctx)
	playing := playableService.GetPlayingFromCtx()

	if util.UuidString(p.ID) == playing {
		p.IsPlaying = true
	}
	p.ProfilePhotoSrc = getProfileSrc(p.UserID)
	p.DisplayTime = getTimeForUi(p.Created)
	return p
}
