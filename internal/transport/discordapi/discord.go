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
