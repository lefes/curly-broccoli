package main

import (
	"fmt"
	"os"

	"github.com/bwmarrin/discordgo"
	"github.com/lefes/curly-broccoli/internal/bot"
	"github.com/lefes/curly-broccoli/internal/config"
	"github.com/lefes/curly-broccoli/internal/logging"
	"github.com/lefes/curly-broccoli/pkg/weather"
)

func main() {
	cfg, err := config.Load()
	if err != nil {
		fmt.Println("config error:", err)
		os.Exit(1)
	}

	logging.Init(cfg.LogLevel)
	log := logging.Get("main")

	session, err := discordgo.New("Bot " + cfg.Token)
	if err != nil {
		log.Error("discord session:", err)
		os.Exit(1)
	}
	session.Identify.Intents = discordgo.IntentsGuildMessages

	weatherClient := weather.NewClient(cfg.WeatherAPIKey, cfg.WeatherBaseURL)
	r := bot.NewRouter(weatherClient)
	session.AddHandler(r.Handler())

	if err := session.Open(); err != nil {
		log.Error("open connection:", err)
		os.Exit(1)
	}
	log.Info("Bot is now running. Press CTRL-C to exit.")
	<-make(chan struct{})
}
