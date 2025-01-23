package services

import (
	"context"
	// "fmt"
	// "github.com/jmarren/deepfried/db"
	// // "github.com/jmarren/deepfried/sqlc"
	// "github.com/jmarren/deepfried/util"
)

type PlayerService struct{}

func NewPlayerService() *PlayerService {
	p := new(PlayerService)
	return p
}

func (p *PlayerService) GetPlayerData(ctx context.Context) {

	// tags, err := db.Query.GetAudioFileTags(ctx, playable.ID)
	// util.EMsg(err, "getting tags for featured")
	// playable.Tags = tags

	// return playable
}

// func (p *PlayerService) GetEditorsPicks(ctx context.Context) []*PlayableElt {
// 	dbData, err := db.Query.GetEditorsPicks(ctx)
// 	util.EMsg(err, "Getting editors picks")
// 	var playables []*PlayableElt
// 	for _, row := range dbData {
// 		playable := NewPlayableElt(&row)
// 		playables = append(playables, playable)
// 	}
// 	return playables
// }
