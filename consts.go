package main

import (
	"flag"
)

var (
	GuildID        = "147313959819542528"
	discordUsers   = make(map[string]bool)
	dbUsers        = make(map[string]bool)
	RemoveCommands = flag.Bool("rmcmd", true, "Remove all commands after shutdowning or not")
)

// Убери это в config
