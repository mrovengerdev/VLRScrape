package main

import (
	"log"
)

func check(err error) {
	if err != nil {
		log.Fatalf("Error: %v", err)
	}
}

func main() {

	// Creates output folder for JSON files.
	createOutputDirectory()

	// Check if endpoint can be reached
	docThread := threadPrep("https://www.vlr.gg/threads")
	// Scrape VLR.gg threads.
	threadScrape(docThread)

	// Check if endpoint can be reached
	docMatch := threadPrep("https://www.vlr.gg/matches")
	// Scrape VLR.gg matches.
	matchScrape(docMatch)

	// Upload output files to Amazon S3 bucket: "vlr-scrape".
	upload()
}
