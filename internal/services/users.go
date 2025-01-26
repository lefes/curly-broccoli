package services

import (
	"fmt"
	"time"

	"github.com/lefes/curly-broccoli/internal/domain"
	"github.com/lefes/curly-broccoli/internal/repository"
)

type UserService struct {
	repo repository.Users
}

func NewUsersService(repo repository.Users) *UserService {
	return &UserService{repo: repo}
}

func (s *UserService) Reset() []*domain.UserActivity {
	return s.repo.Reset()
}

func (s *UserService) CanSendMessage(msg *domain.Message) (*domain.UserActivity, bool) {
	if len(msg.Content) < 5 {
		return nil, false
	}

	if s.repo.IsLimitReached(msg.Author) {
		return nil, false
	}

	activity := s.repo.AddOrUpdateUserActivity(msg.Author)

	now := time.Now()
	if now.Before(activity.NextMessageTime) {
		return activity, false
	}

	if activity.MessageCount == s.repo.GetMaxMessages() {
		fmt.Printf("User %s reached daily limit\n", msg.Username)
		s.repo.MarkLimitReached(msg.Author)
		err := s.repo.UpdateUserPoints(msg.Author, 25)
		if err != nil {
			fmt.Printf("Error updating user points in database: %s\n", err)
			return nil, false
		}
		fmt.Printf("User %s received 25 points\n", msg.Username)
		return activity, true
	}

	return activity, true
}

func (s *UserService) IncrementUserMessageCount(activity *domain.UserActivity) {
	now := time.Now()
	activity.LastMessageTime = now
	activity.NextMessageTime = now.Add(2 * time.Second)
	activity.MessageCount = activity.MessageCount + 1
}
