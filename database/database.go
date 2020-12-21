package database

import (
	"database/sql"
	"log"

	_ "github.com/mattn/go-sqlite3"
)

var db *sql.DB

func Init(databaseFile string) error {
	log.Println("Opening sqlite3 database")
	d, err := sql.Open("sqlite3", databaseFile)
	if err != nil {
		return err
	}
	db = d
	if err := initTables(); err != nil {
		return err
	}
	return nil
}

func Close() {
	if db != nil {
		err := db.Close()
		if err != nil {
			log.Println("Error closing database:", err)
		}
	}
	log.Println("Database closed")
}
