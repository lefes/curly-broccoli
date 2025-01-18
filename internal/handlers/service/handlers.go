package service_handlers

import (
	"github.com/lefes/curly-broccoli/internal/domain"
	domain_handler "github.com/lefes/curly-broccoli/internal/handlers/domain"
)

type MessageHandlerService struct {
	handlers []func(*domain.Message, *domain_handler.HandlerContext) bool
}

func NewMessageHandlerService() *MessageHandlerService {
	return &MessageHandlerService{
		handlers: []func(*domain.Message, *domain_handler.HandlerContext) bool{},
	}
}

func (m *MessageHandlerService) HandleMessage(msg *domain.Message, ctx *domain_handler.HandlerContext) {
	for _, handler := range m.handlers {
		if handler(msg, ctx) {
			return
		}
	}
}

func (m *MessageHandlerService) AddHandler(handler func(*domain.Message, *domain_handler.HandlerContext) bool) {
	m.handlers = append(m.handlers, handler)
}
