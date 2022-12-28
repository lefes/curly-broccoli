package quotes

import (
	"encoding/json"
	"fmt"
	"net/http"
)

type Quote struct {
	Anime     string `json:"anime"`
	Character string `json:"character"`
	Quote     string `json:"quote"`
}

func New() *Quote {
	return &Quote{}
}

func (q *Quote) GetRandomAcademia() string {
	// Get a random anime quote
	resp, err := http.Get("https://animechan.vercel.app/api/random/anime?title=My%20Hero%20Academia")
	if err != nil {
		fmt.Println("error getting anime quote,", err)
		return ""
	}
	defer resp.Body.Close()

	// Decode the response
	var quote Quote
	err = json.NewDecoder(resp.Body).Decode(&quote)
	if err != nil {
		fmt.Println("error decoding anime quote,", err)
		return ""
	}

	return "Random quote from My Hero Academia: " + quote.Quote + " - " + quote.Character
}

func (q *Quote) GetRandom() string {
	// Get a random anime quote
	resp, err := http.Get("https://animechan.vercel.app/api/random")
	if err != nil {
		fmt.Println("error getting anime quote,", err)
		return ""
	}
	defer resp.Body.Close()

	// Decode the response
	var quote Quote
	err = json.NewDecoder(resp.Body).Decode(&quote)
	if err != nil {
		fmt.Println("error decoding anime quote,", err)
		return ""
	}

	return "Random quote from " + quote.Anime + ": " + quote.Quote + " - " + quote.Character
}
