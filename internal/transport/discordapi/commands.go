package discordapi

import (
	"fmt"

	"github.com/bwmarrin/discordgo"
	"github.com/lefes/curly-broccoli/internal/domain"
)

func (ds *DiscordSession) RegisterCommands(commands []domain.SlashCommand, guildID string) ([]*discordgo.ApplicationCommand, error) {
	ds.s.AddHandler(func(s *discordgo.Session, i *discordgo.InteractionCreate) {
		for _, cmd := range commands {
			if cmd.Name == i.ApplicationCommandData().Name {
				cmd.Handler(ds.s, i)
				break
			}
		}
	})
	registeredCommands := make([]*discordgo.ApplicationCommand, len(commands))
	for i, cmd := range commands {
		command, err := ds.s.ApplicationCommandCreate(ds.s.State.User.ID, guildID, &discordgo.ApplicationCommand{
			Name:        cmd.Name,
			Description: cmd.Description,
			Options:     cmd.Options,
		})
		if err != nil {
			return nil, fmt.Errorf("failed to register command %s: %w", cmd.Name, err)
		}
		registeredCommands[i] = command
	}

	return registeredCommands, nil
}

func (ds *DiscordSession) DeleteCommands(commands []*discordgo.ApplicationCommand, guildID string) error {
	for _, cmd := range commands {
		err := ds.s.ApplicationCommandDelete(ds.s.State.User.ID, guildID, cmd.ID)
		if err != nil {
			return fmt.Errorf("failed to delete command '%s': %w", cmd.Name, err)
		}
	}
	return nil
}
