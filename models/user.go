package models

import "time"

type User struct {
	UserName  string
	FullName  string
	Locality  string
	Phone     string
	Rating    string
	DateFound time.Time
}
