package main

import (
	"encoding/json"
	"fmt"
	"net/http"
	"os"
	"path/filepath"
	"strconv"
	"strings"

	"github.com/PuerkitoBio/goquery"
)

const base_url = "https://www.vlr.gg"

// Makes connection to scraping destination and returns document for parsing
func threadPrep(url string) *goquery.Document {
	response, err := http.Get(url)
	check(err)
	defer response.Body.Close()

	if response.StatusCode == 200 {
		fmt.Printf("Success: %d at %s \n", response.StatusCode, url)
	} else {
		fmt.Println("Error: ", err)
	}

	doc, err := goquery.NewDocumentFromReader(response.Body)
	check(err)
	return doc
}

// Checks if string is an int
func isInt(teamName string) bool {
	if _, err := strconv.Atoi(teamName); err == nil {
		return true
	} else {
		return false
	}
}

// Creates output folder to store JSON files
func createOutputDirectory() {
	outputPath, err := os.Getwd()
	check(err)
	// Create the "output" directory
	outputDir := filepath.Join(outputPath, "output")
	err = os.MkdirAll(outputDir, 0755) // Creates directory if it doesn't exist
	check(err)
}

// Scrape threads from vlr.gg/threads
func threadScrape(doc *goquery.Document) {

	type Thread struct {
		ID               int    `json:"id"`
		Title            string `json:"title"`
		URL              string `json:"url"`
		FragCount        int    `json:"frag_count"`
		DatePublished    string `json:"date_published"`
		DatePublishedAgo string `json:"date_published_ago"`
		CommentCount     int    `json:"comment_count"`
	}

	var threads []Thread

	doc.Find("div.thread.wf-module-item.mod-color.mod-left.mod-bg-after-.unread").Each(func(index int, item *goquery.Selection) { // Two wf-cards so double for-loop required

		// Text Processing (string to int)
		tempFragCount, err := strconv.Atoi(strings.TrimSpace(item.Find("span.frag-count").Text()))
		check(err)
		tempID, err := strconv.Atoi(item.Find("div.block.frag.frag-container.noselect.neutral").AttrOr("data-thread-id", ""))
		check(err)

		//Text Processing (trimming and replacing)
		tempCommentCount := strings.TrimSpace(item.Find("span.post-count").Text())
		tempCommentCount = strings.ReplaceAll(tempCommentCount, "\t\t\t\t\t\t\t\t\t\t\t\t\t", " ")
		commentNum, err := strconv.Atoi(strings.Split(tempCommentCount, " ")[0])
		check(err)

		thread := Thread{
			ID:               tempID,
			Title:            strings.TrimSpace(item.Find(".thread-item-header-title").Text()),
			URL:              base_url + item.Find(".thread-item-header-title").AttrOr("href", ""),
			FragCount:        tempFragCount,
			DatePublished:    strings.TrimSpace(item.Find("span.date-full.hide").Text()),
			DatePublishedAgo: strings.TrimSpace(item.Find("span.js-date-toggle.date-eta").Text()),
			CommentCount:     commentNum,
		}
		threads = append(threads, thread)
		for i := 1; i < len(threads); i++ {

		}
	})

	// Converts data format to JSON
	jsonData, err := json.MarshalIndent(threads, "", "    ")
	check(err)

	err = os.WriteFile("output/outputThreads.json", jsonData, 0644)
	check(err)

	fmt.Println("Thread scrape complete.")
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
func matchScrape(doc *goquery.Document) {

	type Match struct {
		ID         int    `json:"id"`
		URL        string `json:"url"`
		Tournament string `json:"tournament"`
		Team1      string `json:"team1"`
		Team2      string `json:"team2"`
		Date       string `json:"date"`
		MatchTime  string `json:"match_time"`
		TimeUntil  string `json:"time_until"` // Time until match
	}

	var matches []Match

	doc.Find("a[class*='mod-color']").Each(func(index int, item *goquery.Selection) {

		// Retrieve ID from URL through string parsing
		tempID := item.AttrOr("href", "")
		tempID = strings.Split(tempID, "/")[1]
		intTempID, err := strconv.Atoi(tempID)
		check(err)

		// Retrieve team names and trim spaces & \t
		tempTeam := strings.TrimSpace(item.Find("div.match-item-vs-team").Text())
		tempTeam = strings.ReplaceAll(tempTeam, "\t", "")

		// Assigns team names to team1 and team2
		tempTeam1 := strings.Split(tempTeam, "\n\n\n")[0]
		tempTeam2 := strings.Split(tempTeam, "\n\n\n")[3]

		// Checks if team2's name is team1's score. If so, then the team name is the second element.
		if isInt(tempTeam2) {
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
		dateDoc := threadPrep(matchURL)

		match := Match{
			Tournament: strings.ReplaceAll(strings.TrimSpace(item.Find("div.match-item-event-series.text-of").Text()), "–", " "),
			ID:         intTempID,
			URL:        matchURL,
			Team1:      tempTeam1,
			Team2:      tempTeam2,
			Date:       dateScrape(dateDoc),
			MatchTime:  strings.TrimSpace(item.Find("div.match-item-time").Text()),
			TimeUntil:  tempTimeUntil,
		}
		matches = append(matches, match)

	})

	// Converts data format to JSON
	jsonData, err := json.MarshalIndent(matches, "", "    ")
	check(err)

	err = os.WriteFile("output/outputMatches.json", jsonData, 0644)
	check(err)

	fmt.Println("Match scrape complete.")
}

// Test output
// fmt.Println(string(jsonData))
