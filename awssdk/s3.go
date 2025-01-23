package awssdk

import (
	"bytes"
	"context"
	"errors"
	"fmt"
	// "io/fs"
	"os"
	"time"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/feature/s3/manager"
	"github.com/aws/aws-sdk-go-v2/service/s3"
	"github.com/jmarren/deepfried/util"
)

var loopS3PublicBucketName = aws.String("loop-public-s3-bucket")

func UploadBufferToS3(buf *bytes.Buffer, dest string, contentType string) error {
	ctx, cancel := context.WithTimeout(context.Background(), 50*time.Second)

	defer cancel()

	fmt.Println("buffer length: ", buf.Len())

	// fmt.Println("buf.Bytes()", string(buf.Bytes()))
	fmt.Println("upload--dest: ", dest)

	output, err := uploader.Upload(ctx, &s3.PutObjectInput{
		Bucket:      aws.String("loop-public-s3-bucket"),
		Key:         aws.String(dest),
		Body:        bytes.NewReader(buf.Bytes()),
		ContentType: aws.String(contentType),
	})

	fmt.Println("output: ", output)
	if err != nil {
		var mu manager.MultiUploadFailure
		if errors.As(err, &mu) {
			// Process error and its associated uploadID
			fmt.Println("Error:", mu)
			_ = mu.UploadID() // retrieve the associated UploadID
		} else {
			// Process error generically
			fmt.Println("Error:", err.Error())
		}
		return err
	}

	return nil
}

func DownloadFromS3(location string) ([]byte, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 20*time.Second)

	defer cancel()
	buf := manager.NewWriteAtBuffer([]byte{})

	numBytes, err := downloader.Download(ctx, buf, &s3.GetObjectInput{
		Bucket: aws.String("loop-public-s3-bucket"),
		Key:    aws.String(location),
	})
	if err != nil {
		fmt.Println("error downloading content: ", err)
		return nil, err
	}
	fmt.Println("bytes received: ", numBytes)
	return buf.Bytes(), nil
}

func UploadStaticAssets() {
	// if os.Getenv("env") != "porcupine!" {
	// 	return
	// }
	jsDirPath := "./static/assets/js/"

	jsDir, err := os.ReadDir(jsDirPath)

	util.EMsg(err, "opening js dir")

	for _, jsFile := range jsDir {
		buf := new(bytes.Buffer)
		name := jsFile.Name()
		file, err := os.ReadFile(fmt.Sprintf("%s%s", jsDirPath, name))
		util.EMsg(err, fmt.Sprintf("reading js file: %s\n", name))
		buf.Write(file)
		destination := fmt.Sprintf("public/js/%s", name)
		err = UploadBufferToS3(buf, destination, "text/javascript")
		util.EMsg(err, "uploading js files to S3")

	}

	cssDirPath := "./static/assets/css/"

	cssDir, err := os.ReadDir(cssDirPath)

	util.EMsg(err, "opening css dir")

	for _, cssFile := range cssDir {
		buf := new(bytes.Buffer)
		name := cssFile.Name()
		file, err := os.ReadFile(fmt.Sprintf("%s%s", cssDirPath, name))
		util.EMsg(err, fmt.Sprintf("reading css file: %s\n", name))
		buf.Write(file)
		destination := fmt.Sprintf("public/css/%s", name)
		err = UploadBufferToS3(buf, destination, "text/css")
		util.EMsg(err, "uploading css files to S3")
	}

	/*
	   fontDirPath := "./static/assets/fonts/"

	   fontDir, err := os.ReadDir(fontDirPath)

	   util.EMsg(err, "opening css dir")

	   	for _, cssFile := range cssDir {
	   		buf := new(bytes.Buffer)
	   		name := cssFile.Name()
	   		file, err := os.ReadFile(fmt.Sprintf("%s%s", cssDirPath, name))
	   		util.EMsg(err, fmt.Sprintf("reading font file: %s\n", name))
	   		buf.Write(file)
	   		destination := fmt.Sprintf("public/css/%s", name)
	   		err = UploadBufferToS3(buf, destination, "text/css")
	   		util.EMsg(err, "uploading css files to S3")
	   	}
	*/
}
