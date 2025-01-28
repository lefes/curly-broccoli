package handlers

import (
	"fmt"

	"github.com/bwmarrin/discordgo"
	"github.com/lefes/curly-broccoli/internal/domain"
	"github.com/lefes/curly-broccoli/internal/logging"
	"github.com/lefes/curly-broccoli/internal/repository"
	"github.com/lefes/curly-broccoli/internal/services"
)

type CommandHandlers struct {
	services *services.Services
	repo     *repository.Repositories // вот это нужно изи убрать
	dSession *discordgo.Session
	logger   *logging.Logger
}

func NewCommandHandlers(services *services.Services, r *repository.Repositories, s *discordgo.Session, l *logging.Logger) *CommandHandlers {
	return &CommandHandlers{
		services: services,
		repo:     r,
		dSession: s,
		logger:   l,
	}
}

func (cmdH *CommandHandlers) CommandsInit() []domain.SlashCommand {
	minValue := float64(1)
	maxValue := float64(7)
	commands := []domain.SlashCommand{
		{
			Name:        "weather",
			Description: "Get weather information for a city",
			Options: []*discordgo.ApplicationCommandOption{
				{
					Name:        "city",
					Description: "City to get weather information for",
					Type:        discordgo.ApplicationCommandOptionString,
					Required:    true,
				},
				{
					Name:        "days",
					Description: "Number of days for the forecast (default: 1)",
					Type:        discordgo.ApplicationCommandOptionInteger,
					Required:    false,
					MinValue:    &minValue,
					MaxValue:    maxValue,
				},
			},
			Handler: cmdH.HandleWeatherCommand,
		},
	}
	return commands
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
		cmdH.logger.Errorf("Failed to send response: %v", err)
	}
}

func (cmdH *CommandHandlers) HandleMessagePoints(msg *domain.Message) bool {
	if msg.Raw.Author.Bot {
		return false
	}

	activity, canSend := cmdH.services.User.CanSendMessage(msg)
	if !canSend {
		return false
	}

	cmdH.services.User.IncrementUserMessageCount(activity)

	err := cmdH.repo.User.UpdateUserDailyMessages(msg.Author, activity.MessageCount)
	if err != nil {
		cmdH.logger.Errorf("Error updating user daily messages in database: %s", err)
		return false
	}
	return true
}

func (cmdH *CommandHandlers) HandleReactionPointsAdd(r *discordgo.MessageReactionAdd) bool {

	message, err := cmdH.dSession.ChannelMessage(r.ChannelID, r.MessageID)
	if err != nil {
		fmt.Printf("Failed to get message: %v\n", err)
		return false
	}

	if !cmdH.services.Discord.IsValidReaction(message, r.UserID) {
		return false
	}

	if !cmdH.services.User.ReactionPoints(message) {
		return false
	}

	return true
}

func (cmdH *CommandHandlers) HandleReactionPointsRemove(r *discordgo.MessageReactionRemove) bool {

	message, err := cmdH.dSession.ChannelMessage(r.ChannelID, r.MessageID)
	if err != nil {
		cmdH.logger.Errorf("Failed to get message: %v", err)
		return false
	}
	if !cmdH.services.User.ReactionPointsRemoval(message) {
		return false
	}

	return true
}
