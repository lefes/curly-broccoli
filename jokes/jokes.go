package jokes

import (
	"net/http"

	"github.com/PuerkitoBio/goquery"
)

const (
	htmlURL = "https://baneks.ru/random"
)

// GetJoke fetches from HTML page and returns a joke
func GetJoke() (string, error) {
	client := http.Client{}
	req, err := http.NewRequest("GET", htmlURL, nil)
	if err != nil {
		return "", err
	}
	resp, err := client.Do(req)
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()

	// Parse HTML and find tag with joke by CSS selector
	doc, err := goquery.NewDocumentFromReader(resp.Body)
	// handle err
	if err != nil {
		return "", err
	}

	joke := doc.Find(".anek-view > article:nth-child(2) > p:nth-child(2)").Text()

	return joke, nil
}
