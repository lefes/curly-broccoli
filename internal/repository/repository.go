package repository

import (
	"database/sql"

	"github.com/bwmarrin/discordgo"
	"github.com/lefes/curly-broccoli/internal/domain"
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

type Role interface {
	PromoteUser(userID string) error
	DemoteUser(userID string) error
	GetUserRole(userID string) (int, error)
}

type UserActivity interface {
	AddOrUpdateUser(userID string) *domain.UserActivity
	MarkLimitReached(userID string)
	Reset() []*domain.UserActivity
	GetMaxMessages() int
	IsLimitReached(userdID string) bool
}

type Repositories struct {
	User         Users
	Transaction  Transactions
	Discord      Discord
	Role         Role
	UserActivity UserActivity
}

func NewRepository(db *sql.DB, discordSession *discordgo.Session) *Repositories {
	activities := domain.NewUserActivities(25)
	return &Repositories{
		User:         NewUsersRepo(db),
		Transaction:  NewTransactionRepo(db),
		Discord:      NewDiscordRepo(discordSession),
		Role:         NewRoleRepo(db),
		UserActivity: NewUserActivityRepo(activities),
	}
}
