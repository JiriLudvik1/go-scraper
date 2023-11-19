package scraper

import (
	"fmt"
	"go-scraper/models"
	"go-scraper/repository"
	"strconv"
	"strings"
	"time"

	"github.com/gocolly/colly"
)

func StartScraping() {
	// Instantiate default collector
	c := colly.NewCollector(
		colly.AllowedDomains("hudebnibazar.cz"),
	)

	c.OnHTML("div[class=InzeratBody]", func(e *colly.HTMLElement) {
		listing, err := mapPropertiesToListing(e)
		if err != nil {
			fmt.Printf("Error while mapping properties to listing: %s\n", err)
		}
		repository.Enqueue(listing)
	})

	c.OnHTML("a[href]", func(e *colly.HTMLElement) {
		link := e.Attr("href")
		c.Visit(e.Request.AbsoluteURL(link))
	})

	// Before making a request print "Visiting ..."
	c.OnRequest(func(r *colly.Request) {
		fmt.Println("Visiting", r.URL.String())
	})

	c.Visit("https://hudebnibazar.cz/")
}

func mapPropertiesToListing(e *colly.HTMLElement) (models.Listing, error) {
	headerText := e.ChildText("b")
	fmt.Printf("Header found: %q\n", headerText)

	link := e.DOM.Find("a").Last().AttrOr("href", "")
	link = "https://hudebnibazar.cz" + link

	rawPrice := e.ChildText("div[class=InzeratCena]")
	formattedPrice, err := formatPrice(rawPrice)
	if err != nil {
		fmt.Printf("Error while parsing price: %s\n", err)
		return models.Listing{}, err
	}

	newListing := models.Listing{
		ID:        getIdFromUrl(link),
		Title:     e.ChildText("div[class=InzeratNadpis]"),
		Price:     formattedPrice,
		Link:      link,
		Intent:    getIntentFromHTML(e),
		DateFound: time.Now(),
	}
	return newListing, nil
}

func getIntentFromHTML(e *colly.HTMLElement) string {
	res := e.ChildText("div.label-nabidka")
	if res != "" {
		return res
	}

	res = e.ChildText("div.label-poptavka")
	if res != "" {
		return res
	}

	res = e.ChildText("div.label-ruzne")
	if res != "" {
		return res
	}

	return ""
}

func getIdFromUrl(url string) string {
	splits := strings.Split(url, "/")
	return splits[len(splits)-2]
}

func formatPrice(price string) (float32, error) {
	if price == "" {
		return 0, nil
	}

	firstSplit := strings.Split(price, " ")[0]
	if firstSplit == "cena" {
		return 0, nil
	}

	splits := strings.Split(price, " Kč")
	cleanedPrice := strings.ReplaceAll(splits[0], " ", "")

	res, err := strconv.ParseFloat(cleanedPrice, 32)
	if err != nil {
		return 0, err
	}

	return (float32)(res), nil
}
