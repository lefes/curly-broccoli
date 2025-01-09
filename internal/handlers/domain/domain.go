package domain_handler

import (
	"github.com/bwmarrin/discordgo"
	"github.com/lefes/curly-broccoli/internal/services"
)

type HandlerContext struct {
	Services *services.Services
	Session  *discordgo.Session
}

func NewHandlerContext(services *services.Services, session *discordgo.Session) *HandlerContext {
	return &HandlerContext{
		Services: services,
		Session:  session,
	}
}
