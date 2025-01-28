package repository

import (
	"database/sql"

	"github.com/lefes/curly-broccoli/internal/domain"
	"github.com/lefes/curly-broccoli/internal/logging"
)

type Users interface {
	CreateUser(user *domain.User) error
	GetAllUsers() ([]*domain.User, error)
	DeleteUser(userId string) error
	UpdateUserPoints(discordID string, points int) error
	UpdateUserDailyMessages(discordID string, dailyMessages int) error
	AddOrUpdateUserActivity(userID string) *domain.UserActivity
	Reset() []*domain.UserActivity
	MarkLimitReached(userID string)
	GetMaxMessages() int
	IsLimitReached(userdID string) bool
}

type Repositories struct {
	User Users
}

func NewRepository(db *sql.DB, l *logging.Logger) *Repositories {
	activities := domain.NewUserActivities(25)
	return &Repositories{
		User: NewUsersRepo(db, activities, l),
	}
}
