package config

import (
	"os"

	"github.com/lefes/curly-broccoli/internal/logging"
	"github.com/sirupsen/logrus"
)

var (
	weatherApiKey string
	weatherApiUrl = "https://weather.visualcrossing.com/VisualCrossingWebServices/rest/services/timeline"
	botToken      string
	cLogger       *logrus.Entry
	GuildID       = "147313959819542528"
	dbPath        string
)

func Init() *Configs {
	cLogger = logging.GetLogger("config")
	conf := &Configs{}
	parseConfig(conf)
	return conf
}

func parseConfig(conf *Configs) {
	conf.Weather.ApiKey = getEnv("WEATHER_API_KEY")
	conf.Weather.ApiUrl = getEnvOptional("WEATHER_API_URL", weatherApiUrl)
	conf.Discord.BotToken = getEnv("DISCORD_BOT_TOKEN")
	conf.Discord.GuildID = getEnvOptional("DISCORD_GUILD_ID", GuildID)
	conf.Storage.DbPath = getEnvOptional("DB_PATH", "data/bot.db")
}

func getEnv(key string) string {
	value, ok := os.LookupEnv(key)
	if !ok {
		cLogger.Fatalf("Missing %s environment variable", key)
	}
	return value
}

func getEnvOptional(key string, defaultVal string) string {
	value, ok := os.LookupEnv(key)
	if !ok {
		return defaultVal
	}
	return value
}
