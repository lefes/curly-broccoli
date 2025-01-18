package domain

import "time"

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
