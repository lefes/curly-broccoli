package cron

import (
	"fmt"
	"time"

	"github.com/lefes/curly-broccoli/internal/services"
)

type CronService struct {
	services *services.Services
}

func NewCronService(services *services.Services) *CronService {
	return &CronService{
		services: services,
	}
}

func (c *CronService) Start() {
	go c.runDailyTasks()
	go c.runWeeklyTasks()
}

func (c *CronService) runDailyTasks() {
	for {
		if err := c.services.Discord.SyncUsers(); err != nil {
			fmt.Printf("Failed to sync users: %v\n", err)
		} else {
			fmt.Println("Users sync has been completed")
		}

		now := time.Now()
		nextDay := now.Add(24 * time.Hour).Truncate(24 * time.Hour)
		durationUntilNextDay := time.Until(nextDay)

		time.Sleep(durationUntilNextDay)

		fmt.Println("Running daily tasks...")

		if err := c.services.User.Reset(); err != nil {
			fmt.Printf("Failed to reset daily limits: %v\n", err)
		} else {
			fmt.Println("Daily limits reset successfully")
		}
	}
}

func (c *CronService) runWeeklyTasks() {
	// TODO: Implement weekly tasks
}
