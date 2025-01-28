package services

import (
	"github.com/bwmarrin/discordgo"
	"github.com/lefes/curly-broccoli/config"
	"github.com/lefes/curly-broccoli/internal/domain"
	"github.com/lefes/curly-broccoli/internal/logging"
	"github.com/lefes/curly-broccoli/internal/repository"
	"github.com/lefes/curly-broccoli/internal/transport/discordapi"
)

type Users interface {
	Reset() []*domain.UserActivity
	CanSendMessage(msg *domain.Message) (*domain.UserActivity, bool)
	IncrementUserMessageCount(activity *domain.UserActivity)
	ReactionPoints(message *discordgo.Message) bool
	ReactionPointsRemoval(message *discordgo.Message) bool
}

type Weather interface {
	GetWeather(city string, days int) (*discordgo.MessageEmbed, error)
}

type Discord interface {
	SyncUsers() error
	IsValidReaction(message *discordgo.Message, reactorID string) bool
}

type Services struct {
	User    Users
	Weather Weather
	Discord Discord
}

func NewServices(repos *repository.Repositories, conf *config.Configs, s *discordapi.DiscordSession, l *logging.Logger) *Services {
	userService := NewUsersService(repos.User, l)
	weatherService := NewWeatherService(&conf.Weather, l)
	discordService := NewDiscordService(&conf.Discord, s, repos, l)
	return &Services{
		User:    userService,
		Weather: weatherService,
		Discord: discordService,
	}
}
