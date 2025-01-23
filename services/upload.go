package services

import (
	"bytes"
	"context"
	"fmt"
	"net/http"
	"slices"
	"time"

	"github.com/jmarren/deepfried/awssdk"
	"github.com/jmarren/deepfried/cache"
	"github.com/jmarren/deepfried/util"
)

type UploadService struct {
	ctx context.Context
}

// creates a new upload service with context
func NewUploadService(ctx context.Context) *UploadService {
	u := new(UploadService)
	u.ctx = ctx
	return u
}

type Upload struct {
	CacheId     string
	UploadType  string
	CurrentSize int
	TotalSize   int
	Destination string
	Chunks      []*Chunk
}

type Chunk struct {
	Index int
	Data  *[]byte
}

// creates a new instance of upload, generates a cacheId and adds the item to the cache
// it should generate a uuid that will be passed to the frontend ?
func (u *UploadService) NewUpload(filename string, totalSize int, uploadType string) {
	upload := new(Upload)
	upload.TotalSize = totalSize
	upload.UploadType = uploadType
	upload.CurrentSize = 0

	user, _ := NewUserService(u.ctx).GetFromCtx()
	upload.CacheId = fmt.Sprintf("%s-%s", UuidString(user.ID), filename)
	upload.Destination = fmt.Sprintf("%s/%s/%s/", uploadType, UuidString(user.ID), filename)

	fmt.Printf("upload added to cache with id: %s\n", upload.CacheId)
	cache.AppCache.Set(upload.CacheId, upload, 10*time.Minute)
}

func (u *UploadService) Append(filename string, uploadType string, data []byte, index int) error {
	upload, err := u.getCacheItem(filename)
	util.EMsg(err, "err getting cache item: %s\n")

	chunk := u.newChunk(data, index)
	upload.Chunks = append(upload.Chunks, chunk)
	upload.CurrentSize += len(*chunk.Data)
	if upload.CurrentSize == upload.TotalSize {
		data, err := upload.assembleChunks()
		if err != nil {
			return err
		}
		upload.CompleteUpload(data)
	}
	return nil
}

func (u *UploadService) newChunk(data []byte, index int) *Chunk {
	c := new(Chunk)
	c.Data = &data
	c.Index = index
	return c
}

func (u *UploadService) getCacheItem(filename string) (*Upload, error) {
	user, auth := NewUserService(u.ctx).GetFromCtx()
	if !auth {
		return nil, fmt.Errorf("user not authenticated")
	}
	cacheId := fmt.Sprintf("%s-%s", UuidString(user.ID), filename)
	cacheItem, found := cache.AppCache.Get(cacheId)
	if !found {
		return nil, fmt.Errorf("item not found in cache")
	}
	upload, ok := cacheItem.(*Upload)
	if !ok {
		return nil, fmt.Errorf("item found in cache is not of type *Upload")
	}
	return upload, nil

}

func (u *Upload) assembleChunks() (*bytes.Buffer, error) {
	slices.SortFunc(u.Chunks, func(a *Chunk, b *Chunk) int {
		if a.Index > b.Index {
			return 1
		} else if a.Index < b.Index {
			return -1
		}
		return 0
	})

	outputBuf := new(bytes.Buffer)

	for i := 0; i < len(u.Chunks); i++ {
		_, err := outputBuf.Write(*u.Chunks[i].Data)
		if err != nil {
			return nil, err
		}
	}
	return outputBuf, nil
}

func (u *Upload) sendFileUploadToS3(dataBuf *bytes.Buffer) {
	dataBytes := dataBuf.Bytes()
	fileType := http.DetectContentType(dataBytes)
	fmt.Printf("fileType: %s\n", fileType)
	err := awssdk.UploadBufferToS3(dataBuf, u.Destination, fileType)
	if err != nil {
		fmt.Printf("Error uploading to S3: %s\n", err)
	}
}

func (u *Upload) CompleteUpload(data *bytes.Buffer) {
	u.sendFileUploadToS3(data)
	if u.UploadType == "audio" {
		visArr, err := ParseAudioArr(data)
		util.EMsg(err, "parsing vis array")
		fmt.Printf("visArr: %v\n", visArr)
		cache.AppCache.Delete(u.CacheId)
		cache.AppCache.Set(fmt.Sprintf("%s-vis", u.CacheId), visArr, 10*time.Minute)
	}
}
