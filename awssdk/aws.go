package awssdk

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"os"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/feature/s3/manager"
	"github.com/aws/aws-sdk-go-v2/service/s3"
	"github.com/aws/aws-sdk-go-v2/service/secretsmanager"
)

type UploadResult struct {
	Path string `json:"path" xml:"path"`
}

var (
	uploader   *manager.Uploader
	downloader *manager.Downloader
	s3Client   *s3.Client
	DbDsn      string
)

func InitAWS() {

	fmt.Println("initializing AWS")
	// AWS SDK
	// Load the Shared AWS Configuration (~/.aws/config)
	ctx := context.Background()
	cfg, err := config.LoadDefaultConfig(ctx, config.WithRegion("us-west-1"))
	if err != nil {
		log.Fatal(err)
	}

	initS3(cfg)
	if os.Getenv("env") == "prod" {
		err = initRds(ctx, cfg)
		if err != nil {
			fmt.Printf("error initializing rds: %s\n", err)
			panic(err)
		}
	}
}

func initS3(cfg aws.Config) {
	// Create an Amazon S3 service s3client
	s3Client = s3.NewFromConfig(cfg)
	uploader = manager.NewUploader(s3Client)
	downloader = manager.NewDownloader(s3Client)
	fmt.Println("initialized s3")
}

func initRds(ctx context.Context, cfg aws.Config) error {
	dbName := "loopdatabase"
	dbHost := os.Getenv("RDS_ENDPOINT")
	secretId := os.Getenv("RDS_SECRET_ID")
	smClient := secretsmanager.NewFromConfig(cfg)
	dbSecretOutput, err := smClient.GetSecretValue(
		ctx,
		&secretsmanager.GetSecretValueInput{
			SecretId: &secretId,
		})
	if err != nil {
		return err
	}

	dbPort := 5432
	type SecretOutput struct {
		Username string
		Password string
	}

	var secret SecretOutput
	secretOutputPtr := &secret

	err = json.Unmarshal([]byte(*dbSecretOutput.SecretString), secretOutputPtr)
	if err != nil {
		return err
	}
	fmt.Printf("secret: %s\n", secret)

	DbDsn = fmt.Sprintf(
		"host=%s port=%d user=%s password=%s dbname=%s",
		dbHost,
		dbPort,
		secret.Username,
		secret.Password,
		dbName,
	)

	fmt.Printf("DbDsn: %s\n", DbDsn)

	return nil
}
