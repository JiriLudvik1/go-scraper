package main

import (
	"go-scraper/repository"
	scraper "go-scraper/scrapers"
)

func main() {
	//DB Stuff
	err := repository.InitializeDB()
	if err != nil {
		panic(err)
	}
	go repository.StartFlushTicker()

	//Scraping stuff
	scraper.StartScraping()
}
