# VLRScrape
A Go (golang) scraper for the vlr.gg website (which loads its credentials from a .env file using the godotenv library).


## Functionality
- Scheduler  
   - Can uncomment scheduler to run this scraper at 6:00AM, 12:00PM, 6:00PM, and 12:00AM every day.  
- Scrape VLR forum threads.  
   - Specify in pageParser argument the header to decide the time table you want to scrape from.  
- Scrape VLR upcoming matches.  
   - Specify in pageParser argument the header to decide the time table you want to scrape from.  
- Scrape VLR rankings per region.  
   - In beta for VLR so current endpoint may be deprecated. Works as of 11/6/2024.  
- Upload the retrieved data to a specified S3 bucket.  
   - Bucket destination stated in .env file.  
- REST API  
   - Following the retrieval of all endpoints, a REST API is enabled which allows for the retrieval of any folder through the base endpoint http://localhost:8080/.


## Installation
go get github.com/mrovengerdev/vlrscrape


## Usage
Create a .env file which stores your environmental variables like so:

AWS_VLR_ACCESS_KEY_ID=enter-your-acess-key-id-here  
AWS_VLR_SECRET_KEY=enter-your-aws-secret-key-here  
AWS_VLR_S3_BUCKET=enter-your-aws-bucket-name-here  
AWS_VLR_S3_REGION=enter-your-aws-region-here  

To run the program:  
- go run .


## TODO
- Improve error handling messages.
- Stats Scraping.
- Optimization/Refactoring for dateScraper (Currently has to access the link for all matches to find the date.)
- Refactor trimming to improve readability.


# Relevant/Contact Information
**Name:** Maxwell Rovenger  
**Github Username:** mrovengerdev  
**Email:** rovenger.max@gmail.com  