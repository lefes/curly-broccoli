package logging

import (
	"strings"

	"github.com/sirupsen/logrus"
)

var baseLogger *logrus.Logger

func Init(level string) {
	baseLogger = logrus.New()
	baseLogger.SetFormatter(&logrus.TextFormatter{FullTimestamp: true})

	switch strings.ToLower(level) {
	case "debug":
		baseLogger.SetLevel(logrus.DebugLevel)
	case "warn":
		baseLogger.SetLevel(logrus.WarnLevel)
	case "error":
		baseLogger.SetLevel(logrus.ErrorLevel)
	default:
		baseLogger.SetLevel(logrus.InfoLevel)
	}
}

func Get(module string) *logrus.Entry {
	if baseLogger == nil {
		Init("info")
	}
	return baseLogger.WithField("module", module)
}
