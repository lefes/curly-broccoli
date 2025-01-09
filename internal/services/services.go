package services

import (
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

type UserActivities interface {
	AddOrUpdateUser(userID string) *domain.UserActivity
	MarkLimitReached(userID string)
	Reset() []*domain.UserActivity
	GetMaxMessages() int
}

type Services struct {
	User         Users
	Transaction  Transactions
	Discord      Discord
	Roles        Roles
	UserActivity UserActivities
}

func NewServices(repos *repository.Repositories) *Services {
	userService := NewUsersService(repos.User)
	transactionService := NewTransactionService(repos.Transaction)
	discordService := NewDiscordService(repos.Discord)
	rolesService := NewRoleService(repos.Role)
	userActivityService := NewUserActivitiesService(repos.UserActivity)
	return &Services{
		User:         userService,
		Transaction:  transactionService,
		Discord:      discordService,
		Roles:        rolesService,
		UserActivity: userActivityService,
	}
}
