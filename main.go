package main

import (
	"fmt"
	"log"
	"os"
)

/**
Plan:
 - Retrieve environmental variables in .aws folder
 - Access S3 bucket: "vlr-scrape"
 - Navigate to the folder "thread"
 - Upload the file "outputThreads.json"
**/

func check(err error) {
	if err != nil {
		fmt.Println("Error occurred:", err)
		log.Fatal("Error occurred:", err)
		os.Exit(1)
	}
}

func main() {
	// Check if endpoint can be reached
	docThread := threadPrep("https://www.vlr.gg/threads")
	// Scrape from VLR.gg threads.
	threadScrape(docThread)

	// Check if endpoint can be reached
	docMatch := threadPrep("https://www.vlr.gg/matches")
	// Scrape from VLR.gg matches.
	matchScrape(docMatch)

	// Retrieve context and search buckets available to given environmental variables.
	// ctx := context.Background()
	// s3c := AWSService.S3Client
	// AWSService.ListBuckets(AWSService{S3Client: }, ctx)
}
