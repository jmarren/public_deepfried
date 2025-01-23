package handlers

import (
	"context"
	// "encoding/json"
	"bytes"
	"fmt"
	// "github.com/aws/aws-sdk-go-v2/service/s3/types"
	"io"
	"net/http"
	// "net/url"
	"strconv"
	"time"

	// "strconv"

	// "github.com/jmarren/deepfried/awssdk"
	"github.com/jmarren/deepfried/awssdk"
	"github.com/jmarren/deepfried/cache"
	"github.com/jmarren/deepfried/services"
	"github.com/jmarren/deepfried/util"
	// "github.com/jmarren/deepfried/util"
)

type PostUpload struct {
	UserService
	UploadService
	// ctx           context.Context
}

// audio/68f25788-73be-40a3-8e0c-55d6da6c3850/splendid+soup_70_A%23_%40lovejames(4).wav    --- AWS
// audio/68f25788-73be-40a3-8e0c-55d6da6c3850/splendid%20soup_70_A#_@lovejames(4).wav	  --- Golang

func NewPostUpload(ctx context.Context, us UserService) *PostUpload {
	pu := new(PostUpload)
	pu.UserService = us
	pu.UploadService = services.NewUploadService(ctx)
	return pu
}

type UploadItem struct {
	TotalSize   int32
	CurrentSize int32
	Data        []byte
	UploadType  string
	ContentType string
	Dest        string
}

func (pu *PostUpload) createNew(u *uploadHeader) *UploadItem {
	user, _ := pu.GetFromCtx()
	userId := user.ID

	dest := fmt.Sprintf("%s/%s/%s", u.UploadType, userId.String(), u.Filename)

	uploadItem := &UploadItem{
		TotalSize:   u.TotalSize,
		CurrentSize: 0,
		Data:        []byte{},
		UploadType:  u.UploadType,
		ContentType: u.ContentType,
		Dest:        dest,
	}

	cacheId := fmt.Sprintf("%s-%s", u.FrontendId, u.UploadType)
	cache.AppCache.Set(cacheId, uploadItem, time.Minute*10)
	return uploadItem
}

type uploadHeader struct {
	Index       int32
	ChunkSize   int32
	TotalSize   int32
	Filename    string
	FrontendId  string
	UploadType  string
	ContentType string
}

func (pu *PostUpload) ParseReqHeaders(r *http.Request) *uploadHeader {
	indexStr := r.Header.Get("X-Index")
	index, err := strconv.Atoi(indexStr)
	index32 := int32(index)
	util.EMsg(err, "converting 'X-Index' to int")

	chunkSizeStr := r.Header.Get("Content-Length")
	chunkSize, err := strconv.Atoi(chunkSizeStr)
	chunkSize32 := int32(chunkSize)

	totalSizeStr := r.Header.Get("X-total-size")
	totalSize, err := strconv.Atoi(totalSizeStr)
	totalSize32 := int32(totalSize)
	util.EMsg(err, "converting 'X-Index' to int")

	filename := r.Header.Get("X-file-name")
	frontendUploadId := r.Header.Get("X-upload-id")
	uploadType := r.Header.Get("X-upload-type")
	contentType := r.Header.Get("Content-Type")

	fmt.Printf("\tindex: %d\n\tchunkSize %d\n\ttotalSize: %d\n\tfilename: %s\n\tX-upload-id: %s\n\tuploadType: %s\n\t", index, chunkSize32, totalSize, filename, frontendUploadId, uploadType)

	u := new(uploadHeader)
	u.Index = index32
	u.ChunkSize = chunkSize32
	u.TotalSize = totalSize32
	u.Filename = filename
	u.FrontendId = frontendUploadId
	u.UploadType = uploadType
	u.ContentType = contentType

	return u
}

