package handlers

import (
	"fmt"
	"time"

	"github.com/bwmarrin/discordgo"
	"github.com/lefes/curly-broccoli/internal/domain"
	"github.com/lefes/curly-broccoli/internal/repository"
	"github.com/lefes/curly-broccoli/internal/services"
	"github.com/lefes/curly-broccoli/internal/transport/discordapi"
)

type CommandHandlers struct {
	services *services.Services
	dSession *discordapi.DiscordSession
	repo     *repository.Repositories
}

func NewCommandHandlers(services *services.Services, dSession *discordapi.DiscordSession, r *repository.Repositories) *CommandHandlers {
	return &CommandHandlers{
		services: services,
		dSession: dSession,
		repo:     r,
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

// Move business logic to discrord service
func (cmdH *CommandHandlers) HandlePoints(msg *domain.Message) bool {
	if msg.Raw.Author.Bot {
		return false
	}

	if len(msg.Content) < 5 {
		fmt.Println("Message too short")
		return false
	}

	activity := cmdH.repo.User.AddOrUpdateUserActivity(msg.Author)

	if cmdH.services.User.IsLimitReached(msg.Author) {
		fmt.Printf("User %s has reached the daily limit. Skipping.\n", msg.Username)
		return false
	}

	now := time.Now()

	if now.Before(activity.NextMessageTime) {
		return false
	}

	activity.LastMessageTime = now
	activity.NextMessageTime = now.Add(2 * time.Second)
	activity.MessageCount = activity.MessageCount + 1

	if activity.MessageCount >= cmdH.repo.User.GetMaxMessages() {
		fmt.Printf("User %s reached daily limit\n", msg.Username)
		cmdH.repo.User.MarkLimitReached(msg.Author)
		err := cmdH.repo.User.UpdateUserPoints(msg.Author, 25)
		if err != nil {
			fmt.Printf("Error updating user points in database: %s\n", err)
			return false
		}
		fmt.Printf("User %s received 25 points\n", msg.Username)
		return true
	}

	err := cmdH.repo.User.UpdateUserDailyMessages(msg.Author, activity.MessageCount)
	if err != nil {
		fmt.Printf("Error updating user daily messages in database: %s\n", err)
		return false
	}
	return true
}
