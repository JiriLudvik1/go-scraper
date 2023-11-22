package parsing

import (
	"fmt"
	"go-scraper/models"
	"strconv"
	"strings"
	"time"

	"github.com/gocolly/colly"
)

func ParseHTMLToListing(e *colly.HTMLElement, userName string) (*models.Listing, error) {
	headerText := e.ChildText("h1")
	rawPrice := e.ChildText("div[class=InzeratCena]")
	formattedPrice, err := formatPrice(rawPrice)
	if err != nil {
		fmt.Printf("Error while parsing price: %s\n", err)
		return nil, err
	}
	link := e.Request.URL.Path
	body := e.ChildText("div.InzeratText.dont-break-out")

	newListing := models.Listing{
		ID:        getIdFromUrl(link),
		Title:     headerText,
		Price:     formattedPrice,
		Link:      "https://hudebnibazar.cz" + link,
		Intent:    getIntentFromHTML(e),
		DateFound: time.Now(),
		Views:     getViewsFromHTML(e),
		Username:  userName,
		Body:      body,
	}
	return &newListing, nil
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
