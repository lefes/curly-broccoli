package repository

import (
	"log"

	"github.com/bwmarrin/discordgo"
	"github.com/lefes/curly-broccoli/internal/domain"
)

type DiscordRepo struct {
	s *discordgo.Session
}

func NewDiscordRepo(s *discordgo.Session) *DiscordRepo {
	return &DiscordRepo{s: s}
}

func (r *DiscordRepo) GetAllUsers(guildID string) (*domain.DiscordMembers, error) {
	allMembers := &domain.DiscordMembers{}
	lastUserID := ""

	for {
		members, err := r.s.GuildMembers(guildID, lastUserID, 1000)
		if err != nil {
			log.Printf("Error fetching members: %v", err)
			return nil, err
		}

		if len(members) == 0 {
			break
		}

		allMembers.Members = append(allMembers.Members, members...)

		lastUserID = members[len(members)-1].User.ID
	}

	return allMembers, nil
}

func (r *DiscordRepo) GetDiscordIdByUsername(username string) (string, error) {
	user, err := r.s.User(username)
	if err != nil {
		return "", err
	}

	return user.ID, nil
}
