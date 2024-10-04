package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"
	"strconv"
	"strings"

	"github.com/PuerkitoBio/goquery"
)

func check(err error) {
	if err != nil {
		fmt.Println("Error occurred:", err)
		log.Fatal("Error occurred:", err)
		os.Exit(1)
	}
}

func threadPrep() *goquery.Document {
	response, err := http.Get("https://www.vlr.gg/threads")
	check(err)
	defer response.Body.Close()

	if response.StatusCode == 200 {
		fmt.Println("Success: ", response.StatusCode)
	} else {
		fmt.Println("Error: ", err)
	}

	doc, err := goquery.NewDocumentFromReader(response.Body)
	check(err)
	return doc
}

func threadScrape(doc *goquery.Document) {

	type Thread struct {
		Title            string `json:"title"`
		URL              string `json:"url"`
		FragCount        int    `json:"frag_count"`
		DatePublished    string `json:"date_published"`
		DatePublishedAgo string `json:"date_published_ago"`
		CommentCount     string `json:"comment_count"`
	}

	var threads []Thread

	base_url := "https://www.vlr.gg"
	doc.Find("div#thread-list.wf-card").Each(func(index int, item *goquery.Selection) {
		doc.Find("div.thread.wf-module-item.mod-color.mod-left.mod-bg-after-.unread").Each(func(index2 int, item2 *goquery.Selection) {

			// Converts string text to int
			tempFragCount, err := strconv.Atoi(strings.TrimSpace(item2.Find("span.frag-count").Text()))
			check(err)

			tempCommentCount := strings.TrimSpace(item2.Find("span.post-count").Text())

			thread := Thread{
				Title:            strings.TrimSpace(item2.Find(".thread-item-header-title").Text()),
				URL:              base_url + item2.Find(".thread-item-header-title").AttrOr("href", ""),
				FragCount:        tempFragCount,
				DatePublished:    strings.TrimSpace(item2.Find("span.date-full.hide").Text()),
				DatePublishedAgo: strings.TrimSpace(item2.Find("span.js-date-toggle.date-eta").Text()),
				CommentCount:     strings.ReplaceAll(tempCommentCount, "\t\t\t\t\t\t\t\t\t\t\t\t\t", " "),
			}
			threads = append(threads, thread)
		})
	})

	// Converts data format to JSON
	jsonData, err := json.MarshalIndent(threads, "", "    ")
	check(err)

	err = os.WriteFile("outputThreads.json", jsonData, 0644)
	check(err)

	// Test output
	// fmt.Println(string(jsonData))
}
