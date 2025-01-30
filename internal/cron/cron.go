package cron

import (
	"time"

	"github.com/lefes/curly-broccoli/internal/logging"
	"github.com/lefes/curly-broccoli/internal/services"
)

type CronService struct {
	services *services.Services
	logger   *logging.Logger
}

func NewCronService(services *services.Services, logger *logging.Logger) *CronService {
	return &CronService{
		services: services,
		logger:   logger,
	}
}

func (c *CronService) Start() {
	go c.runDailyTasks()
	go c.runWeeklyTasks()
}

func (c *CronService) runDailyTasks() {
	for {
		if err := c.services.Discord.SyncUsers(); err != nil {
			c.logger.Infof("Failed to sync users: %v\n", err)
		} else {
			c.logger.Info("Users sync has been completed")
		}

		now := time.Now()
		nextDay := now.Add(24 * time.Hour).Truncate(24 * time.Hour)
		durationUntilNextDay := time.Until(nextDay)

		time.Sleep(durationUntilNextDay)

		c.logger.Info("Running daily tasks...")

		if err := c.services.User.Reset(); err != nil {
			c.logger.Infof("Failed to reset daily limits: %v\n", err)
		} else {
			c.logger.Info("Daily limits reset successfully")
		}

		if err := c.services.Discord.RespectActiveUsers(); err != nil {
			c.logger.Infof("Failed to respect active users: %v\n", err)
		} else {
			c.logger.Info("Added daily user respect to active users")
		}

	}
}

func (c *CronService) runWeeklyTasks() {
	// TODO: Implement weekly tasks
}
