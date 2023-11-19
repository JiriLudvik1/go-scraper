package models

import "time"

type Listing struct {
	ID        string
	Title     string
	Price     float32
	Link      string
	Intent    string
	DateFound time.Time
	Views     int
}
