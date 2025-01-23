package handlers

import (
	"context"
	"encoding/json"
	"fmt"
	// "io"
	"github.com/jmarren/deepfried/cache"
	"github.com/jmarren/deepfried/consts"
	"github.com/jmarren/deepfried/services"
	"github.com/jmarren/deepfried/util"
	"net/http"
	"strconv"
)

type PostAudioForm struct {
	userService   UserService
	uploadService UploadService
}

func NewPostAudioForm(ctx context.Context, us UserService) *PostAudioForm {
	pa := new(PostAudioForm)
	pa.userService = us
	pa.uploadService = services.NewUploadService(ctx)
	return pa
}

type AudioUploadRequest struct {
	AudioFileName string
	ArtworkName   string
	Title         string
	Tags          []string
	Bpm           int32
	StemFileNames []string
	MusicalKey    consts.MusicalKey
	KeySig        consts.MusicalKeySignature
	MajorMinor    consts.MajorMinor
	UsageRights   string
	FrontendId    string
}

func (pa *PostAudioForm) ParseRequest(r *http.Request) (*AudioUploadRequest, error) {

	frontendId := r.FormValue("upload-id")

	fmt.Printf("frontendId: %s\n", frontendId)

	bpm, err := strconv.Atoi(r.FormValue("bpm"))
	util.EMsg(err, "converting bpm to int")

	artworkName := r.FormValue("artwork_file_name")

	title := r.FormValue("title")
	if title == "" {
		return nil, fmt.Errorf("please name your upload")
	}

	tags := []string{}
	tagOne := r.FormValue("tag-1")
	if tagOne != "" {
		tags = append(tags, tagOne)
	}
	tagTwo := r.FormValue("tag-2")
	if tagTwo != "" {
		tags = append(tags, tagTwo)
	}
	tagThree := r.FormValue("tag-3")
	if tagThree != "" {
		tags = append(tags, tagThree)
	}
	fmt.Printf("tags: %v\n", tags)

	audioFileName := r.FormValue("audio_file_name")
	if audioFileName == "" {
		return nil, fmt.Errorf("Upload failed. Please try again later")
	}

	stemFileNamesRaw := r.FormValue("stem_file_names")
	var stemFileNames []string
	_ = json.Unmarshal([]byte(stemFileNamesRaw), &stemFileNames)

	musicalKeyRaw := r.FormValue("musical-key")
	musicalKey := consts.MusicalKey(musicalKeyRaw)

	keySigRaw := r.FormValue("musical-key-signature")
	keySig := consts.MusicalKeySignature(keySigRaw)

	majorMinorRaw := r.FormValue("major-minor")
	majorMinor := consts.MajorMinor(majorMinorRaw)

	usageRights := r.FormValue("usage")
	a := new(AudioUploadRequest)
	a.Title = title
	a.ArtworkName = artworkName
	a.Tags = tags
	a.StemFileNames = stemFileNames
	a.Bpm = int32(bpm)
	a.AudioFileName = audioFileName
	a.MusicalKey = musicalKey
	a.KeySig = keySig
	a.MajorMinor = majorMinor
	a.UsageRights = usageRights
	a.FrontendId = frontendId

	return a, nil

}

func (pa *PostAudioForm) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	user, auth := pa.userService.GetFromCtx()
	if !auth {
		w.WriteHeader(http.StatusUnauthorized)
	}

	a, err := pa.ParseRequest(r)
	if err != nil {
		msg := fmt.Sprintf("%s", err)
		fmt.Fprintf(w, "%s", msg)
	}

	visCacheId := fmt.Sprintf("%s-vis", a.FrontendId)

	fmt.Printf("vis cacheId: %s\n", visCacheId)

	res, found := cache.AppCache.Get(visCacheId)
	if !found {
		fmt.Println("ERROR: visArr not found in cache")
		fmt.Fprintf(w, "an error occurred. please try again later")
	}
	visArr, ok := res.([]int32)
	if !ok {
		fmt.Println("ERROR: visArr from cache not of type []int32")
		fmt.Fprintf(w, "an error occurred. please try again later")
	}

	audioFileService := services.NewAudioService(ctx)

	audioFileService.WriteAudioFileToDb(user.ID, a.Title, a.AudioFileName, a.Bpm, a.MusicalKey, a.KeySig, a.MajorMinor, a.UsageRights, 12308, visArr, a.ArtworkName, 123893, a.Tags, a.StemFileNames)

	if err != nil {
		fmt.Printf("err querying to add audio file: %s\n", err)
		fmt.Fprintf(w, "an error occurred")
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	w.Header().Set("HX-Retarget", "#selected-files")
	fmt.Fprintf(w, "sucess!")

}
