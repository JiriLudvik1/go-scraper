package main

import (
	"go-scraper/repository"
	scraper "go-scraper/scrapers"
)

func main() {
	//DB Stuff
	repository.InitializeIgnoreCache()
	err := repository.InitializeDB()
	if err != nil {
		panic(err)
	}
	go repository.StartFlushTicker()

	//Scraping stuff
	scraper.StartScraping()
}
