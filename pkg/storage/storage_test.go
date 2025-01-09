package storage

import (
	"os"
	"testing"

	"github.com/lefes/curly-broccoli/pkg/logging"
)

const testDBPath = "test.db"

func TestMain(m *testing.M) {
	logging.InitLogger()
	InitLogger()

	code := m.Run()

	_ = os.Remove(testDBPath)

	os.Exit(code)
}

func TestInitDB(t *testing.T) {
	_ = os.Remove(testDBPath)

	db, err := InitDB(testDBPath)
	if err != nil {
		t.Fatalf("Failed to initialize database: %v", err)
	}
	defer func() {
		_ = db.Close()
		_ = os.Remove(testDBPath)
	}()

	if err := db.Ping(); err != nil {
		t.Fatalf("Database ping failed: %v", err)
	}
}

func TestCreateTables(t *testing.T) {
	_ = os.Remove(testDBPath)

	db, err := InitDB(testDBPath)
	if err != nil {
		t.Fatalf("Failed to initialize database: %v", err)
	}
	defer func() {
		_ = db.Close()
		_ = os.Remove(testDBPath)
	}()

	if err := CreateTables(); err != nil {
		t.Fatalf("Failed to create tables: %v", err)
	}

	tables := []string{"users", "roles", "transactions", "lotteries", "lottery_tickets"}
	for _, table := range tables {
		query := `SELECT name FROM sqlite_master WHERE type='table' AND name=?;`
		row := db.QueryRow(query, table)

		var tableName string
		if err := row.Scan(&tableName); err != nil {
			t.Errorf("Table %s not found: %v", table, err)
		} else if tableName != table {
			t.Errorf("Expected table %s but found %s", table, tableName)
		}
	}
}
