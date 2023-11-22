package scraper

import (
	"fmt"
	parsing "go-scraper/parsing"
	"go-scraper/repository"
	"strings"

	"github.com/gocolly/colly"
)

func StartScraping() {
	// Instantiate default collector
	c := colly.NewCollector(
		colly.AllowedDomains("hudebnibazar.cz"),
		colly.MaxDepth(2),
		colly.Async(true),
	)
	c.Limit(&colly.LimitRule{DomainGlob: "*", Parallelism: 10})

	c.OnHTML("td[class=content]", func(e *colly.HTMLElement) {
		if e.ChildText("div.InzeratText.dont-break-out") == "" {
			// this is not listing, get the hell out of here
			c.OnHTML("a[href]", func(e *colly.HTMLElement) {
				link := e.Attr("href")
				if isRelevantUrl(link) {
					c.Visit(e.Request.AbsoluteURL(link))
				}
			})
			return
		}

		// MAP USER
		user, err := parsing.ParseHTMLToUser(e)
		if err != nil {
			fmt.Printf("Error while mapping properties to user: %s\n", err)
		}

		if !repository.IsUserIgnored(user) {
			var usrErr error
			go func() {
				usrErr = repository.UpsertUser(user)
			}()

			if usrErr != nil {
				fmt.Printf("Error while saving user: %s\n", err)
			}
			fmt.Printf("User upsert done: %s\n", user.UserName)
		}

		// MAP LISTING
		listing, err := parsing.ParseHTMLToListing(e, user.UserName)
		if err != nil {
			fmt.Printf("Error while mapping properties to listing: %s\n", err)
		}
		fmt.Printf("Listing upsert done: %s\n", listing.ID)

		if repository.IsListingIgnored(listing) {
			fmt.Printf("Listing already in cache: %s\n", e.Request.URL.Path)
			return
		}

		var listingError error
		go func() {
			listingError = repository.UpsertListing(listing)
		}()
		if listingError != nil {
			fmt.Printf("Error while saving listing: %s\n", err)
		}
	})

	c.OnHTML("div[class=InzeratNadpis]", func(e *colly.HTMLElement) {
		link, _ := e.DOM.Parent().Attr("href")
		if isRelevantUrl(link) {
			c.Visit(e.Request.AbsoluteURL(link))
		}
	})

	// Before making a request print "Visiting ..."
	c.OnRequest(func(r *colly.Request) {
		fmt.Println("Visiting", r.URL.String())
	})

	c.Visit("https://hudebnibazar.cz/")
	c.Wait()
}

func isRelevantUrl(url string) bool {
	if strings.Contains(url, "uzivatel") {
		return false
	}

	if strings.Contains(url, "img_cache") {
		return false
	}

	if strings.Contains(url, "inzerat") {
		return false
	}

	return true
}
