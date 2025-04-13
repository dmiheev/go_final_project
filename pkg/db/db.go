package db

import (
	"database/sql"
	"fmt"
	"log"
	_ "modernc.org/sqlite"
	"os"
)

const (
	schema = `CREATE TABLE scheduler
                (
                    id      INTEGER PRIMARY KEY AUTOINCREMENT,
                    date    CHAR(8) NOT NULL DEFAULT "",
                    title   CHAR(255),
                    comment TEXT,
                    repeat  CHAR(128)
                );
              CREATE INDEX idx_scheduler_date ON scheduler (date);`
)

func Init() (*sql.DB, error) {
	dbFile := "./scheduler.db"
	envFile := os.Getenv("TODO_DBFILE")
	if len(envFile) > 0 {
		dbFile = envFile
	}

	_, err := os.Stat(dbFile)

	var install bool
	if err != nil {
		install = true
	}

	db, err := sql.Open("sqlite", dbFile)
	if err != nil {
		return nil, fmt.Errorf("error while open db: %w", err)
	}

	if install {
		if err = createTable(db); err != nil {
			return nil, err
		}
	}

	return db, nil
}

func createTable(db *sql.DB) error {
	if _, err := db.Exec(schema); err != nil {
		log.Fatalf("failed to create table: %v", err)
		return err
	}
	return nil
}
