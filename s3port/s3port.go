package s3port

import (
	"context"
	"fmt"
	"log"
	"os"
	"path/filepath"

	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/credentials"
	"github.com/aws/aws-sdk-go-v2/feature/s3/manager"
	"github.com/aws/aws-sdk-go-v2/service/s3"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/joho/godotenv"
)

type AWSService struct {
	S3Client *s3.Client
}

type fileWalk chan string

// Traverses files
func (f fileWalk) Walk(path string, info os.FileInfo, err error) error {
	if err != nil {
		return err
	}
	if !info.IsDir() {
		f <- path
	}
	return nil
}

// Traverses through all files in the output folder through fileWalk and uploads them to Amazon S3.
// Does so by retrieving credentials, creating a new S3 client, and parsing through the output folder.
func Upload() {
	// Retrieve S3 credentials from .env file via godotenv.
	err := godotenv.Load()
	if err != nil {
		log.Fatal("Error loading .env file")
	}

	accessKey := os.Getenv("AWS_VLR_ACCESS_KEY_ID")
	secretKey := os.Getenv("AWS_VLR_SECRET_KEY")
	s3Bucket := os.Getenv("AWS_VLR_S3_BUCKET")
	s3Region := os.Getenv("AWS_VLR_S3_REGION")

	fmt.Println(accessKey, secretKey, s3Bucket, s3Region)

	if accessKey == "" || secretKey == "" {
		log.Fatal("AWS credentials not found in .env file")
	}

	// Retrieve path to output file.
	localPath, err := os.Getwd()
	if err != nil {
		log.Fatalf("Error: %v", err)
	}
	localPath = filepath.Join(localPath, "output")

	// Creates SDK configuration
	cfg, err := config.LoadDefaultConfig(context.TODO(),
		config.WithRegion(s3Region),
		config.WithCredentialsProvider(credentials.NewStaticCredentialsProvider(accessKey, secretKey, "")))
	if err != nil {
		log.Fatalln("error:", err)
	}

	client := s3.NewFromConfig(cfg)

	walker := make(fileWalk)
	go func() {
		// Gather the files to upload by walking the path recursively
		if err := filepath.Walk(localPath, walker.Walk); err != nil {
			log.Fatalln("Walk failed:", err)
		}
		close(walker)
	}()

	// For each file found, walking through output file, upload to Amazon S3
	uploader := manager.NewUploader(client)
	for path := range walker {
		rel, err := filepath.Rel(localPath, path)
		if err != nil {
			log.Fatalln("Unable to get relative path:", path, err)
		}
		file, err := os.Open(path)
		if err != nil {
			log.Println("Failed opening file", path, err)
			continue
		}
		defer file.Close()

		// Upload the file to S3 bucket given in .env file.
		result, err := uploader.Upload(context.TODO(), &s3.PutObjectInput{
			Bucket: aws.String(s3Bucket),
			Key:    aws.String(rel), // WIP: Don't need folders in my S3 bucket.
			Body:   file,
			ACL:    "public-read",
		})
		if err != nil {
			log.Fatalln("Failed to upload", path, err)
		}

		log.Printf("Uploaded the file: %s to S3 location: %s\n", path, result.Location)
	}
}
