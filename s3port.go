package main

import (
	"context"
	"errors"
	"io"
	"log"
	"os"

	"github.com/aws/aws-sdk-go-v2/service/s3"
	"github.com/aws/aws-sdk-go-v2/service/s3/types"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/smithy-go"
)

type AWSService struct {
	S3Client *s3.Client
}

func fileChecker(fileName string) {
	file, error := os.Open(fileName)
	if error != nil {
		log.Println("There was an error opening the file.")
	} else {
		defer file.Close()
	}
}

// ListBuckets lists the buckets in the current account.
// ListBuckets lists the buckets in the current account.
func (service AWSService) ListBuckets(ctx context.Context) ([]types.Bucket, error) {
	result, err := service.S3Client.ListBuckets(ctx, &s3.ListBucketsInput{})
	var buckets []types.Bucket
	if err != nil {
		log.Printf("Couldn't list buckets for your account. Here's why: %v\n", err)
	} else {
		buckets = result.Buckets
	}
	return buckets, err
}

// func (service AWSService) ListBuckets(ctx context.Context) ([]types.Bucket, error) {
// 	fmt.Println("Here 1")
// 	result, err := service.S3Client.ListBuckets(ctx, &s3.ListBucketsInput{}) // ERROR
// 	fmt.Println("Here 2")
// 	var buckets []types.Bucket
// 	fmt.Println("Here 3")
// 	if err != nil {
// 		log.Printf("Couldn't list buckets for your account. Here's why: %v\n", err)
// 	} else {
// 		buckets = result.Buckets
// 	}
// 	fmt.Println("Here 4")
// 	return buckets, err
// }

// ListObjects lists the objects in a bucket.
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

// BucketExists checks whether a bucket exists in the current account.
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

// UploadFile reads from a file and puts the data into an object in a bucket.
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

// UploadLargeObject uses an upload manager to upload data to an object in a bucket.
// The upload manager breaks large data into parts and uploads the parts concurrently.
// func (basics BucketBasics) UploadLargeObject(ctx context.Context, bucketName string, objectKey string, largeObject []byte) error {
// 	largeBuffer := bytes.NewReader(largeObject)
// 	var partMiBs int64 = 10
// 	uploader := manager.NewUploader(basics.S3Client, func(u *manager.Uploader) {
// 		u.PartSize = partMiBs * 1024 * 1024
// 	})
// 	_, err := uploader.Upload(ctx, &s3.PutObjectInput{
// 		Bucket: aws.String(bucketName),
// 		Key:    aws.String(objectKey),
// 		Body:   largeBuffer,
// 	})
// 	if err != nil {
// 		log.Printf("Couldn't upload large object to %v:%v. Here's why: %v\n",
// 			bucketName, objectKey, err)
// 	}

// 	return err
// }

// DownloadFile gets an object from a bucket and stores it in a local file.
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

// DownloadLargeObject uses a download manager to download an object from a bucket.
// The download manager gets the data in parts and writes them to a buffer until all of
// the data has been downloaded.
// func (basics BucketBasics) DownloadLargeObject(ctx context.Context, bucketName string, objectKey string) ([]byte, error) {
// 	var partMiBs int64 = 10
// 	downloader := manager.NewDownloader(basics.S3Client, func(d *manager.Downloader) {
// 		d.PartSize = partMiBs * 1024 * 1024
// 	})
// 	buffer := manager.NewWriteAtBuffer([]byte{})
// 	_, err := downloader.Download(ctx, buffer, &s3.GetObjectInput{
// 		Bucket: aws.String(bucketName),
// 		Key:    aws.String(objectKey),
// 	})
// 	if err != nil {
// 		log.Printf("Couldn't download large object from %v:%v. Here's why: %v\n",
// 			bucketName, objectKey, err)
// 	}
// 	return buffer.Bytes(), err
// }
