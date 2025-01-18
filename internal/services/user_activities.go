package services

import (
	"github.com/lefes/curly-broccoli/internal/domain"
	"github.com/lefes/curly-broccoli/internal/repository"
)

type UserActivitiesService struct {
	repo repository.UserActivity
}

func NewUserActivitiesService(repo repository.UserActivity) *UserActivitiesService {
	return &UserActivitiesService{repo: repo}
}

func (s *UserActivitiesService) AddOrUpdateUser(userID string) *domain.UserActivity {
	return s.repo.AddOrUpdateUser(userID)
}

func (s *UserActivitiesService) MarkLimitReached(userID string) {
	s.repo.MarkLimitReached(userID)
}

func (s *UserActivitiesService) Reset() []*domain.UserActivity {
	return s.repo.Reset()
}

func (s *UserActivitiesService) GetMaxMessages() int {
	return s.repo.GetMaxMessages()
}

func (s *UserActivitiesService) IsLimitReached(userID string) bool {
	return s.repo.IsLimitReached(userID)
}
