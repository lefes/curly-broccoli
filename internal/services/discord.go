package services

import (
	"fmt"

	"github.com/lefes/curly-broccoli/config"
	"github.com/lefes/curly-broccoli/internal/domain"
	"github.com/lefes/curly-broccoli/internal/repository"
	"github.com/lefes/curly-broccoli/internal/transport/discordapi"
)

type DiscordService struct {
	config  *config.DiscordService
	session *discordapi.DiscordSession
	repo    *repository.Repositories
}

func NewDiscordService(conf *config.DiscordService, s *discordapi.DiscordSession, r *repository.Repositories) *DiscordService {
	return &DiscordService{
		config:  conf,
		session: s,
		repo:    r,
	}
}

func (ds *DiscordService) SyncUsers() error {
	discordUsers := make(map[string]bool)
	dbUsers := make(map[string]bool)

	members, err := ds.session.GetAllUsers(ds.config.GuildID)
	if err != nil {
		return fmt.Errorf("Failed to get all users: %w", err)
	}

	users, err := ds.repo.User.GetAllUsers()
	if err != nil {
		return fmt.Errorf("Error getting all users from database: %w", err)
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
				return fmt.Errorf("Error adding user to database: %w", err)
			} else {
				fmt.Println("Added user to database:", member.User.ID)
			}
		}
	}

	for userID := range dbUsers {
		if !discordUsers[userID] {
			err := ds.repo.User.DeleteUser(userID)
			if err != nil {
				return fmt.Errorf("Error deleting user from database: %w", err)

			} else {
				fmt.Println("Deleted user from database:", userID)
			}
		}
	}

	return nil
}
