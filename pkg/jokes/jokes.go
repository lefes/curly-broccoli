package jokes

import (
	"errors"
	"net/http"
	"time"

	"github.com/PuerkitoBio/goquery"
)

const (
	htmlURL = "https://baneks.ru/random"
)

// GetJoke fetches from HTML page and returns a joke
func GetJoke() (string, error) {
	client := http.Client{Timeout: 10 * time.Second}
	req, err := http.NewRequest("GET", htmlURL, nil)
	if err != nil {
		return "", err
	}
	req.Header.Set("User-Agent", "curly-broccoli-bot/1.0 (+https://github.com/lefes/curly-broccoli)")

	resp, err := client.Do(req)
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return "", errors.New("non-200 response from jokes source")
	}

	// Parse HTML and find tag with joke by CSS selector
	doc, err := goquery.NewDocumentFromReader(resp.Body)
	if err != nil {
		return "", err
	}

	// Primary selector used historically
	joke := doc.Find(".anek-view > article:nth-child(2) > p:nth-child(2)").Text()
	if joke == "" {
		// Fallback: try first article paragraph
		joke = doc.Find("article p").First().Text()
	}
	if joke == "" {
		return "", errors.New("failed to parse joke from page")
	}
	return joke, nil
}
