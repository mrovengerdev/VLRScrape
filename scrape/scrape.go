package scrape

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"
	"strconv"
	"strings"
	"time"

	"github.com/PuerkitoBio/goquery"
	"github.com/mrovengerdev/vlrscrape/paginator"
	"github.com/mrovengerdev/vlrscrape/scrapetools"
)

// Struct Catalog
type Thread struct {
	ID               int    `json:"id"`
	Title            string `json:"title"`
	ThreadURL        string `json:"thread_url"`
	FragCount        int    `json:"frag_count"`
	DatePublished    string `json:"date_published"`
	DatePublishedAgo string `json:"date_published_ago"`
	CommentCount     int    `json:"comment_count"`
}

type Match struct {
	ID             int    `json:"id"`
	MatchURL       string `json:"match_url"`
	Tournament     string `json:"tournament"`
	Team1          string `json:"team1"`
	Team2          string `json:"team2"`
	Date           string `json:"date"`
	MatchTime      string `json:"match_time"`
	TimeUntilMatch string `json:"time_until_match"` // Time until match
}

type Ranking struct {
	Rank     int    `json:"rank"`
	Region   string `json:"region"`
	TeamName string `json:"team_name"`
	ELO      int    `json:"elo"`
	TeamURL  string `json:"team_url"`
}

const base_url = "https://www.vlr.gg"

// Makes connection to scraping destination and returns document for parsing.
func ScrapePrep(url string) *goquery.Document {
	response, err := http.Get(url)
	if err != nil {
		log.Fatalf("Error: %v", err)
	}
	defer response.Body.Close()

	if response.StatusCode == 200 {
		fmt.Printf("Success: %d at %s \n", response.StatusCode, url)
	} else {
		fmt.Println("Error: ", err)
	}

	doc, err := goquery.NewDocumentFromReader(response.Body)
	if err != nil {
		log.Fatalf("Error: %v", err)
	}
	return doc
}

// Scrape threads from vlr.gg/threads. Returns JSON data as []byte.
func threadScrape(currentPage int, doc *goquery.Document) []byte {

	var threads []Thread

	// Needs performance improvement:
	// Retrieves the first 3 threads only for the first page, since they repeat on every page.
	doc.Find("div.thread.wf-module-item.mod-color.mod-left.mod-bg-after-.unread").Each(func(index int, item *goquery.Selection) { // Two wf-cards so double for-loop required
		// Ignore first 3 posts since they are always the same
		if currentPage == 1 {
			// Upvote count processing (string to int)
			tempFragCount, err := strconv.Atoi(strings.TrimSpace(item.Find("span.frag-count").Text()))
			if err != nil {
				log.Fatalf("Error: %v", err)
			}
			tempID, err := strconv.Atoi(item.Find("div.block.frag.frag-container.noselect.neutral").AttrOr("data-thread-id", ""))
			if err != nil {
				log.Fatalf("Error: %v", err)
			}

			// Comment count processing (string to int)
			tempCommentCount := strings.TrimSpace(item.Find("span.post-count").Text())
			tempCommentCount = strings.ReplaceAll(tempCommentCount, "\t\t\t\t\t\t\t\t\t\t\t\t\t", " ")
			commentNum, err := strconv.Atoi(strings.Split(tempCommentCount, " ")[0])
			if err != nil {
				log.Fatalf("Error: %v", err)
			}

			thread := Thread{
				ID:               tempID,
				Title:            strings.TrimSpace(item.Find(".thread-item-header-title").Text()),
				ThreadURL:        base_url + item.Find(".thread-item-header-title").AttrOr("href", ""),
				FragCount:        tempFragCount,
				DatePublished:    strings.TrimSpace(item.Find("span.date-full.hide").Text()),
				DatePublishedAgo: strings.TrimSpace(item.Find("span.js-date-toggle.date-eta").Text()),
				CommentCount:     commentNum,
			}

			threads = append(threads, thread)
		} else {
			if index > 2 {
				// Upvote count processing (string to int)
				tempFragCount, err := strconv.Atoi(strings.TrimSpace(item.Find("span.frag-count").Text()))
				if err != nil {
					log.Fatalf("Error: %v", err)
				}
				tempID, err := strconv.Atoi(item.Find("div.block.frag.frag-container.noselect.neutral").AttrOr("data-thread-id", ""))
				if err != nil {
					log.Fatalf("Error: %v", err)
				}

				// Comment count processing (string to int)
				tempCommentCount := strings.TrimSpace(item.Find("span.post-count").Text())
				tempCommentCount = strings.ReplaceAll(tempCommentCount, "\t\t\t\t\t\t\t\t\t\t\t\t\t", " ")
				commentNum, err := strconv.Atoi(strings.Split(tempCommentCount, " ")[0])
				if err != nil {
					log.Fatalf("Error: %v", err)
				}

				thread := Thread{
					ID:               tempID,
					Title:            strings.TrimSpace(item.Find(".thread-item-header-title").Text()),
					ThreadURL:        base_url + item.Find(".thread-item-header-title").AttrOr("href", ""),
					FragCount:        tempFragCount,
					DatePublished:    strings.TrimSpace(item.Find("span.date-full.hide").Text()),
					DatePublishedAgo: strings.TrimSpace(item.Find("span.js-date-toggle.date-eta").Text()),
					CommentCount:     commentNum,
				}

				threads = append(threads, thread)
			}
		}
	})

	// Converts data format to JSON
	jsonData, err := json.MarshalIndent(threads, "", "    ")
	if err != nil {
		log.Fatalf("Error: %v", err)
	}

	fmt.Println("Thread scrape complete.")

	return jsonData
}

