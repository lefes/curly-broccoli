package discordapi

import (
	"fmt"

	"github.com/bwmarrin/discordgo"
	"github.com/lefes/curly-broccoli/internal/domain"
)

func RegisterCommands(s *discordgo.Session, commands []domain.SlashCommand, guildID string) ([]*discordgo.ApplicationCommand, error) {
	s.AddHandler(func(s *discordgo.Session, i *discordgo.InteractionCreate) {
		for _, cmd := range commands {
			if cmd.Name == i.ApplicationCommandData().Name {
				cmd.Handler(s, i)
				break
			}
		}
	})

	registeredCommands := make([]*discordgo.ApplicationCommand, len(commands))
	for i, cmd := range commands {
		command, err := s.ApplicationCommandCreate(s.State.User.ID, guildID, &discordgo.ApplicationCommand{
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

func DeleteCommands(s *discordgo.Session, commands []*discordgo.ApplicationCommand, guildID string) error {
	for _, cmd := range commands {
		err := s.ApplicationCommandDelete(s.State.User.ID, guildID, cmd.ID)
		if err != nil {
			return fmt.Errorf("failed to delete command '%s': %w", cmd.Name, err)
		}
		fmt.Printf("Command '%s' deleted successfully\n", cmd.Name)
	}
	return nil
}
