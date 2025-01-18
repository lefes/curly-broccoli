package logging

import (
	"github.com/sirupsen/logrus"
)

var Logger *logrus.Logger

func InitLogger() {
	Logger = logrus.New()

	Logger.SetFormatter(&logrus.TextFormatter{
		FullTimestamp: true,
	})
	Logger.SetLevel(logrus.InfoLevel)
}

func GetLogger(module string) *logrus.Entry {
	return Logger.WithFields(logrus.Fields{
		"module": module,
	})
}
