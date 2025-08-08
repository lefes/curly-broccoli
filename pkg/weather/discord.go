package weather

import (
	"fmt"
	"strconv"
	"strings"
	"time"

	"github.com/bwmarrin/discordgo"
)

var cityShortcuts = map[string]string{
	"мск": "Moscow",
	"спб": "Saint Petersburg",
	"екб": "Yekaterinburg",
	"нск": "Novosibirsk",
	"крд": "Krasnodar",
	"соч": "Sochi",
}

func sendHelpMessage(session *discordgo.Session, message *discordgo.MessageCreate) error {
	shortcutsDescription := ""
	for shortcut, city := range cityShortcuts {
		shortcutsDescription += fmt.Sprintf("`%s`: %s\n", shortcut, city)
	}

	embed := &discordgo.MessageEmbed{
		Title:       "Weather Command Help",
		Description: "Команда `!weather` используется для получения информации о погоде.\n\n**Важно:** Названия городов должны быть на английском языке.",
		Color:       0x3498db, // Blue
		Fields: []*discordgo.MessageEmbedField{
			{
				Name:  "🔹 `!weather <city>`",
				Value: "Получить текущую погоду в указанном городе.",
			},
			{
				Name:  "🔹 `!weather <city> <days>`",
				Value: "Получить прогноз погоды на указанное количество дней (максиум 7).",
			},
			{
				Name:  "🔹 Шорткаты городов",
				Value: shortcutsDescription,
			},
		},
		Footer: &discordgo.MessageEmbedFooter{
			Text: "Пример: !weather Moscow\nПример: !weather Moscow 3",
		},
	}

	_, err := session.ChannelMessageSendEmbed(message.ChannelID, embed)
	return err
}

func sendForecastWeatherMessage(session *discordgo.Session, message *discordgo.MessageCreate, weather *WeatherResponse) error {
	fields := []*discordgo.MessageEmbedField{}
	for _, day := range weather.Days {
		fields = append(fields, &discordgo.MessageEmbedField{
			Name: day.Datetime,
			Value: fmt.Sprintf(
				"🌡️ **Max:** %.1f°C\n"+
					"🌡️ **Min:** %.1f°C\n"+
					"🌤️ **Condition:** %s\n"+
					"🌅 **Sunrise:** %s\n"+
					"🌇 **Sunset:** %s",
				day.TempMax, day.TempMin, day.Condition, day.SunRise, day.SunSet),
			Inline: false,
		})
	}

	embed := &discordgo.MessageEmbed{
		Title:       fmt.Sprintf("Weather Forecast for %s", weather.ResolvedAddress),
		Description: "Here's the weather forecast for the next days:",
		Color:       0x1abc9c,
		Fields:      fields,
		Timestamp:   time.Now().Format(time.RFC3339),
	}

	_, err := session.ChannelMessageSendEmbed(message.ChannelID, embed)
	return err
}

func sendCurrentWeatherMessage(session *discordgo.Session, message *discordgo.MessageCreate, weather *WeatherResponse) error {
	embed := &discordgo.MessageEmbed{
		Title:       "Current Weather",
		Description: fmt.Sprintf("Here's the current weather for **%s**", weather.ResolvedAddress),
		Color:       0x3498db,
		Fields: []*discordgo.MessageEmbedField{
			{
				Name:   "📍 Location",
				Value:  weather.ResolvedAddress,
				Inline: true,
			},
			{
				Name:   "🌡️ Temperature",
				Value:  fmt.Sprintf("**%.1f°C** (Max: %.1f°C, Min: %.1f°C)", weather.CurrentConditions.Temp, weather.Days[0].TempMax, weather.Days[0].TempMin),
				Inline: true,
			},
			{
				Name:   "🌤️ Condition",
				Value:  weather.CurrentConditions.Condition,
				Inline: true,
			},
			{
				Name:   "📅 Date",
				Value:  weather.Days[0].Datetime,
				Inline: true,
			},
			{
				Name:   "🌅 Sunrise",
				Value:  weather.Days[0].SunRise,
				Inline: true,
			},
			{
				Name:   "🌇 Sunset",
				Value:  weather.Days[0].SunSet,
				Inline: true,
			},
		},
		Timestamp: time.Now().Format(time.RFC3339),
	}

	_, err := session.ChannelMessageSendEmbed(message.ChannelID, embed)
	return err
}

func resolveShortCut(input string) string {
	normalizedInput := strings.ToLower(input)

	if fullName, exists := cityShortcuts[normalizedInput]; exists {
		return fullName
	}

	return input
}

func HandleWeatherMessage(client Client, session *discordgo.Session, message *discordgo.MessageCreate, cmdMatches []string) error {
	if cmdMatches[2] == "" {
		return sendHelpMessage(session, message)
	}

	city := resolveShortCut(cmdMatches[2])
	days := 0
	if cmdMatches[3] != "" {
		parsedDays, err := strconv.Atoi(cmdMatches[3])
		if err != nil {
			session.ChannelMessageSend(message.ChannelID, "Ошибка: Количество дней должно быть числом.")
			return err
		}
		days = parsedDays
	}
	if days == 0 {
		weather, err := client.CurrentWeather(city)
		if err != nil {
			session.ChannelMessageSend(message.ChannelID, "Ошибка: Не удалось получить текущую погоду.")
			return err
		}
		err = sendCurrentWeatherMessage(session, message, weather)
		if err != nil {
			session.ChannelMessageSend(message.ChannelID, "Ошибка: Не удалось отправить сообщение с текущей погодой.")
			return err
		}
		return nil
	}

	weather, err := client.ForecastWeather(city, days)
	if err != nil {
		session.ChannelMessageSend(message.ChannelID, "Ошибка: Не удалось получить прогноз погоды.")
		return err
	}
	err = sendForecastWeatherMessage(session, message, weather)
	if err != nil {
		session.ChannelMessageSend(message.ChannelID, "Ошибка: Не удалось отправить сообщение с прогнозом погоды.")
		return err
	}

	return nil
}
