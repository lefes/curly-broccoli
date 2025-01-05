package weather

import (
	"fmt"

	"github.com/lefes/curly-broccoli/pkg/logging"
	"github.com/sirupsen/logrus"
)

var logger *logrus.Entry

func InitWeatherLogger() {
	logger = logging.GetLogger("weather")

}

func GetCurrentWeather(apiKey, baseURL, city string) (error, *WeatherResponse) {
	client := &Client{APIKey: apiKey, BaseURL: baseURL}

	weather, err := client.CurrentWeather(city)
	if err != nil {
		logger.Error("Error fetching current weather:", err)
		return err, nil
	}

	logger.Infof("Current weather in %s: %.1f°C, Condition: %s\n",
		weather.ResolvedAddress, weather.Days[0].TempMax, weather.Days[0].Condition)
	return nil, weather
}

func GetForecastWeather(apiKey, baseURL, city string, days int) (error, *WeatherResponse) {
	client := &Client{APIKey: apiKey, BaseURL: baseURL}

	forecast, err := client.ForecastWeather(city, days)
	if err != nil {
		logger.Error("Error fetching forecast:", err)
		return err, nil
	}

	fmt.Printf("Forecast for: %s\n", forecast.ResolvedAddress)
	for _, day := range forecast.Days {
		logger.Infof("Date: %s, TempMax: %.1f°C, TempMin: %.1f°C, Condition: %s\n",
			day.Datetime, day.TempMax, day.TempMin, day.Condition)
	}
	return nil, forecast
}
