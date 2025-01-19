package main

import (
	"flag"
	"log"
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

func loggersInit() {
	logging.InitLogger()
	mLogger = logging.GetLogger("main")
	sLogger = logging.GetLogger("storage")
	dLogger = logging.GetLogger("discord")
	//wLogger := logging.GetLogger("weather")
}

func main() {
	loggersInit()
	mainConfig := config.Init()
	mLogger.Info("Starting application")

	db, err := storage.InitDB(dbPath)
	if err != nil {
		sLogger.Fatalf("Failed to initialize database: %v", err)
	}
	sLogger.Info("Database connection initialized")

	repo := repository.NewRepository(db)
	services := services.NewServices(repo, mainConfig)
	session, err := services.Discord.Open()
	if err != nil {
		dLogger.Errorf("Failed to open discord session: %v", err)
	}
	defer session.Close()
	/*  err, _ = services.Discord.BotRegister() */
	/* if err != nil { */
	/* dLogger.Fatalf("Failed to register bot: %v", err) */
	/* } */

	handlers := discordapi.NewCommandHandlers(services)
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

	registeredCommands, err := discordapi.RegisterCommands(session, commands, mainConfig.Discord.GuildID)
	if err != nil {
		dLogger.Fatalf("Failed to register commands: %v", err)
	}

	stop := make(chan os.Signal, 1)
	signal.Notify(stop, os.Interrupt)
	mLogger.Info("Press Ctrl+C to exit")
	<-stop

	if *RemoveCommands {
		mLogger.Println("Removing commands...")
		for _, v := range registeredCommands {
			err := session.ApplicationCommandDelete(session.State.User.ID, GuildID, v.ID)
			if err != nil {
				log.Panicf("Cannot delete '%v' command: %v", v.Name, err)
			}
		}
	}
}
