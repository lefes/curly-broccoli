package services

import (
	"time"

	"github.com/bwmarrin/discordgo"
	"github.com/lefes/curly-broccoli/config"
	"github.com/lefes/curly-broccoli/internal/domain"
	"github.com/lefes/curly-broccoli/internal/logging"
	"github.com/lefes/curly-broccoli/internal/repository"
)

type UserService struct {
	repo   repository.Users
	logger *logging.Logger
}

func NewUsersService(repo repository.Users, l *logging.Logger) *UserService {
	return &UserService{repo: repo, logger: l}
}

func (s *UserService) IsAdmin(userID string) bool {
	_, isAdmin := config.AdminUsers[userID]
	return isAdmin
}

func (s *UserService) Reset() []*domain.UserActivity {
	return s.repo.Reset()
}

func (s *UserService) WillReachPointLimit(userID string, points int) bool {
	userPoints, err := s.repo.GetTodayPoints(userID)
	if err != nil {
		s.logger.Errorf("Error getting user points from database: %s", err)
	}

	afterPoints := userPoints + points

	if afterPoints > config.DayPointsLimit {
		return false
	}
	return true
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
		s.logger.Infof("User %s reached daily limit", msg.Username)
		s.repo.MarkLimitReached(msg.Author)

		if s.WillReachPointLimit(msg.Author, 25) {
			return nil, false
		}

		err := s.repo.AddUserPoints(msg.Author, 25)
		if err != nil {
			s.logger.Errorf("Error updating user points in database: %s", err)
			return nil, false
		}

		err = s.repo.AddDayPoints(msg.Author, 25)
		if err != nil {
			s.logger.Errorf("Error updating user daily points in database: %s", err)
			return nil, false
		}
		s.logger.Infof("User %s received 25 points", msg.Username)
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

func (s *UserService) ReactionPoints(message *discordgo.Message) bool {
	messageAuthor := message.Author.ID

	if s.WillReachPointLimit(messageAuthor, 1) {
		return false
	}

	err := s.repo.AddDayPoints(messageAuthor, 1)
	if err != nil {
		s.logger.Errorf("Error updating user daily points in database: %s", err)
		return false
	}

	err = s.repo.AddUserPoints(messageAuthor, 1)
	if err != nil {
		s.logger.Errorf("Error updating user points in database: %s", err)
		return false
	}

	return true
}

func (s *UserService) ReactionPointsRemoval(message *discordgo.Message) bool {

	messageAuthor := message.Author.ID

	err := s.repo.RemoveDayPoints(messageAuthor, 1)
	if err != nil {
		s.logger.Errorf("Error updating user daily points in database: %s", err)
		return false
	}

	err = s.repo.RemoveUserPoints(messageAuthor, 1)
	if err != nil {
		s.logger.Errorf("Error updating user points in database: %s", err)
		return false
	}

	return true
}
