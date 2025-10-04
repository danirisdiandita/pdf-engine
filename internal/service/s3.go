package service

import (
	"fmt"
	"os"
	"path/filepath"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/credentials"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/s3"
	"github.com/aws/aws-sdk-go/service/s3/s3manager"
	"github.com/danirisdiandita/pdf-engine/internal/config"
)

var s3client *s3.S3

func init() {
	cfg := config.Load()
	sess, err := session.NewSession(&aws.Config{
		Region:      aws.String(cfg.S3Region),
		Endpoint:    aws.String(cfg.S3Url),
		Credentials: credentials.NewStaticCredentials(cfg.S3AccessKeyId, cfg.S3SecretAccessKey, ""),
	})
	if err != nil {
		fmt.Printf("Failed to create AWS session: %v\n", err)
		return
	}
	s3client = s3.New(sess)
}

func DownloadFile(config *config.Config, objectKey string) (string, error) {

	filePath := fmt.Sprintf("tmp/%s", objectKey)

	dirName := filepath.Dir(filePath)
	if _, err := os.Stat(dirName); os.IsNotExist(err) {
		if err := os.MkdirAll(dirName, 0755); err != nil {
			return "", fmt.Errorf("failed to create directory %s: %v", dirName, err)
		}
	}

	sess, err := session.NewSession(&aws.Config{
		Region:      aws.String(config.S3Region),
		Endpoint:    aws.String(config.S3Url),
		Credentials: credentials.NewStaticCredentials(config.S3AccessKeyId, config.S3SecretAccessKey, ""),
	})
	if err != nil {
		return "", fmt.Errorf("failed to create AWS session: %v", err)
	}

	file, err := os.Create(filePath)
	if err != nil {
		return "", fmt.Errorf("failed to create file %s: %v", filePath, err)
	}
	defer file.Close()

	downloader := s3manager.NewDownloader(sess)
	_, err = downloader.Download(file, &s3.GetObjectInput{
		Bucket: aws.String(config.S3BucketName),
		Key:    aws.String(objectKey),
	})
	if err != nil {
		return "", fmt.Errorf("failed to download file: %v", err)
	}

	return filePath, nil
}

func UploadFile(config *config.Config, filePath string, objectKey string) error {
	sess, err := session.NewSession(&aws.Config{
		Region:      aws.String(config.S3Region),
		Endpoint:    aws.String(config.S3Url),
		Credentials: credentials.NewStaticCredentials(config.S3AccessKeyId, config.S3SecretAccessKey, ""),
	})
	if err != nil {
		return fmt.Errorf("failed to create AWS session: %v", err)
	}

	file, err := os.Open(filePath)
	if err != nil {
		return fmt.Errorf("failed to open file %s: %v", filePath, err)
	}
	defer file.Close()

	fmt.Println("uploading file to s3", filePath, objectKey)

	uploader := s3manager.NewUploader(sess)
	_, err = uploader.Upload(&s3manager.UploadInput{
		Bucket: aws.String(config.S3BucketName),
		Key:    aws.String(objectKey),
		Body:   file,
	})
	if err != nil {
		return fmt.Errorf("failed to upload file: %v", err)
	}
	return nil
}
