package services

import (
	"fmt"

	"github.com/lefes/curly-broccoli/internal/transport/http/weatherapi"
)

type WeatherService struct {
	Client *weatherapi.Client
}

func NewWeatherService(client *weatherapi.Client) *WeatherService {
	return &WeatherService{Client: client}
}

func (s *WeatherService) CurrentWeather(city string) (*weatherapi.WeatherResponse, error) {
	response, err := s.Client.CurrentWeather(city)
	if err != nil {
		return nil, fmt.Errorf("failed to fetch current weather: %w", err)
	}
	return response, nil
}

func (s *WeatherService) ForecastWeather(city string, days int) (*weatherapi.WeatherResponse, error) {
	if days < 1 || days > 7 {
		return nil, fmt.Errorf("days must be between 1 and 7")
	}
	response, err := s.Client.ForecastWeather(city, days)
	if err != nil {
		return nil, fmt.Errorf("failed to fetch forecast weather: %w", err)
	}
	return response, nil
}
