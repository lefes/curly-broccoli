package main

import (
	"flag"
	"os"
	"os/signal"

	"github.com/bwmarrin/discordgo"
	"github.com/lefes/curly-broccoli/config"
	"github.com/lefes/curly-broccoli/internal/cron"
	"github.com/lefes/curly-broccoli/internal/domain"
	"github.com/lefes/curly-broccoli/internal/logging"
	"github.com/lefes/curly-broccoli/internal/repository"
	"github.com/lefes/curly-broccoli/internal/services"
	"github.com/lefes/curly-broccoli/internal/storage"
	"github.com/lefes/curly-broccoli/internal/transport/discordapi"
	"github.com/lefes/curly-broccoli/internal/transport/handlers"
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

	db, err := storage.InitDB(mainConfig.Storage.DbPath)
	if err != nil {
		logger.Fatalf("Failed to initialize database: %v", err)
	}
	logger.Info("Database connection initialized")

	err = dSession.Start(&mainConfig.Discord)
	if err != nil {
		logger.Errorf("Failed to open discord session: %v", err)
	}

	repo := repository.NewRepository(db)
	services := services.NewServices(repo, mainConfig, &dSession)
	if err != nil {
		logger.Errorf("Failed to sync users: %v", err)
	}
	cronService := cron.NewCronService(services)
	cronService.Start()

	handlers := handlers.NewCommandHandlers(services, repo)

	commands := handlers.CommandsInit()
	registeredCommands, err := dSession.RegisterCommands(commands, mainConfig.Discord.GuildID)
	if err != nil {
		logger.Fatalf("Failed to register commands: %v", err)
	}

	dSession.WatchMessages(func(m *discordgo.MessageCreate) { // унести туда же где и обработчики команд
		msg := domain.Message{
			Raw:       m,
			ID:        m.ID,
			Username:  m.Author.Username,
			Content:   m.Content,
			Author:    m.Author.ID,
			Channel:   m.ChannelID,
			ChannelID: m.ChannelID,
		}
		handlers.HandleMessagePoints(&msg)
	})

	stop := make(chan os.Signal, 1)
	signal.Notify(stop, os.Interrupt)
	logger.Info("Press Ctrl+C to exit")
	<-stop

	// Logic after bot has been stopped with Ctrl+C
	dSession.Stop()
	if *RemoveCommands {
		logger.Println("Removing commands...")
		err := dSession.DeleteCommands(registeredCommands, mainConfig.Discord.GuildID)
		if err != nil {
			logger.Fatalf("Failed to remove commands: %v", err)
		}
	}
}
