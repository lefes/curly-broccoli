package service_handlers

import (
	"github.com/lefes/curly-broccoli/internal/domain"
	domain_handler "github.com/lefes/curly-broccoli/internal/handlers/domain"
)

type MessageHandler interface {
	AddHandler(handler func(*domain.Message, *domain_handler.HandlerContext) bool)
	HandleMessage(msg *domain.Message, ctx *domain_handler.HandlerContext)
}

type MessageService struct {
	MessageHandler MessageHandler
}

func NewMessageHandler() *MessageService {
	messageHandlerService := NewMessageHandlerService()
	return &MessageService{
		MessageHandler: messageHandlerService,
	}
}
