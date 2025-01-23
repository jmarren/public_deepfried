package services

import (
	"context"
	// "fmt"
	// "github.com/jackc/pgx/v5/pgtype"
	"github.com/jmarren/deepfried/db"
	"github.com/jmarren/deepfried/sqlc"
	"github.com/jmarren/deepfried/util"
	"strconv"
)

type SearchService struct{}

func NewSearchService() *SearchService {
	s := new(SearchService)
	return s
}

type AudioSearchResults struct {
	AudioFiles *[]AudioSearchRow
}

type AudioSearchRow struct {
	*PlayableElt
	VisArr          []int32
	ProfilePhotoSrc string
	DisplayKey      string
	StemFileNames   []string
	TagArray        []string
}

type SearchAudioFilters struct {
	Page                   int
	Keyword                string
	BpmRadioStr            string
	ExactBpmStr            string
	MaxBpmStr              string
	MinBpmStr              string
	MusicalKeyStr          string
	MajorMinorStr          string
	MusicalKeySignatureStr string
	IncludesStemsOnlyStr   string
}

func (s *SearchService) GetAudioSearchResults(ctx context.Context, filtersStrs *SearchAudioFilters) *AudioSearchResults {

	var bpmRadio string
	var exactBpm int
	var minBpm int
	var maxBpm int
	var err error

	if filtersStrs.ExactBpmStr != "" {
		exactBpm, err = strconv.Atoi(filtersStrs.ExactBpmStr)
		util.EMsg(err, "converting exactBpmStr to integer")
	} else {
		exactBpm = -1
	}

	if filtersStrs.MinBpmStr != "" {
		minBpm, err = strconv.Atoi(filtersStrs.MinBpmStr)
		util.EMsg(err, "converting minBpmStr to integer")
	} else {
		minBpm = -1
	}

	if filtersStrs.MaxBpmStr != "" {
		maxBpm, err = strconv.Atoi(filtersStrs.MaxBpmStr)
		util.EMsg(err, "converting maxBpmStr to integer")
	} else {
		maxBpm = -1
	}

	stemsOnly := false
	if filtersStrs.IncludesStemsOnlyStr == "on" {
		stemsOnly = true
	}

	if filtersStrs.BpmRadioStr == "use-exact" {
		bpmRadio = "use-exact"
	} else if filtersStrs.BpmRadioStr == "use-range" {
		bpmRadio = "use-range"
	} else {
		minBpm = -1
		maxBpm = -1
		exactBpm = -1
		bpmRadio = ""
	}

	dbRows, err := db.Query.SearchAudioFiles(ctx,
		sqlc.SearchAudioFilesParams{
			BpmRadio:   bpmRadio,
			ExactBpm:   int32(exactBpm),
			MinBpm:     int32(minBpm),
			MaxBpm:     int32(maxBpm),
			MusicalKey: filtersStrs.MusicalKeyStr,
			KeySig:     filtersStrs.MusicalKeySignatureStr,
			MajorMinor: filtersStrs.MajorMinorStr,
			StemsOnly:  stemsOnly,
			Keyword:    filtersStrs.Keyword,
			PageOffset: int32((filtersStrs.Page * 20) - 20),
		})
	util.EMsg(err, "searching for audio files")

	var audioFileRows []AudioSearchRow

	for _, row := range dbRows {
		playableElt := NewPlayableElt(&row.Playable, ctx)
		profilePhotoSrc := getProfileSrc(row.Playable.UserID)
		displayKey := getDisplayMusicalKey(row.MusicalKey, row.MusicalKeySignature, row.MajorMinor)
		audioFileRow := &AudioSearchRow{
			playableElt,
			row.VisArr,
			profilePhotoSrc,
			displayKey,
			row.StemFileNames,
			row.TagArray,
		}

		audioFileRows = append(audioFileRows, *audioFileRow)
	}

	return &AudioSearchResults{
		&audioFileRows,
	}
}

type UserSearchResults struct {
	Users []*UserSearchRow
}

type UserSearchRow struct {
	*sqlc.UserSearchRow
	ProfilePhotoSrc string
}

func (s *SearchService) GetUserSearchResults(ctx context.Context, keyword *string) *UserSearchResults {
	dbRows, err := db.Query.SearchForUsers(ctx, *keyword)
	util.EMsg(err, "searching for users")

	var userSearchRows []*UserSearchRow

	for _, row := range dbRows {
		profilePhotoSrc := getProfileSrc(row.ID)
		userSearchRows = append(userSearchRows, &UserSearchRow{
			&row,
			profilePhotoSrc,
		})
	}
	return &UserSearchResults{
		userSearchRows,
	}
}

func (s *SearchService) GetDropdownItems(ctx context.Context, keyword *string) []*sqlc.SearchKeywordForDropdownRow {
	res, err := db.Query.SearchKeywordForDropdown(ctx, *keyword)
	util.EMsg(err, "searching for dropdown items")

	var data []*sqlc.SearchKeywordForDropdownRow

	for _, item := range res {
		itemPtr := &item
		itemPtr.ArtworkSrc = getArtSrc(item.UserID, item.ArtworkSrc)
		data = append(data, itemPtr)
	}

	return data
}
