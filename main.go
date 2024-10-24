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

	// Scrape from VLR.gg threads.
	pageParser("https://www.vlr.gg/threads", "/?t=1w", "outputThreads")

	// Scrape from VLR.gg matches.
	pageParser("https://www.vlr.gg/matches", "/?", "outputMatches")

	// Upload output files to Amazon S3 bucket: "vlr-scrape".
	upload()
}

/*
Plan on features to add:
- Add scheduler for scraping so bucket stays up to date.
	- Concurrency for faster scraping if possible
	- Ranking --> Events --> Stats scraping
*/
