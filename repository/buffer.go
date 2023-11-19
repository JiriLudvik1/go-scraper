package repository

import (
	"go-scraper/models"
	"time"
)

var buffer []models.Listing = make([]models.Listing, 10)
var lastIndex int = 0
var lastUpdate time.Time

func Enqueue(listing models.Listing) {
	if lastIndex == 10 {
		FlushSave()
	}

	buffer[lastIndex] = listing
	lastIndex++
	lastUpdate = time.Now()
}

func FlushSave() {
	err := SaveListings(buffer)
	if err != nil {
		panic(err)
	}

	resetBuffer(&buffer)
	lastIndex = 0
}

func StartFlushTicker() {
	// define an interval and the ticker for this interval
	interval := time.Duration(1) * time.Second
	// create a new Ticker
	tk := time.NewTicker(interval)
	// start the ticker by constructing a loop
	for range tk.C {
		if shouldFlushBuffer() {
			FlushSave()
		}
	}
}

func resetBuffer(buff *[]models.Listing) {
	for i := 0; i < len(*buff); i++ {
		(*buff)[i] = models.Listing{}
	}
}

func shouldFlushBuffer() bool {
	if buffer[0] == (models.Listing{}) {
		return false
	}

	return time.Since(lastUpdate) > time.Duration(5)*time.Second
}
