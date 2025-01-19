package repository

import (
	"database/sql"

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
	AddOrUpdateUserActivity(userID string) *domain.UserActivity
	Reset() []*domain.UserActivity
	MarkLimitReached(userID string)
	GetMaxMessages() int
	IsLimitReached(userdID string) bool
}

type Transactions interface {
	CreateTransaction(transaction *domain.Transaction) error
	GetTransactionsByUserID(userID int) ([]*domain.Transaction, error)
	GetAllTransactions() ([]*domain.Transaction, error)
}

type Role interface {
	PromoteUser(userID string) error
	DemoteUser(userID string) error
	GetUserRole(userID string) (int, error)
}

type Repositories struct {
	User        Users
	Transaction Transactions
	Role        Role
}

func NewRepository(db *sql.DB) *Repositories {
	activities := domain.NewUserActivities(25)
	return &Repositories{
		User:        NewUsersRepo(db, activities),
		Transaction: NewTransactionRepo(db),
		Role:        NewRoleRepo(db),
	}
}
