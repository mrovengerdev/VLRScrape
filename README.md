# VLRSCRAPE
A Go (golang) scraper for the vlr.gg website (which loads its credentials from a .env file using the godotenv library).


## Functionality
 - Scrape VLR forum threads.  
    - Specify in pageParser argument the header to decide the time table you want to scrape from.  
 - Scrape VLR upcoming matches.  
    - Specify in pageParser argument the header to decide the time table you want to scrape from.  
 - Upload the retrieved data to a specified S3 bucket.  
    - Bucket destination stated in .env file.  


## Installation
go get github.com/mrovengerdev/vlrscrape


## Usage
Create a .env file which stores your environmental variables like so:

AWS_ACCESS_KEY_ID= enter-your-acess-key-id-here  
AWS_SECRET_KEY= enter-your-aws-secret-key-here  
AWS_S3_BUCKET= enter-your-aws-bucket-name-here  

To run the program:  
- go run .


# Relevant/Contact Information
**Name:** Maxwell Rovenger  
**Github Username:** mrovengerdev  
**Email:** rovenger.max@gmail.com  

# Future Features
- Make error handling unique instead of boiler plate.
- Add scheduler for scraping so bucket stays up to date.  
- Ranking --> Events --> Stats scraping  
- Optimization/Refactoring for dateScraper  