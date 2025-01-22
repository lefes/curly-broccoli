package services

import (
	"github.com/bwmarrin/discordgo"
	"github.com/lefes/curly-broccoli/config"
	"github.com/lefes/curly-broccoli/internal/domain"
	"github.com/lefes/curly-broccoli/internal/repository"
	"github.com/lefes/curly-broccoli/internal/transport/discordapi"
)

type Users interface {
	CreateUser(user *domain.User) error
	GetUserByDiscordID(discordID string) (*domain.User, error)
	GetAllUsers() ([]*domain.User, error)
	DeleteUser(userId string) error
	UpdateUserRole(discordID string, roleID int) error
	UpdateUserPoints(discordID string, points int) error
	UpdateUserRespect(discordID string, respect int) error
	UpdateUserDailyMessages(discordID string, dailyMessages int) error
	AddOrUpdateUserActivity(userID string) *domain.UserActivity
	Reset() []*domain.UserActivity
	MarkLimitReached(userID string)
	GetMaxMessages() int
	IsLimitReached(userdID string) bool
	CanSendMessage(msg *domain.Message) (*domain.UserActivity, bool)
	IncrementUserMessageCount(activity *domain.UserActivity)
}

type Weather interface {
	GetWeather(city string, days int) (*discordgo.MessageEmbed, error)
}

type Transactions interface {
	CreateTransaction(transaction *domain.Transaction) error
	GetTransactionsByUserID(userID int) ([]*domain.Transaction, error)
	GetAllTransactions() ([]*domain.Transaction, error)
}

type Roles interface {
	PromoteUser(userID string) error
	DemoteUser(userID string) error
	GetUserRole(userID string) (int, error)
}

type Discord interface {
	SyncUsers() error
}

type MessageHandler interface {
	HandleMessage(message *domain.Message, ctx *MessageHandlerContext)
	AddHandler(handler func(*domain.Message, *MessageHandlerContext) bool)
	HandlePoints(msg *domain.Message, ctx *MessageHandlerContext) bool
}

type Services struct {
	User        Users
	Transaction Transactions
	Roles       Roles
	Weather     Weather
	MsgHandler  MessageHandler
	Discord     Discord
}

func NewServices(repos *repository.Repositories, conf *config.Configs, s *discordapi.DiscordSession) *Services {
	userService := NewUsersService(repos.User)
	transactionService := NewTransactionService(repos.Transaction)
	rolesService := NewRoleService(repos.Role)
	weatherService := NewWeatherService(&conf.Weather)
	msgHandlerService := NewMessageHandlerService()
	discordService := NewDiscordService(&conf.Discord, s, repos)
	return &Services{
		User:        userService,
		Transaction: transactionService,
		Roles:       rolesService,
		Weather:     weatherService,
		MsgHandler:  msgHandlerService,
		Discord:     discordService,
	}
}
