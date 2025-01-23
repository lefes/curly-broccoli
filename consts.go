package main

import (
	"flag"

	"github.com/sirupsen/logrus"
)

var (
	GuildID        = "147313959819542528"
	discordUsers   = make(map[string]bool)
	dbUsers        = make(map[string]bool)
	logger         *logrus.Entry
	RemoveCommands = flag.Bool("rmcmd", true, "Remove all commands after shutdowning or not")
)

// Убери это в config
