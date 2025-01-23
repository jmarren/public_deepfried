package services

import (
	"context"
	"fmt"
	"github.com/jackc/pgx/v5/pgtype"
	"github.com/jmarren/deepfried/db"
	"github.com/jmarren/deepfried/sqlc"
	"github.com/jmarren/deepfried/util"
)

type TagService struct {
	ctx context.Context
}

func NewTagService(ctx context.Context) *TagService {
	t := new(TagService)
	t.ctx = ctx
	return t
}

func (t *TagService) GetMostPopular() []string {
	fmt.Println("getting most popular tags!")

	tags := []string{}

	dbTags, err := db.Query.GetTagsOrderedByCount(t.ctx)
	if dbTags == nil || err != nil {
		return []string{"hi"}
	}

	fmt.Printf("dbTags: %v\n", dbTags)

	util.EMsg(err, "getting most popular tags")
	for _, tag := range dbTags {
		tags = append(tags, tag)
	}

	fmt.Printf("tags: %v\n", tags)
	return tags
}

func (t *TagService) AddTags(tags []string, audioFileId pgtype.UUID) {
	for _, tag := range tags {
		t.AddTag(tag, audioFileId)
	}
}

func (t *TagService) AddTag(tag string, audioFileId pgtype.UUID) {
	tagCount, err := db.Query.GetTagCount(t.ctx, tag)
	util.EMsg(err, "adding tag to db")

	var tagId int32
	if tagCount == 0 {
		tagId, err = db.Query.InsertTag(t.ctx, tag)
		util.EMsg(err, "inserting new tag")
	} else {
		tagId, err = db.Query.GetTagId(t.ctx, tag)
		util.EMsg(err, "getting tag id")
	}

	err = db.Query.AddAudioFileTag(t.ctx, sqlc.AddAudioFileTagParams{
		TagID:       tagId,
		AudioFileID: audioFileId,
	})
	util.EMsg(err, "adding audio file tag")

}
