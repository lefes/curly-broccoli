package services

import (
	"github.com/bwmarrin/discordgo"
	"github.com/lefes/curly-broccoli/config"
	"github.com/lefes/curly-broccoli/internal/domain"
	"github.com/lefes/curly-broccoli/internal/repository"
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
}

type Weather interface {
	GetWeather(city string, days int) (*discordgo.MessageEmbed, error)
}

type Transactions interface {
	CreateTransaction(transaction *domain.Transaction) error
	GetTransactionsByUserID(userID int) ([]*domain.Transaction, error)
	GetAllTransactions() ([]*domain.Transaction, error)
}

type Discord interface {
	//GetAllUsers(guildID string) (*domain.DiscordMembers, error)
	//BotRegister() (error, *discordgo.Session)
	Open() (*discordgo.Session, error)
}

type Roles interface {
	PromoteUser(userID string) error
	DemoteUser(userID string) error
	GetUserRole(userID string) (int, error)
}

type MessageHandler interface {
	HandleMessage(message *domain.Message, ctx *MessageHandlerContext)
	AddHandler(handler func(*domain.Message, *MessageHandlerContext) bool)
	HandlePoints(msg *domain.Message, ctx *MessageHandlerContext) bool
}

type Services struct {
	User        Users
	Transaction Transactions
	Discord     Discord
	Roles       Roles
	Weather     Weather
	MsgHandler  MessageHandler
}

func NewServices(repos *repository.Repositories, conf *config.Configs) *Services {
	userService := NewUsersService(repos.User)
	transactionService := NewTransactionService(repos.Transaction)
	discordService := NewDiscordService(&conf.Discord)
	rolesService := NewRoleService(repos.Role)
	weatherService := NewWeatherService(&conf.Weather)
	msgHandlerService := NewMessageHandlerService()
	return &Services{
		User:        userService,
		Transaction: transactionService,
		Discord:     discordService,
		Roles:       rolesService,
		Weather:     weatherService,
		MsgHandler:  msgHandlerService,
	}
}
