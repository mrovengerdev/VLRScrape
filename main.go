package main

import (
	"fmt"
)

/**
Plan:
 - Retrieve environmental variables in .aws folder
 - Access S3 bucket: "vlr-scrape"
 - Navigate to the folder "thread"
 - Upload the file "outputThreads.json"
**/

func main() {
	// Check if endpoint can be reached
	doc := threadPrep("https://www.vlr.gg/threads")

	// Scrape from VLR.gg threads.
	threadScrape(doc)

	// Retrieve context and search buckets available to given environmental variables.
	// ctx := context.Background()
	// AWSService.ListBuckets(AWSService{}, ctx)

	fmt.Println("Scrape complete.")
}
