package domain

import "github.com/bwmarrin/discordgo"

type DiscordMembers struct {
	Members []*discordgo.Member
}
