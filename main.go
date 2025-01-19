package main

import (
	"flag"
	"os"
	"os/signal"

	"github.com/bwmarrin/discordgo"
	"github.com/lefes/curly-broccoli/config"
	"github.com/lefes/curly-broccoli/internal/domain"
	"github.com/lefes/curly-broccoli/internal/logging"
	"github.com/lefes/curly-broccoli/internal/repository"
	"github.com/lefes/curly-broccoli/internal/services"
	"github.com/lefes/curly-broccoli/internal/storage"
	"github.com/lefes/curly-broccoli/internal/transport/discordapi"
)

func init() { flag.Parse() }

func main() {
	logging.InitLogger()
	mainConfig := config.Init()
	dSession := discordapi.DiscordSession{}
	err := dSession.Start(&mainConfig.Discord)
	if err != nil {
		logging.GetLogger("bot").Fatalf("Failed to open discord session: %v", err)
	}
	logger := logging.GetLogger("bot")
	logger.Info("Starting application")

	db, err := storage.InitDB(dbPath)
	if err != nil {
		logger.Fatalf("Failed to initialize database: %v", err)
	}
	logger.Info("Database connection initialized")

	repo := repository.NewRepository(db)
	services := services.NewServices(repo, mainConfig)
	err = dSession.Start(&mainConfig.Discord)
	if err != nil {
		logger.Errorf("Failed to open discord session: %v", err)
	}

	handlers := discordapi.NewCommandHandlers(services, &dSession)
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
			Handler: handlers.HandleWeatherCommand,
		},
	}

	registeredCommands, err := dSession.RegisterCommands(commands, mainConfig.Discord.GuildID)
	if err != nil {
		logger.Fatalf("Failed to register commands: %v", err)
	}

	stop := make(chan os.Signal, 1)
	signal.Notify(stop, os.Interrupt)
	logger.Info("Press Ctrl+C to exit")
	<-stop
	dSession.Stop()

	if *RemoveCommands {
		logger.Println("Removing commands...")
		err := dSession.DeleteCommands(registeredCommands, mainConfig.Discord.GuildID)
		if err != nil {
			logger.Fatalf("Failed to remove commands: %v", err)
		}
	}
}
