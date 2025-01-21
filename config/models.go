package config

type Configs struct {
	Weather WeatherService
	Discord DiscordService
	Storage StorageConfig
}

type StorageConfig struct {
	DbPath string
}

type WeatherService struct {
	ApiKey string
	ApiUrl string
}

type DiscordService struct {
	BotToken string
	GuildID  string
}
