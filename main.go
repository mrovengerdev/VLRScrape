package main

import (
	"github.com/mrovengerdev/vlrscrape/s3port"
	"github.com/mrovengerdev/vlrscrape/scrape"
)

func main() {
	// Creates output folder for JSON files.
	scrape.CreateOutputDirectory()

	// Scrape from VLR.gg threads.
	scrape.PageParser("https://www.vlr.gg/threads", "/?t=1w", "outputThreads")

	// Scrape from VLR.gg matches.
	scrape.PageParser("https://www.vlr.gg/matches", "/?", "outputMatches")

	// Upload output files to Amazon S3 bucket: "vlr-scrape".
	s3port.Upload()
}
