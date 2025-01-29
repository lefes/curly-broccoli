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
	WillReachPointLimit(userID string, points int) bool
}

type Weather interface {
	GetWeather(city string, days int) (*discordgo.MessageEmbed, error)
}

type Discord interface {
	SyncUsers() error
	IsValidReaction(message *discordgo.Message, reactorID string) bool
}

type Roles interface {
	WillGetPromotion(userID string, respectToAdd int) (bool, *domain.Role)
	GetUserRole(userID string) (*domain.Role, error)
}

type Services struct {
	User    Users
	Weather Weather
	Discord Discord
	Roles   Roles
}

func NewServices(repos *repository.Repositories, conf *config.Configs, s *discordapi.DiscordSession, l *logging.Logger) *Services {
	userService := NewUsersService(repos.User, l)
	weatherService := NewWeatherService(&conf.Weather, l)
	discordService := NewDiscordService(&conf.Discord, s, repos, l)
	rolesService := NewRoleService(repos.Roles, l)
	return &Services{
		User:    userService,
		Weather: weatherService,
		Discord: discordService,
		Roles:   rolesService,
	}
}
