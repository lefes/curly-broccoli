package domain

import "github.com/bwmarrin/discordgo"

type Message struct {
	ID        string
	Username  string
	Content   string
	Author    string
	Channel   string
	ChannelID string
	Raw       *discordgo.MessageCreate
}
