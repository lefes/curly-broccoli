package storage

import (
	"database/sql"
	"fmt"

	"github.com/lefes/curly-broccoli/pkg/logging"
	_ "github.com/mattn/go-sqlite3"
	"github.com/sirupsen/logrus"
)

var logger *logrus.Entry

var db *sql.DB

func InitLogger() {
	if logger == nil {
		logger = logging.GetLogger("storage")
	}
}

func InitDB(filePath string) (*sql.DB, error) {
	var err error
	db, err = sql.Open("sqlite3", filePath)
	if err != nil {
		return nil, fmt.Errorf("error initializing database: %v", err)
	}

	if err := db.Ping(); err != nil {
		return nil, fmt.Errorf("error connecting to database: %v", err)
	}

	logger.Info("Database connection initialized.")
	return db, nil
}
