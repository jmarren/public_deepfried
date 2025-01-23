package services

import (
	"context"
	"github.com/jackc/pgx/v5/pgtype"
	"github.com/jmarren/deepfried/db"
	"github.com/jmarren/deepfried/sqlc"
	"github.com/jmarren/deepfried/util"
)

type StemService struct {
	ctx context.Context
}

func NewStemService(ctx context.Context) *StemService {
	s := new(StemService)
	s.ctx = ctx
	return s
}

func (s *StemService) AddStems(stemFileNames []string, audioFileId pgtype.UUID) {
	for _, stemFileName := range stemFileNames {
		s.AddStem(stemFileName, audioFileId)
	}
}

func (s *StemService) AddStem(stemFileName string, audioFileId pgtype.UUID) {
	err := db.Query.AddStem(s.ctx, sqlc.AddStemParams{
		AudioFileID:  audioFileId,
		StemFileName: stemFileName,
	})
	util.EMsg(err, "adding stem file")
}