// Retrieves match dates for matchScrape
// TODO: Refactor:
// Currently, retrieving date requires connecting to every single match's match page.
func dateScrape(doc *goquery.Document) string {
	currentDate := strings.TrimSpace(doc.Find("div.moment-tz-convert").Text())
	currentDate = strings.ReplaceAll(currentDate, "\t\t\t\t\n\n\t\t\t\t\t\t\t\n\t\t\t\t\t\t", " ")
	return currentDate
}

// Scrape matches from vlr.gg/matches
func matchScrape(doc *goquery.Document) []byte {

	var matches []Match

	doc.Find("a[class*='mod-color']").Each(func(index int, item *goquery.Selection) {

		// Retrieve ID from URL through string parsing
		tempID := item.AttrOr("href", "")
		tempID = strings.Split(tempID, "/")[1]
		intTempID, err := strconv.Atoi(tempID)
		if err != nil {
			log.Fatalf("Error: %v", err)
		}

		// Retrieve team names and trim spaces & \t
		tempTeam := strings.TrimSpace(item.Find("div.match-item-vs-team").Text())
		tempTeam = strings.ReplaceAll(tempTeam, "\t", "")

		// Assigns team names to team1 and team2
		tempTeam1 := strings.Split(tempTeam, "\n\n\n")[0]
		tempTeam2 := strings.Split(tempTeam, "\n\n\n")[3]

		// Checks if team2's name is team1's score. If so, then the team name is the second element.
		if scrapetools.IsInt(tempTeam2) {
			tempTeam2 = strings.Split(tempTeam, "\n\n\n")[2]
			tempTeam2 = strings.ReplaceAll(tempTeam2, "\n", "")
		}

		// If no time until match, then it is live
		tempTimeUntil := strings.TrimSpace(item.Find("div.ml-eta").Text())
		if tempTimeUntil == "" {
			tempTimeUntil = "Live"
		}

		// Retrieve current match URL for date scraping
		matchURL := base_url + item.AttrOr("href", "")

		// For each match, got to the match page, and retrieve the match date at the top right.
		dateDoc := ScrapePrep(matchURL)

		match := Match{
			Tournament:     strings.ReplaceAll(strings.TrimSpace(item.Find("div.match-item-event-series.text-of").Text()), "–", " "),
			ID:             intTempID,
			MatchURL:       matchURL,
			Team1:          tempTeam1,
			Team2:          tempTeam2,
			Date:           dateScrape(dateDoc),
			MatchTime:      strings.TrimSpace(item.Find("div.match-item-time").Text()),
			TimeUntilMatch: tempTimeUntil,
		}
		matches = append(matches, match)
	})

	// Converts data format to JSON
	jsonData, err := json.MarshalIndent(matches, "", "    ")
	if err != nil {
		log.Fatalf("Error: %v", err)
	}

	fmt.Println("Match scrape complete.")

	return jsonData
}

