package domain

import (
	"sync"
	"time"
)

type User struct {
	ID            int
	DiscordID     string
	Username      string
	RoleID        int
	Points        int
	Respect       int
	DailyMessages int
	LastActivity  time.Time
}

type UserActivity struct {
	UserID          string
	LastMessageTime time.Time
	NextMessageTime time.Time
	MessageCount    int
}

type UserActivities struct {
	Mu              sync.Mutex
	Activities      map[string]*UserActivity
	LimitReachedIDs map[string]bool
	MaxMessages     int
}

func NewUserActivities(maxMessages int) *UserActivities {
	return &UserActivities{
		Activities:      make(map[string]*UserActivity),
		LimitReachedIDs: make(map[string]bool),
		MaxMessages:     maxMessages,
	}
}
