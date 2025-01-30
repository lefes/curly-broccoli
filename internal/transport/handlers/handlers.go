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
			Name:        "points",
			Description: "Add or remove points",
			Options: []*discordgo.ApplicationCommandOption{
				{
					Name:        "add",
					Description: "Add points to a user",
					Type:        discordgo.ApplicationCommandOptionSubCommand,
					Options: []*discordgo.ApplicationCommandOption{
						{
							Name:        "user",
							Description: "The user to give points to",
							Type:        discordgo.ApplicationCommandOptionUser,
							Required:    true,
						},
						{
							Name:        "points",
							Description: "The number of points to add",
							Type:        discordgo.ApplicationCommandOptionInteger,
							Required:    true,
						},
					},
				},
				{
					Name:        "remove",
					Description: "Remove points from a user",
					Type:        discordgo.ApplicationCommandOptionSubCommand,
					Options: []*discordgo.ApplicationCommandOption{
						{
							Name:        "user",
							Description: "The user to remove points from",
							Type:        discordgo.ApplicationCommandOptionUser,
							Required:    true,
						},
						{
							Name:        "points",
							Description: "The number of points to remove",
							Type:        discordgo.ApplicationCommandOptionInteger,
							Required:    true,
						},
					},
				},
			},
			Handler: cmdH.HandlePointsCommand,
		},
		{
			Name:        "хтоя",
			Description: "Get you current role",
			Handler:     cmdH.HandleWhoAmICommand,
		},
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
		{
			Name:        "respect",
			Description: "Add or remove respect points",
			Options: []*discordgo.ApplicationCommandOption{
				{
					Name:        "add",
					Description: "Add respect points to a user",
					Type:        discordgo.ApplicationCommandOptionSubCommand,
					Options: []*discordgo.ApplicationCommandOption{
						{
							Name:        "user",
							Description: "The user to give respect points to",
							Type:        discordgo.ApplicationCommandOptionUser,
							Required:    true,
						},
						{
							Name:        "points",
							Description: "The number of respect points to add",
							Type:        discordgo.ApplicationCommandOptionInteger,
							Required:    true,
						},
					},
				},
				{
					Name:        "remove",
					Description: "Remove respect points from a user",
					Type:        discordgo.ApplicationCommandOptionSubCommand,
					Options: []*discordgo.ApplicationCommandOption{
						{
							Name:        "user",
							Description: "The user to remove respect points from",
							Type:        discordgo.ApplicationCommandOptionUser,
							Required:    true,
						},
						{
							Name:        "points",
							Description: "The number of respect points to remove",
							Type:        discordgo.ApplicationCommandOptionInteger,
							Required:    true,
						},
					},
				},
			},
			Handler: cmdH.HandleRespectCommand,
		},
	}
	return commands
}

func (cmdH *CommandHandlers) HandlePointsCommand(s *discordgo.Session, i *discordgo.InteractionCreate) {
	if !cmdH.services.User.IsAdmin(i.Member.User.ID) {
		s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
			Type: discordgo.InteractionResponseChannelMessageWithSource,
			Data: &discordgo.InteractionResponseData{
				Content: "You don't have permission to use this command.",
			},
		})
		return
	}

	options := i.ApplicationCommandData().Options

	if len(options) == 0 {
		s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
			Type: discordgo.InteractionResponseChannelMessageWithSource,
			Data: &discordgo.InteractionResponseData{
				Content: "Invalid command usage. Please use /points add or /points remove.",
			},
		})
		return
	}

	subcommand := options[0]
	var userID string
	var points int64

	for _, opt := range subcommand.Options {
		switch opt.Name {
		case "user":
			userID = opt.UserValue(nil).ID
		case "points":
			points = opt.IntValue()
		}
	}

	if userID == "" || points == 0 {
		s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
			Type: discordgo.InteractionResponseChannelMessageWithSource,
			Data: &discordgo.InteractionResponseData{
				Content: "Please specify a user and the number of points.",
			},
		})
		return
	}

	switch subcommand.Name {
	case "add":
		err := cmdH.repo.User.AddUserPoints(userID, int(points))
		cmdH.repo.User.AddDayPoints(userID, int(points))
		if err != nil {
			s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
				Type: discordgo.InteractionResponseChannelMessageWithSource,
				Data: &discordgo.InteractionResponseData{
					Content: fmt.Sprintf("Failed to add respect points: %v", err),
				},
			})
			return
		}
		s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
			Type: discordgo.InteractionResponseChannelMessageWithSource,
			Data: &discordgo.InteractionResponseData{
				Content: fmt.Sprintf("Added %d points to <@%s>", points, userID),
			},
		})

	case "remove":
		err := cmdH.repo.User.RemoveUserPoints(userID, int(points))
		if err != nil {
			s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
				Type: discordgo.InteractionResponseChannelMessageWithSource,
				Data: &discordgo.InteractionResponseData{
					Content: fmt.Sprintf("Failed to remove points: %v", err),
				},
			})
			return
		}
		cmdH.repo.User.RemoveDayPoints(userID, int(points))
		s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
			Type: discordgo.InteractionResponseChannelMessageWithSource,
			Data: &discordgo.InteractionResponseData{
				Content: fmt.Sprintf("Removed %d points from <@%s>", points, userID),
			},
		})

	default:
		s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
			Type: discordgo.InteractionResponseChannelMessageWithSource,
			Data: &discordgo.InteractionResponseData{
				Content: "Invalid subcommand. Use /points add or /points remove.",
			},
		})
	}
}

