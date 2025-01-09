package services

import (
	"github.com/lefes/curly-broccoli/internal/domain"
	"github.com/lefes/curly-broccoli/internal/repository"
)

type DiscordService struct {
	repo repository.Discord
}

func NewDiscordService(repo repository.Discord) *DiscordService {
	return &DiscordService{repo: repo}
}

func (d *DiscordService) GetAllUsers(guildID string) (*domain.DiscordMembers, error) {
	return d.repo.GetAllUsers(guildID)
}
