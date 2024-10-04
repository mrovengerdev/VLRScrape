package main

import (
	"context"
	"fmt"
	"log"

	"github.com/aws/aws-sdk-go-v2/service/s3"
	"github.com/aws/aws-sdk-go-v2/service/s3/types"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/smithy-go"
)

/**
Plan:
 - Retrieve environmental variables in .aws folder
 - Access S3 bucket: "vlr-scrape"
 - Navigate to the folder "thread"
 - Upload the file "outputThreads.json"
**/

func main() {
doc := threadPrep()
	ctx := context.Background()
	sess := session.Must(session.NewSession())
	svc := s3.New(sess)
	
	buckets, err := svc.ListBucketsWithContext(ctx, nil)
	if err != nil {
		log.Printf("Couldn't list buckets for your account. Here's why: %v\n", err)
	buckets, err := service.ListBuckets(ctx)
	fmt.Printf("Buckets: %v\n", buckets)

	AWSService.ListBuckets(AWSService{}, ctx)
	// AWSService.UploadFile(ctx, "vlr-scrape", objectkey, "outputThreads.json")

	fmt.Println("Port to S3 complete.")
}