package services

import (
	"github.com/lefes/curly-broccoli/internal/domain"
	"github.com/lefes/curly-broccoli/internal/repository"
	"github.com/lefes/curly-broccoli/internal/transport/http/weatherapi"
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

type WeatherAPI interface {
	CurrentWeather(city string) (*weatherapi.WeatherResponse, error)
	ForecastWeather(city string, days int) (*weatherapi.WeatherResponse, error)
}

type Transactions interface {
	CreateTransaction(transaction *domain.Transaction) error
	GetTransactionsByUserID(userID int) ([]*domain.Transaction, error)
	GetAllTransactions() ([]*domain.Transaction, error)
}

type Discord interface {
	GetAllUsers(guildID string) (*domain.DiscordMembers, error)
}

type Roles interface {
	PromoteUser(userID string) error
	DemoteUser(userID string) error
	GetUserRole(userID string) (int, error)
}

type Services struct {
	User        Users
	Transaction Transactions
	Discord     Discord
	Roles       Roles
	WAPI        WeatherAPI
}

func NewServices(repos *repository.Repositories, wClient *weatherapi.Client) *Services {
	userService := NewUsersService(repos.User)
	transactionService := NewTransactionService(repos.Transaction)
	discordService := NewDiscordService(repos.Discord)
	rolesService := NewRoleService(repos.Role)
	weatherService := NewWeatherService(wClient)
	return &Services{
		User:        userService,
		Transaction: transactionService,
		Discord:     discordService,
		Roles:       rolesService,
		WAPI:        weatherService,
	}
}