// Scrape leaderboard rankings and team info from vlr.gg/rankings
func rankingScrape(doc *goquery.Document, region string) {

	var rankings []Ranking

	doc.Find("div.rank-item.wf-card.fc-flex").Each(func(index int, item *goquery.Selection) {

		tempRank, err := strconv.Atoi(strings.TrimSpace(item.Find("div.rank-item-rank-num").Text()))
		if err != nil {
			log.Fatalf("Error: %v", err)
		}

		tempELO, err := strconv.Atoi(item.Find("div.rank-item-rating").AttrOr("data-sort-value", ""))
		if err != nil {
			log.Fatalf("Error: %v", err)
		}

		ranking := Ranking{
			Rank:     tempRank,
			TeamName: item.Find("a.rank-item-team.fc-flex").AttrOr("data-sort-value", ""),
			ELO:      tempELO,
			Region:   scrapetools.Filter(region, "-", " "),
			TeamURL:  base_url + item.Find("a.rank-item-team.fc-flex").AttrOr("href", ""),
		}

		rankings = append(rankings, ranking)
	})

	// Converts data format to JSON
	jsonData, err := json.MarshalIndent(rankings, "", "    ")
	if err != nil {
		log.Fatalf("Error: %v", err)
	}

	// Writes JSON data into new/existing JSON file.
	os.WriteFile("output/ranking/output"+region+"Rankings"+".json", jsonData, 0644)
}

// Scrapes the rankings from all regions by using the rankingScrape for each region.
func AllRankingScrape(doc *goquery.Document) {
	doc.Find("a.wf-nav-item.mod-collapsible").Each(func(index int, item *goquery.Selection) {

		// Retrieve the region name and filter out any unnecessary characters.
		region := strings.TrimSpace(item.Find("span.normal").Text())
		region = scrapetools.Filter(region, " ", "-")

		fmt.Println(region)

		// Use rankingScrape at the region URL.
		if region != "World" && region != "" {
			rankingURL := base_url + "/rankings/" + region
			rankingDoc := ScrapePrep(rankingURL)
			rankingScrape(rankingDoc, region)
		}
	})

	fmt.Println("Ranking scrape complete.")
}

// Retrieves the number of the last page of threads containing unique threads.
// Only way since pages out of bounds will still contain the top 4 posts.
// Verify that the doc.Find() location works for future scrapes.
func findLastPage(doc *goquery.Document) int {
	lastPage := ""
	doc.Find("a.btn.mod-page").Each(func(index int, item *goquery.Selection) {
		lastPage = item.Text()
	})
	lastPageInt, err := strconv.Atoi(lastPage)
	if err != nil {
		log.Fatalf("Error: %v", err)
	}

	return lastPageInt
}

// Conducts scraping for the total number of pages available to the given section_url.
// For every new scrape added, the switch statement must be edited to cover it.
func PageParser(section_url string, header string, outputFileName string, paginator *paginator.Paginator) {

	// Stores page of scraped data per index
	var totalScrape = [][]byte{}
	// Stores the current page of scraped data
	var currentPageScrape []byte

	currentPage := 1
	pageExists := true

	// The class that gives the last page changes when scraping the last page. So before looping, it must be retrieved.
	prepDocument := ScrapePrep(section_url + header)
	lastPage := findLastPage(prepDocument)

	// For every page, scrape the data and append it to the totalScrape slice.
	for pageExists {
		if currentPage <= lastPage {
			url := fmt.Sprintf("%s%s&page=%d", section_url, header, currentPage)

			// Wait for permission from the limiter
			if err := paginator.Limiter.Wait(paginator.Context); err != nil {
				fmt.Println("Rate limiter context canceled or timed out:", err)
				break
			}

			document := ScrapePrep(url)

			switch section_url {
			case "https://www.vlr.gg/threads":
				currentPageScrape = threadScrape(currentPage, document)
			case "https://www.vlr.gg/matches":
				currentPageScrape = matchScrape(document)
			default:
				fmt.Println("Invalid base URL.")
			}

			// Join all pages of data.
			totalScrape = append(totalScrape, currentPageScrape)
			currentPage++

		} else {
			pageExists = false
		}
	}

	// Create the output file.
	timeStamp := time.Now().Format("2006-01-02_15-04-05")
	file, err := os.Create("output/" + outputFileName + "_" + timeStamp + ".json")
	if err != nil {
		log.Fatal(err)
	}
	defer file.Close()

	// Write the JSON data to a file.
	for i := 0; i < len(totalScrape); i++ {
		if _, err := file.Write(totalScrape[i]); err != nil {
			log.Fatal(err)
		}
	}

	scrapetools.FileFix("output/" + outputFileName + "_" + timeStamp + ".json")
}
