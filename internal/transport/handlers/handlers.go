package handlers

import (
	"fmt"

	"github.com/bwmarrin/discordgo"
	"github.com/lefes/curly-broccoli/internal/services"
	"github.com/lefes/curly-broccoli/internal/transport/discordapi"
)

type CommandHandlers struct {
	services *services.Services
	dSession *discordapi.DiscordSession
}

func NewCommandHandlers(services *services.Services, dSession *discordapi.DiscordSession) *CommandHandlers {
	return &CommandHandlers{
		services: services,
		dSession: dSession,
	}
}

func (cmdH *CommandHandlers) HandleWeatherCommand(s *discordgo.Session, i *discordgo.InteractionCreate) {
	options := i.ApplicationCommandData().Options

	city := ""
	days := 1
	for _, opt := range options {
		switch opt.Name {
		case "city":
			city = opt.StringValue()
		case "days":
			days = int(opt.IntValue())
		}
	}

	weather, err := cmdH.services.Weather.GetWeather(city, days)
	if err != nil {
		errMsg := fmt.Sprintf("Failed to fetch weather data: %v", err)
		s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
			Type: discordgo.InteractionResponseChannelMessageWithSource,
			Data: &discordgo.InteractionResponseData{
				Content: errMsg,
			},
		})
		return
	}

	err = s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
		Type: discordgo.InteractionResponseChannelMessageWithSource,
		Data: &discordgo.InteractionResponseData{
			Embeds: []*discordgo.MessageEmbed{weather},
		},
	})
	if err != nil {
		fmt.Printf("Failed to send response: %v\n", err)
	}
}
