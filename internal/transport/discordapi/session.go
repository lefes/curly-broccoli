package discordapi

import (
	"fmt"

	"github.com/bwmarrin/discordgo"
	"github.com/lefes/curly-broccoli/config"
)

type DiscordSession struct {
	s *discordgo.Session
}

func (ds *DiscordSession) Start(dc *config.DiscordService) error {
	session, err := discordgo.New("Bot " + dc.BotToken)
	if err != nil {
		return fmt.Errorf("Error occured while creating discord session: %s", err)
	}
	ds.s = session
	err = ds.s.Open()
	if err != nil {
		return fmt.Errorf("Error occured while opening discord session: %s", err)
	}

	return nil
}

func (ds *DiscordSession) Stop() error {
	return ds.s.Close()
}

func (ds *DiscordSession) WatchMessages(handler func(*discordgo.MessageCreate)) {
	ds.s.AddHandler(func(s *discordgo.Session, m *discordgo.MessageCreate) {
		handler(m)
	})
}

func (ds *DiscordSession) WatchReactions(
	onAdd func(*discordgo.MessageReactionAdd) bool,
	onRemove func(*discordgo.MessageReactionRemove) bool,
) {
	ds.s.AddHandler(func(s *discordgo.Session, r *discordgo.MessageReactionAdd) {
		if onAdd != nil {
			onAdd(r)
		}
	})

	ds.s.AddHandler(func(s *discordgo.Session, r *discordgo.MessageReactionRemove) {
		if onRemove != nil {
			onRemove(r)
		}
	})
}

func (ds *DiscordSession) GetAllUsers(guildID string) (*DiscordMembers, error) {
	members, err := ds.s.GuildMembers(guildID, "", 1000)
	if err != nil {
		return nil, fmt.Errorf("Failed to get guild members: %v", err)
	}
	return &DiscordMembers{Members: members}, nil
}

func (ds *DiscordSession) GetSession() *discordgo.Session {
	return ds.s
}
