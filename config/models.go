package config

type Configs struct {
	Weather WeatherService
	Discord DiscordService
}

type WeatherService struct {
	ApiKey string
	ApiUrl string
}

type DiscordService struct {
	BotToken string
	GuildID  string
}
