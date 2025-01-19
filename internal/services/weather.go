package services

import (
	"fmt"
	"time"

	"github.com/bwmarrin/discordgo"
	"github.com/lefes/curly-broccoli/config"
	"github.com/lefes/curly-broccoli/internal/transport/http/weatherapi"
)

type WeatherService struct {
	config *config.WeatherService
	Client *weatherapi.Client
}

func NewWeatherService(conf *config.WeatherService) *WeatherService {
	c := weatherapi.NewClient(conf.ApiKey, conf.ApiUrl)
	return &WeatherService{config: conf, Client: c}
}

func (w *WeatherService) GetWeather(city string, days int) (*discordgo.MessageEmbed, error) {
	embed := &discordgo.MessageEmbed{}
	if city == "" {
		return nil, fmt.Errorf("city cannot be empty")
	}
	if days < 1 || days > 7 {
		return nil, fmt.Errorf("days must be between 1 and 7")
	}
	if days == 1 {
		response, err := w.Client.CurrentWeather(city)
		if err != nil {
			return nil, fmt.Errorf("failed to fetch current weather: %w", err)
		}
		embed = &discordgo.MessageEmbed{
			Title:       "🌤️ Current Weather",
			Description: fmt.Sprintf("Here's the current weather for **%s**", response.ResolvedAddress),
			Color:       0x3498db,
			Fields: []*discordgo.MessageEmbedField{
				{
					Name:   "📍 Location",
					Value:  response.ResolvedAddress,
					Inline: true,
				},
				{
					Name:   "🌡️ Temperature",
					Value:  fmt.Sprintf("**%.1f°C** (Max: %.1f°C, Min: %.1f°C)", response.CurrentConditions.Temp, response.Days[0].TempMax, response.Days[0].TempMin),
					Inline: true,
				},
				{
					Name:   "🌤️ Condition",
					Value:  response.CurrentConditions.Condition,
					Inline: true,
				},
				{
					Name:   "📅 Date",
					Value:  response.Days[0].Datetime,
					Inline: true,
				},
			},
			Timestamp: time.Now().Format(time.RFC3339),
		}
		return embed, nil
	} else {
		response, err := w.Client.ForecastWeather(city, days)
		if err != nil {
			return nil, fmt.Errorf("failed to fetch forecast weather: %w", err)
		}
		for _, day := range response.Days {
			embed.Fields = append(embed.Fields, &discordgo.MessageEmbedField{
				Name:   fmt.Sprintf("📅 %s", day.Datetime),
				Value:  fmt.Sprintf("Max: %.1f°C, Min: %.1f°C\n🌤️ %s", day.TempMax, day.TempMin, day.Condition),
				Inline: false,
			})
		}
		return embed, nil
	}
}
