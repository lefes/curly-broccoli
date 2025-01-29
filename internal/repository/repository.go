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
	AddUserPoints(discordID string, points int) error
	RemoveUserPoints(discordID string, points int) error
	UpdateUserDailyMessages(discordID string, dailyMessages int) error
	AddOrUpdateUserActivity(userID string) *domain.UserActivity
	Reset() []*domain.UserActivity
	MarkLimitReached(userID string)
	GetMaxMessages() int
	IsLimitReached(userdID string) bool
	AddDayPoints(discordID string, points int) error
	RemoveDayPoints(discordID string, points int) error
	GetTodayPoints(discordID string) (int, error)
}

type Roles interface {
	GetUserRespect(discordID string) (int, error)
	AddUserRespect(discordID string, respect int) error
	RemoveUserRespect(discordID string, respect int) error
	UpdateUserRole(discordID string, roleID int) error
	GetUserRole(discordID string) (int, error)
}

type Repositories struct {
	User  Users
	Roles Roles
}

func NewRepository(db *sql.DB, l *logging.Logger) *Repositories {
	activities := domain.NewUserActivities(25)
	return &Repositories{
		User:  NewUsersRepo(db, activities, l),
		Roles: NewRolesRepo(db, l),
	}
}
