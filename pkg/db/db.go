package db

import (
	"database/sql"
	_ "modernc.org/sqlite"
	"os"
)

var dbConn *sql.DB

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

func Init() error {
	dbFile := "./scheduler.db"
	_, err := os.Stat(dbFile)

	var install bool
	if err != nil {
		install = true
	}

	if install {
		dbConn, err = open()

		defer dbConn.Close()

		if err != nil {
			return err
		}
		_, err = dbConn.Exec(schema)
		if err != nil {
			return err
		}
	}

	return nil
}

func open() (*sql.DB, error) {
	dbFile := "./scheduler.db"
	envFile := os.Getenv("TODO_DBFILE")
	if len(envFile) > 0 {
		dbFile = envFile
	}
	db, err := sql.Open("sqlite", dbFile)
	if err != nil {
		return nil, err
	}
	return db, nil
}
