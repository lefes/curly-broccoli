package weather

import (
	"encoding/json"
	"fmt"
	"net/http"
	"net/url"
	"strconv"
)

type Client struct {
	APIKey  string
	BaseURL string
}

func NewClient(apiKey, baseUrl string) *Client {
	return &Client{
		APIKey:  apiKey,
		BaseURL: baseUrl,
	}
}

func (c *Client) CurrentWeather(city string) (*WeatherResponse, error) {
	cityEncoded := url.QueryEscape(city)
	url := fmt.Sprintf("%s/%s/today?key=%s&unitGroup=metric&include=current&contentType=json&lang=ru", c.BaseURL, cityEncoded, c.APIKey)

	resp, err := http.Get(url)
	if err != nil {
		return nil, fmt.Errorf("failed to fetch current weather: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("bad response: %s", resp.Status)
	}

	var result WeatherResponse
	if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
		return nil, fmt.Errorf("failed to decode response: %w", err)
	}

	return &result, nil
}

func (c *Client) ForecastWeather(city string, days int) (*WeatherResponse, error) {
	cityEncoded := url.QueryEscape(city)
	if days < 1 || days > 7 {
		return nil, fmt.Errorf("days must be between 1 and 7")
	}
	daysStr := strconv.Itoa(days)
	url := fmt.Sprintf("%s/%s/next%sdays?unitGroup=metric&include=days&key=%s&contentType=json&lang=ru",
		c.BaseURL, cityEncoded, daysStr, c.APIKey)

	resp, err := http.Get(url)
	if err != nil {
		return nil, fmt.Errorf("failed to fetch forecast weather: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("bad response: %s", resp.Status)
	}

	var result WeatherResponse
	if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
		return nil, fmt.Errorf("failed to decode response: %w", err)
	}

	return &result, nil
}
