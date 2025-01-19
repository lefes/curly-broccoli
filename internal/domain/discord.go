package domain

import (
	"github.com/bwmarrin/discordgo"
)

type DiscordMembers struct {
	Members []*discordgo.Member
}

type Message struct {
	ID        string
	Username  string
	Content   string
	Author    string
	Channel   string
	ChannelID string
	Raw       *discordgo.MessageCreate
}

type SlashCommand struct {
	Name        string
	Description string
	Options     []*discordgo.ApplicationCommandOption
	Handler     func(s *discordgo.Session, i *discordgo.InteractionCreate)
}
