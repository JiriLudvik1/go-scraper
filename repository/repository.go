package repository

import (
	"database/sql"
	"go-scraper/models"
	"os"

	_ "github.com/mattn/go-sqlite3"
)

var connString = "./listings.db"
var existingIds = []string{}

func InitializeDB() error {
	_, err := os.Stat(connString)
	if os.IsNotExist(err) {
		f, err := os.Create(connString)
		if err != nil {
			panic(err)
		}
		defer f.Close()
	}

	db := openDbConnection()
	defer db.Close()
	_, err = db.Exec("CREATE TABLE IF NOT EXISTS listings (id TEXT PRIMARY KEY, title TEXT, price REAL, link TEXT, intent TEXT, date_found DATETIME)")
	if err != nil {
		return err
	}

	rows, err := db.Query("SELECT id FROM listings")
	if err != nil {
		return err
	}

	for rows.Next() {
		var id string
		err = rows.Scan(&id)
		if err != nil {
			return err
		}
		existingIds = append(existingIds, id)
	}
	return nil
}

func SaveListing(listing *models.Listing) error {
	db := openDbConnection()
	defer db.Close()
	_, err := db.Exec("INSERT INTO listings (id, title, price, link, intent, date_found) VALUES (?, ?, ?, ?, ?, ?)",
		listing.ID, listing.Title, listing.Price, listing.Link, listing.Intent, listing.DateFound)
	return err
}

func SaveListings(listings []models.Listing) error {
	db := openDbConnection()
	defer db.Close()

	// Begin a transaction
	tx, err := db.Begin()
	if err != nil {
		return err
	}

	// Rollback the transaction in case of an error and close the database
	defer func() {
		if err != nil {
			tx.Rollback()
			return
		}
		tx.Commit()
	}()

	// Prepare the INSERT statement outside the loop for better performance
	stmt, err := tx.Prepare("INSERT INTO listings (id, title, price, link, intent, date_found) VALUES (?, ?, ?, ?, ?, ?)")
	if err != nil {
		return err
	}
	defer stmt.Close()

	var addedIds []string
	// Iterate through the listings and execute the INSERT statement for each one
	for _, listing := range listings {
		if listing.ID == "" {
			continue
		}
		if existsInSlice(existingIds, listing.ID) {
			continue
		}

		_, err := stmt.Exec(listing.ID, listing.Title, listing.Price, listing.Link, listing.Intent, listing.DateFound)
		if err != nil {
			if err.Error() == "UNIQUE constraint failed: listings.id" {
				continue
			}

			return err
		}
		addedIds = append(addedIds, listing.ID)
	}

	existingIds = append(existingIds, addedIds...)
	return nil
}

func existsInSlice(slice []string, id string) bool {
	for _, item := range slice {
		if item == id {
			return true
		}
	}

	return false
}

func openDbConnection() *sql.DB {
	db, err := sql.Open("sqlite3", connString)
	if err != nil {
		panic(err)
	}

	return db
}
