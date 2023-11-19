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

	c.OnHTML("td[class=content]", func(e *colly.HTMLElement) {
		if e.ChildText("div.InzeratText.dont-break-out") == "" {
			// this is not listing, get the hell out of here
			return
		}

		listing, err := mapPropertiesToListing(e)
		if err != nil {
			fmt.Printf("Error while mapping properties to listing: %s\n", err)
		}
		repository.Enqueue(listing)
	})

	c.OnHTML("div[class=InzeratNadpis]", func(e *colly.HTMLElement) {
		link, _ := e.DOM.Parent().Attr("href")
		if isRelevantUrl(link) {
			c.Visit(e.Request.AbsoluteURL(link))
		}
	})

	c.OnHTML("a[href]", func(e *colly.HTMLElement) {
		link := e.Attr("href")
		if isRelevantUrl(link) {
			c.Visit(e.Request.AbsoluteURL(link))
		}
	})

	// Before making a request print "Visiting ..."
	c.OnRequest(func(r *colly.Request) {
		fmt.Println("Visiting", r.URL.String())
	})

	c.Visit("https://hudebnibazar.cz/")
}

func mapPropertiesToListing(e *colly.HTMLElement) (models.Listing, error) {
	headerText := e.ChildText("h1")
	fmt.Printf("Header found: %q\n", headerText)

	rawPrice := e.ChildText("div[class=InzeratCena]")
	formattedPrice, err := formatPrice(rawPrice)
	if err != nil {
		fmt.Printf("Error while parsing price: %s\n", err)
		return models.Listing{}, err
	}
	link := e.Request.URL.Path

	newListing := models.Listing{
		ID:        getIdFromUrl(link),
		Title:     headerText,
		Price:     formattedPrice,
		Link:      "https://hudebnibazar.cz" + link,
		Intent:    getIntentFromHTML(e),
		DateFound: time.Now(),
		Views:     getViewsFromHTML(e),
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

	firstPart := strings.Split(price, " Kč")[0]
	soloPrice := strings.Split(firstPart, "Cena: ")[1]
	cleanedPrice := strings.ReplaceAll(soloPrice, " ", "")

	res, err := strconv.ParseFloat(cleanedPrice, 32)
	if err != nil {
		return 0, err
	}

	return (float32)(res), nil
}

func getViewsFromHTML(e *colly.HTMLElement) int {
	var views int

	e.ForEach("div", func(_ int, s *colly.HTMLElement) {
		containsViews := strings.Contains(s.Text, "Zobrazeno")
		hasCorrectParent := s.DOM.Parent().AttrOr("class", "") == "InzeratBodyDetail"

		if !containsViews || !hasCorrectParent {
			return
		}

		text := s.Text
		numberStr := strings.TrimSpace(strings.Split(text, " ")[1])
		numberStr = strings.Replace(numberStr, "x", "", -1)
		views, _ = strconv.Atoi(numberStr)
	})

	return views
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
