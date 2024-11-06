package main

import (
	"github.com/mrovengerdev/vlrscrape/restAPI"
	"github.com/mrovengerdev/vlrscrape/s3port"
	"github.com/mrovengerdev/vlrscrape/scrape"
	"github.com/mrovengerdev/vlrscrape/scrapetools"
)

func main() {

	// Creates output folder for JSON files.
	scrapetools.CreateOutputDirectory()

	// Scrape from VLR.gg threads. Change 2nd argument to specify time frame.
	scrape.PageParser("https://www.vlr.gg/threads", "/?t=1w", "outputThreads")

	// Scrape from VLR.gg matches.
	scrape.PageParser("https://www.vlr.gg/matches", "/?", "outputMatches")

	// Scrape from VLR.gg teams.
	prepDocument := scrape.ScrapePrep("https://www.vlr.gg/teams" + "")
	scrape.RankingScrape(prepDocument)

	// Upload output files to Amazon S3 bucket: "vlr-scrape".
	s3port.Upload()

	// Enables REST API endpoint throuhg localhost.
	restAPI.Get()

	// Scheduled version of the main method.
	// Scheduler runs the program at 6:00AM, 12:00PM, 6:00PM, and 12:00AM.
	// c := cron.New()
	// c.AddFunc("0 */6 * * *", func() {
	// 	// Creates output folder for JSON files.
	// 	scrape.CreateOutputDirectory()

	// 	// Scrape from VLR.gg threads.
	// 	scrape.PageParser("https://www.vlr.gg/threads", "/?t=1w", "outputThreads")

	// 	// Scrape from VLR.gg matches.
	// 	scrape.PageParser("https://www.vlr.gg/matches", "/?", "outputMatches")

	// 	// Upload output files to Amazon S3 bucket: "vlr-scrape".
	// 	s3port.Upload()
	// })
	// c.Start()
}
