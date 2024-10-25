package scrapetools

import (
	"log"
	"os"
	"path/filepath"
	"regexp"
	"strconv"
	"strings"
)

// Checks if string is an int
func IsInt(teamName string) bool {
	if _, err := strconv.Atoi(teamName); err == nil {
		return true
	} else {
		return false
	}
}

func Filter(input string, filterValue string, replaceValue string) string {
	// Use regex to match one or more tab characters
	re := regexp.MustCompile(filterValue + `+`)
	// Replace all matches with a single space
	return re.ReplaceAllString(input, replaceValue)
}

// Properly joins all of the JSON files into one via handling excessive brackets.
func FileFix(fileName string) {
	file, err := os.ReadFile(fileName)
	if err != nil {
		log.Fatalf("Error: %v", err)
	}

	allText := string(file)
	allText = strings.ReplaceAll(allText, "    }\n][\n    {", "    },\n    {")

	err = os.WriteFile(fileName, []byte(allText), 0644)
	if err != nil {
		log.Fatalf("Error: %v", err)
	}
}

// Creates output folder to store JSON files
func CreateOutputDirectory() {
	outputPath, err := os.Getwd()
	if err != nil {
		log.Fatalf("Error: %v", err)
	}
	// Create the "output" directory
	outputDir := filepath.Join(outputPath, "output")
	err = os.MkdirAll(outputDir, 0755) // Creates directory if it doesn't exist
	if err != nil {
		log.Fatalf("Error: %v", err)
	}
}
