package services

import (
	"fmt"
	"time"

	"github.com/bwmarrin/discordgo"
	"github.com/lefes/curly-broccoli/internal/domain"
)

type MessageHandlerContext struct {
	Services *Services
	Session  *discordgo.Session
}

func NewHandlerContext(services *Services, session *discordgo.Session) *MessageHandlerContext {
	return &MessageHandlerContext{
		Services: services,
		Session:  session,
	}
}

type MessageHandlerService struct {
	handlers []func(*domain.Message, *MessageHandlerContext) bool
}

func NewMessageHandlerService() *MessageHandlerService {
	return &MessageHandlerService{
		handlers: []func(*domain.Message, *MessageHandlerContext) bool{},
	}
}

func (m *MessageHandlerService) HandleMessage(msg *domain.Message, ctx *MessageHandlerContext) {
	for _, handler := range m.handlers {
		if handler(msg, ctx) {
			return
		}
	}
}

func (m *MessageHandlerService) AddHandler(handler func(*domain.Message, *MessageHandlerContext) bool) {
	m.handlers = append(m.handlers, handler)
}

func (m *MessageHandlerService) HandlePoints(msg *domain.Message, ctx *MessageHandlerContext) bool {
	if msg.Raw.Author.Bot {
		return false
	}

	if len(msg.Content) < 5 {
		fmt.Println("Message too short")
		return false
	}

	activity := ctx.Services.User.AddOrUpdateUserActivity(msg.Author)

	if ctx.Services.User.IsLimitReached(msg.Author) {
		fmt.Printf("User %s has reached the daily limit. Skipping.\n", msg.Author)
		return false
	}

	now := time.Now()
	fmt.Println("Last message time:", activity.LastMessageTime)
	fmt.Println("Next allowed message time:", activity.NextMessageTime)

	if now.Before(activity.NextMessageTime) {
		fmt.Println("Message too soon")
		return false
	}

	activity.LastMessageTime = now
	activity.NextMessageTime = now.Add(2 * time.Second)
	activity.MessageCount = activity.MessageCount + 1

	if activity.MessageCount >= ctx.Services.User.GetMaxMessages() {
		fmt.Println("User reached daily limit")
		ctx.Services.User.MarkLimitReached(msg.Author)
		err := ctx.Services.User.UpdateUserPoints(msg.Author, 25)
		if err != nil {
			fmt.Printf("Error updating user points in database: %s\n", err)
			return false
		}
		return true
	}

	err := ctx.Services.User.UpdateUserDailyMessages(msg.Author, activity.MessageCount)
	if err != nil {
		fmt.Printf("Error updating user daily messages in database: %s\n", err)
		return false
	}

	fmt.Printf("User %s sent a valid message. Daily messages updated.\n", msg.Author)
	return true
}