func (cmdH *CommandHandlers) HandleWhoAmICommand(s *discordgo.Session, i *discordgo.InteractionCreate) {
	userID := i.Member.User.ID
	role, err := cmdH.services.Roles.GetUserRole(userID)
	if err != nil {
		s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
			Type: discordgo.InteractionResponseChannelMessageWithSource,
			Data: &discordgo.InteractionResponseData{
				Content: "Failed to get user role",
			},
		})
	}

	s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
		Type: discordgo.InteractionResponseChannelMessageWithSource,
		Data: &discordgo.InteractionResponseData{
			Content: fmt.Sprintf("<@%s> has the role %s", userID, role.Name),
		},
	})
}

func (cmdH *CommandHandlers) HandleRespectCommand(s *discordgo.Session, i *discordgo.InteractionCreate) {
	if !cmdH.services.User.IsAdmin(i.Member.User.ID) {
		s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
			Type: discordgo.InteractionResponseChannelMessageWithSource,
			Data: &discordgo.InteractionResponseData{
				Content: "You don't have permission to use this command.",
			},
		})
		return
	}
	options := i.ApplicationCommandData().Options

	if len(options) == 0 {
		s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
			Type: discordgo.InteractionResponseChannelMessageWithSource,
			Data: &discordgo.InteractionResponseData{
				Content: "Invalid command usage. Please use /respect add or /respect remove.",
			},
		})
		return
	}

	subcommand := options[0]
	var userID string
	var points int64

	for _, opt := range subcommand.Options {
		switch opt.Name {
		case "user":
			userID = opt.UserValue(nil).ID
		case "points":
			points = opt.IntValue()
		}
	}

	if userID == "" || points == 0 {
		s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
			Type: discordgo.InteractionResponseChannelMessageWithSource,
			Data: &discordgo.InteractionResponseData{
				Content: "Please specify a user and the number of respect points.",
			},
		})
		return
	}

	// TODO: MERGE PROMOTION CHECKING WITH ADDING RESPECT It SHOULD BE INSIDE ADDUSERRESPECT
	switch subcommand.Name {
	case "add":
		promoted, newRole := cmdH.services.Roles.WillGetPromotion(userID, int(points))
		err := cmdH.repo.Roles.AddUserRespect(userID, int(points))
		if err != nil {
			s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
				Type: discordgo.InteractionResponseChannelMessageWithSource,
				Data: &discordgo.InteractionResponseData{
					Content: fmt.Sprintf("Failed to add respect points: %v", err),
				},
			})
			return
		}
		err = cmdH.repo.Roles.AddDayUserRespect(userID, int(points))
		if err != nil {
			cmdH.logger.Errorf("Failed to add daily respect points: %v", err)
		}
		s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
			Type: discordgo.InteractionResponseChannelMessageWithSource,
			Data: &discordgo.InteractionResponseData{
				Content: fmt.Sprintf("Added %d respect points to <@%s>", points, userID),
			},
		})
		if promoted {
			err = cmdH.repo.Roles.UpdateUserRole(userID, newRole.ID)
			if err != nil {
				cmdH.logger.Errorf("Failed to update user role: %v", err)
			}
			s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
				Type: discordgo.InteractionResponseChannelMessageWithSource,
				Data: &discordgo.InteractionResponseData{
					Content: fmt.Sprintf("<@%s> has been promoted to %s!", userID, newRole.Name),
				},
			})
		}

	case "remove":
		err := cmdH.repo.Roles.RemoveUserRespect(userID, int(points))
		if err != nil {
			s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
				Type: discordgo.InteractionResponseChannelMessageWithSource,
				Data: &discordgo.InteractionResponseData{
					Content: fmt.Sprintf("Failed to remove respect points: %v", err),
				},
			})
			return
		}
		err = cmdH.repo.Roles.RemoveDayUserRespect(userID, int(points))
		if err != nil {
			cmdH.logger.Errorf("Failed to remove daily respect points: %v", err)
		}
		s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
			Type: discordgo.InteractionResponseChannelMessageWithSource,
			Data: &discordgo.InteractionResponseData{
				Content: fmt.Sprintf("Removed %d respect points from <@%s>", points, userID),
			},
		})

	default:
		s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
			Type: discordgo.InteractionResponseChannelMessageWithSource,
			Data: &discordgo.InteractionResponseData{
				Content: "Invalid subcommand. Use /respect add or /respect remove.",
			},
		})
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
