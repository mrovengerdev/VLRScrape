# VLRScrape
A Go (golang) scraper for the vlr.gg website (which loads its credentials from a .env file using the godotenv library).


## Functionality
- Scheduler
   - Can uncomment scheduler to run this scraper at 6:00AM, 12:00PM, 6:00PM, and 12:00AM every day.
- Scrape VLR forum threads.  
   - Specify in pageParser argument the header to decide the time table you want to scrape from.  
- Scrape VLR upcoming matches.  
   - Specify in pageParser argument the header to decide the time table you want to scrape from.  
- Scrape VLR leaderboards.
   - Note: Region does not specify the league they play in. Just their origin region.
- Upload the retrieved data to a specified S3 bucket.  
   - Bucket destination stated in .env file.  
- REST API
   - Following the retrieval of all endpoints, a REST API is enabled which allows for the retrieval of any folder through the base endpoint http://localhost:8080/.


## Installation
go get github.com/mrovengerdev/vlrscrape


## Usage
Create a .env file which stores your environmental variables like so:

AWS_ACCESS_KEY_ID= enter-your-acess-key-id-here  
AWS_SECRET_KEY= enter-your-aws-secret-key-here  
AWS_S3_BUCKET= enter-your-aws-bucket-name-here  

To run the program:  
- go run .


## Future Features
- Improve error handling messages.
- Events --> Stats scraping.
- Optimization/Refactoring for dateScraper
- Remove need for TrimSpace by adding .AttrOr


# Relevant/Contact Information
**Name:** Maxwell Rovenger  
**Github Username:** mrovengerdev  
**Email:** rovenger.max@gmail.com  