func (pu *PostUpload) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	_, auth := pu.GetFromCtx()
	if !auth {
		w.WriteHeader(http.StatusUnauthorized)
		return
	}

	uploadHeaders := pu.ParseReqHeaders(r)
	cacheId := fmt.Sprintf("%s-%s", uploadHeaders.FrontendId, uploadHeaders.UploadType)
	res, found := cache.AppCache.Get(cacheId)

	uploadItem, ok := res.(*UploadItem)

	if !found || !ok || uploadItem == nil {
		fmt.Println("not found in cache")
		uploadItem = pu.createNew(uploadHeaders)
	}

	data, err := io.ReadAll(r.Body)
	util.EMsg(err, "reading body")

	// append data to uploadItem
	uploadItem.Data = append(uploadItem.Data, data...)
	// incrememnt currentsize
	uploadItem.CurrentSize += int32(len(data))

	fmt.Printf("\t\tCurrentSize: %d\n\t\t\tTotalSize: %d\n", uploadItem.CurrentSize, uploadHeaders.TotalSize)

	// if current is equal to total, complete upload
	if uploadItem.TotalSize == uploadItem.CurrentSize {
		fmt.Println("totalsize is equal to current size")
		err = pu.CompleteUpload(uploadItem, uploadHeaders.FrontendId)
		if err == nil {
			fmt.Fprintf(w, "complete")
		}
	} else {
		fmt.Fprintf(w, "success")
	}
}

func (pu *PostUpload) CompleteUpload(u *UploadItem, frontendId string) error {
	fmt.Println("completing upload...")
	buf := new(bytes.Buffer)
	buf.Write(u.Data)
	err := awssdk.UploadBufferToS3(buf, u.Dest, u.ContentType)

	if u.UploadType == "audio" {
		pu.HandleAudioArr(u, frontendId)
	}
	util.EMsg(err, "uploading buffer to s3")
	if err != nil {
		fmt.Printf("uploaded item to %s\n", u.Dest)
	}

	cache.AppCache.Delete(frontendId)
	return err

}

func (pu *PostUpload) HandleAudioArr(u *UploadItem, frontendId string) {
	buf := new(bytes.Buffer)
	buf.Write(u.Data)
	visArr, err := services.ParseAudioArr(buf)

	util.EMsg(err, "parsing vis array")
	fmt.Printf("visArr: %v\n", visArr)

	visCacheId := fmt.Sprintf("%s-vis", frontendId)
	cache.AppCache.Set(visCacheId, visArr, 10*time.Minute)
	fmt.Printf("added vis array to cache with id: %s\n", visCacheId)
}

//
// func (pu *PostUpload) getDest(frontendId string) {
//
// }

// cache item key is frontendId
// cache holds awsUploadId

// type UploadCacheItem struct {
// 	MainFileName string
// }

//	'Content-Type': fileType,
//	'X-index': index,
//	'X-total-size': fileSize,
//	'X-file-name': fileName,
//	'X-Upload-Id':  fileId,
//
// // type user_ctx_key_type string
//
// const user_ctx_key user_ctx_key_type = "user"

// func CompleteUpload(ctx context.Context, dest string, uploadId *string) error {
/*
func (pu *PostUpload) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	fmt.Println("new post upload")
	uploadId := r.Header.Get("X-Upload-Id")

	if uploadId == "" {
		w.WriteHeader(http.StatusNotAcceptable)
	}

	switch r.PathValue("action") {
	case "new":
		pu.HandleNew(w, r)
		return
	case "append":
		pu.HandleAppend(w, r)
		return
	}
	w.WriteHeader(http.StatusNotFound)
}
*/
/*
func (pu *PostUpload) HandleAppend(w http.ResponseWriter, r *http.Request) {

	filename := r.Header.Get("X-file-name")
	chunkIndex, err := strconv.Atoi(r.Header.Get("X-index"))
	uploadType := r.PathValue("type")
	body, err := io.ReadAll(r.Body)
	if err != nil {
		fmt.Printf("Error reading body %s", err)
		w.WriteHeader(http.StatusInternalServerError)
	}
	err = pu.uploadService.Append(filename, uploadType, body, chunkIndex)
	util.EMsg(err, "error appending to upload item from cache")
}

func (pu *PostUpload) HandleNew(w http.ResponseWriter, r *http.Request) {
	uploadType := r.PathValue("type")
	if uploadType != "audio" && uploadType != "artwork" {
		w.WriteHeader(http.StatusNotAcceptable)
		return
	}

	filename := r.FormValue("filename")
	totalSizeStr := r.FormValue("totalSize")
	totalSize, err := strconv.Atoi(totalSizeStr)
	util.EMsg(err, "converting total size to int")

	pu.uploadService.NewUpload(filename, totalSize, uploadType)

	if uploadType == "audio" {
		w.Header().Set("HX-Trigger-After-Settle", "upload-audio-files")
		fmt.Printf("file_%s", filename)
		NewUploadForm(r.Context(), filename).ServeHTTP(w, r)
	}

	if uploadType == "artwork" {
		pu.uploadService.NewUpload(filename, totalSize, uploadType)
	}
}
*/
