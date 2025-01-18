package services

import (
	"github.com/lefes/curly-broccoli/internal/domain"
	"github.com/lefes/curly-broccoli/internal/repository"
)

type UserService struct {
	repo repository.Users
}

func NewUsersService(repo repository.Users) *UserService {
	return &UserService{repo: repo}
}

func (s *UserService) GetAllUsers() ([]*domain.User, error) {
	return s.repo.GetAllUsers()
}

func (s *UserService) CreateUser(user *domain.User) error {
	return s.repo.CreateUser(user)
}

func (s *UserService) DeleteUser(userID string) error {
	return s.repo.DeleteUser(userID)
}

func (s *UserService) GetUserByDiscordID(discordID string) (*domain.User, error) {
	return s.repo.GetUserByDiscordID(discordID)
}

func (s *UserService) UpdateUserRole(discordID string, roleID int) error {
	return s.repo.UpdateUserRole(discordID, roleID)
}

func (s *UserService) UpdateUserPoints(discordID string, points int) error {
	return s.repo.UpdateUserPoints(discordID, points)
}

func (s *UserService) UpdateUserRespect(discordID string, respect int) error {
	return s.repo.UpdateUserRespect(discordID, respect)
}

func (s *UserService) UpdateUserDailyMessages(discordID string, dailyMessages int) error {
	return s.repo.UpdateUserDailyMessages(discordID, dailyMessages)
}

func (s *UserService) AddOrUpdateUserActivity(userID string) *domain.UserActivity {
	return s.repo.AddOrUpdateUserActivity(userID)
}

func (s *UserService) MarkLimitReached(userID string) {
	s.repo.MarkLimitReached(userID)
}

func (s *UserService) Reset() []*domain.UserActivity {
	return s.repo.Reset()
}

func (s *UserService) GetMaxMessages() int {
	return s.repo.GetMaxMessages()
}

func (s *UserService) IsLimitReached(userID string) bool {
	return s.repo.IsLimitReached(userID)
}
