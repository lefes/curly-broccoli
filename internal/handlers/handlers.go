package handlers

import (
	"fmt"
	"time"

	"github.com/bwmarrin/discordgo"
	"github.com/lefes/curly-broccoli/internal/domain"
	"github.com/lefes/curly-broccoli/internal/services"
)

func HandlePoints(msg *domain.Message, svc *services.Services, session *discordgo.Session) bool {
	if msg.Raw.Author.Bot {
		return false
	}

	if len(msg.Content) < 5 {
		return false
	}

	activity := svc.UserActivity.AddOrUpdateUser(msg.Author)

	now := time.Now()
	fmt.Println("Last message time:", activity.LastMessageTime)
	fmt.Println("Next allowed message time:", activity.NextMessageTime)

	if now.Before(activity.NextMessageTime) {
		fmt.Println("Message too soon")
		return false
	}

	activity.LastMessageTime = now
	activity.NextMessageTime = now.Add(2 * time.Second)
	activity.MessageCount += 1

	if activity.MessageCount >= svc.UserActivity.GetMaxMessages() {
		fmt.Println("User reached daily limit")
		svc.UserActivity.MarkLimitReached(msg.Author)
		err := svc.User.UpdateUserPoints(msg.Author, 25)
		if err != nil {
			fmt.Printf("Error updating user points in database: %s\n", err)
			return false
		}
		return true
	}

	err := svc.User.UpdateUserDailyMessages(msg.Author, activity.MessageCount)
	if err != nil {
		fmt.Printf("Error updating user points in database: %s\n", err)
		return false
	}

	fmt.Printf("User %s sent a valid message. Points updated.\n", msg.Author)
	return true
}
