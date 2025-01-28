package services

import (
	"github.com/bwmarrin/discordgo"
	"github.com/lefes/curly-broccoli/config"
	"github.com/lefes/curly-broccoli/internal/domain"
	"github.com/lefes/curly-broccoli/internal/logging"
	"github.com/lefes/curly-broccoli/internal/repository"
	"github.com/lefes/curly-broccoli/internal/transport/discordapi"
)

type DiscordService struct {
	config  *config.DiscordService
	session *discordapi.DiscordSession
	repo    *repository.Repositories
	logger  *logging.Logger
}

func NewDiscordService(conf *config.DiscordService, s *discordapi.DiscordSession, r *repository.Repositories, l *logging.Logger) *DiscordService {
	return &DiscordService{
		config:  conf,
		session: s,
		repo:    r,
		logger:  l,
	}
}

func (ds *DiscordService) SyncUsers() error {
	discordUsers := make(map[string]bool)
	dbUsers := make(map[string]bool)

	members, err := ds.session.GetAllUsers(ds.config.GuildID)
	if err != nil {
		return ds.logger.Errorf("Failed to get all users: %w", err)
	}

	users, err := ds.repo.User.GetAllUsers()
	if err != nil {
		return ds.logger.Errorf("Error getting all users from database: %w", err)
	}

	for _, member := range members.Members {
		discordUsers[member.User.ID] = true
	}

	for _, user := range users {
		dbUsers[user.DiscordID] = true
	}

	for _, member := range members.Members {
		if !dbUsers[member.User.ID] {
			err := ds.repo.User.CreateUser(&domain.User{
				DiscordID:     member.User.ID,
				Username:      member.User.Username,
				RoleID:        1,
				Points:        0,
				Respect:       0,
				DailyMessages: 0,
			})
			if err != nil {
				return ds.logger.Errorf("Error adding user to database: %w", err)
			} else {
				ds.logger.Info("Added user to database:", member.User.Username)
			}
		}
	}

	for userID := range dbUsers {
		if !discordUsers[userID] {
			err := ds.repo.User.DeleteUser(userID)
			if err != nil {
				return ds.logger.Errorf("Error deleting user from database: %w", err)

			} else {
				ds.logger.Info("Deleted user from database:", userID)
			}
		}
	}

	return nil
}

func (ds *DiscordService) IsValidReaction(message *discordgo.Message, reactorID string) bool {
	/*  s := ds.session.GetSession() */
	/* if reactorID == message.Author.ID || reactorID == s.State.User.ID { */
	/* return false */
	/* } */

	messageReactions := message.Reactions
	totalReactions := 0

	for _, reaction := range messageReactions {
		totalReactions += reaction.Count
	}

	if totalReactions > 5 {
		return false
	}

	return true
}
