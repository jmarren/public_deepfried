package services

import (
	"context"
	"fmt"
	"github.com/google/uuid"
	"github.com/jackc/pgx/v5/pgtype"
	"github.com/jmarren/deepfried/db"
	"github.com/jmarren/deepfried/sqlc"
	"github.com/jmarren/deepfried/util"
	// "github.com/jmarren/deepfried/sqlc"
	// "github.com/jmarren/deepfried/util"
	// "strconv"
)

type DownloadsService struct {
	ctx context.Context
}

func NewDownloadsService(ctx context.Context) *DownloadsService {
	d := new(DownloadsService)
	d.ctx = ctx
	return d
}

func (d *DownloadsService) AddUserDownload(userId pgtype.UUID, audioIdStr string) {
	// userService := NewUserService(d.Context)
	// user := userService.GetFromCtx()
	// if user == nil {
	// 	fmt.Println("user is nil")
	// 	return
	// }

	fmt.Printf("audioIdStr: %s\n", audioIdStr)
	fmt.Printf("userId: %s\n", userId.String())
	parsedId, err := uuid.Parse(audioIdStr)

	util.EMsg(err, "parsing download audio id")
	if err != nil {
		return
	}

	pgId := pgtype.UUID{
		Bytes: parsedId,
		Valid: true,
	}

	err = db.Query.AddUserDownload(d.ctx, sqlc.AddUserDownloadParams{
		UserID:  userId,
		AudioID: pgId,
	})

	util.EMsg(err, "adding download to db")

}
