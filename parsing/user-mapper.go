package parsing

import (
	"go-scraper/models"
	"strings"
	"time"

	"github.com/gocolly/colly"
)

func ParseHTMLToUser(e *colly.HTMLElement) (*models.User, error) {
	userDiv := e.DOM.Find("div[class=user-right]")
	userName := userDiv.Find("span[class=muted-text]").Text()
	userName = strings.ReplaceAll(userName, "(", "")
	userName = strings.ReplaceAll(userName, ")", "")

	newUser := models.User{
		UserName:  userName,
		DateFound: time.Now(),
	}

	e.ForEach("div", func(_ int, s *colly.HTMLElement) {
		if strings.Contains(s.Text, "Jméno:") {
			splits := strings.Split(s.Text, ":")
			valueSplits := strings.Split(splits[1], " (")
			newUser.FullName = strings.TrimSpace(valueSplits[0])
		}

		if strings.Contains(s.Text, "Lokalita:") {
			splits := strings.Split(s.Text, ":")
			newUser.Locality = strings.TrimSpace(splits[1])
		}

		if strings.Contains(s.Text, "Telefon:") {
			splits := strings.Split(s.Text, ":")
			newUser.Phone = strings.TrimSpace(splits[1])
		}

		if strings.Contains(s.Text, "Hodnocení:") {
			splits := strings.Split(s.Text, ":")
			newUser.Rating = strings.TrimSpace(splits[1])
		}
	})

	return &newUser, nil
}
