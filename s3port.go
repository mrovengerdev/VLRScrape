package main

import (
	"context"
	"errors"
	"io"
	"log"
	"os"
	"path/filepath"

	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/credentials"
	"github.com/aws/aws-sdk-go-v2/feature/s3/manager"
	"github.com/aws/aws-sdk-go-v2/service/s3"
	"github.com/aws/aws-sdk-go-v2/service/s3/types"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/smithy-go"
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

// Traverses through all files in the output folder and uploads them to Amazon S3.
// Does so by retrieving credentials, creating a new S3 client, and parsing through the output folder.
func upload() {
	// Retrieve S3 credentials from .env file via godotenv.
	err := godotenv.Load()
	if err != nil {
		log.Fatal("Error loading .env file")
	}

	accessKey := os.Getenv("AWS_ACCESS_KEY_ID")
	secretKey := os.Getenv("AWS_SECRET_KEY")
	s3Bucket := os.Getenv("AWS_S3_BUCKET")

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
		config.WithRegion("us-west-2"),
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

		// Upload the file to S3 bucket: vlr-scrape.
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

// Lists the buckets in the current account.
func (service AWSService) ListBuckets(ctx context.Context) ([]types.Bucket, error) {
	result, err := service.S3Client.ListBuckets(ctx, &s3.ListBucketsInput{}) // I'm referencing a nil pointer
	var buckets []types.Bucket
	if err != nil {
		log.Printf("Couldn't list buckets for your account. Here's why: %v\n", err)
	} else {
		buckets = result.Buckets
	}
	return buckets, err
}

// Lists the objects in a bucket.
func (service AWSService) ListObjects(ctx context.Context, bucketName string) ([]types.Object, error) {
	result, err := service.S3Client.ListObjectsV2(ctx, &s3.ListObjectsV2Input{
		Bucket: aws.String(bucketName),
	})
	var contents []types.Object
	if err != nil {
		log.Printf("Couldn't list objects in bucket %v. Here's why: %v\n", bucketName, err)
	} else {
		contents = result.Contents
	}
	return contents, err
}

// Checks whether a bucket exists in the current account.
func (service AWSService) BucketExists(ctx context.Context, bucketName string) (bool, error) {
	_, err := service.S3Client.HeadBucket(ctx, &s3.HeadBucketInput{
		Bucket: aws.String(bucketName),
	})
	exists := true
	if err != nil {
		var apiError smithy.APIError
		if errors.As(err, &apiError) {
			switch apiError.(type) {
			case *types.NotFound:
				log.Printf("Bucket %v is available.\n", bucketName)
				exists = false
				err = nil
			default:
				log.Printf("Either you don't have access to bucket %v or another error occurred. "+
					"Here's what happened: %v\n", bucketName, err)
			}
		}
	} else {
		log.Printf("Bucket %v exists and you already own it.", bucketName)
	}

	return exists, err
}

// Reads from a file and puts the data into an object in a bucket.
func (service AWSService) UploadFile(ctx context.Context, bucketName string, objectKey string, fileName string) error {
	file, err := os.Open(fileName)
	if err != nil {
		log.Printf("Couldn't open file %v to upload. Here's why: %v\n", fileName, err)
	} else {
		defer file.Close()
		_, err = service.S3Client.PutObject(ctx, &s3.PutObjectInput{
			Bucket: aws.String(bucketName),
			Key:    aws.String(objectKey),
			Body:   file,
		})
		if err != nil {
			log.Printf("Couldn't upload file %v to %v:%v. Here's why: %v\n",
				fileName, bucketName, objectKey, err)
		}
	}
	return err
}

// Gets an object from a bucket and stores it in a local file.
func (basics AWSService) DownloadFile(ctx context.Context, bucketName string, objectKey string, fileName string) error {
	result, err := basics.S3Client.GetObject(ctx, &s3.GetObjectInput{
		Bucket: aws.String(bucketName),
		Key:    aws.String(objectKey),
	})
	if err != nil {
		log.Printf("Couldn't get object %v:%v. Here's why: %v\n", bucketName, objectKey, err)
		return err
	}
	defer result.Body.Close()
	file, err := os.Create(fileName)
	if err != nil {
		log.Printf("Couldn't create file %v. Here's why: %v\n", fileName, err)
		return err
	}
	defer file.Close()
	body, err := io.ReadAll(result.Body)
	if err != nil {
		log.Printf("Couldn't read object body from %v. Here's why: %v\n", objectKey, err)
	}
	_, err = file.Write(body)
	return err
}
