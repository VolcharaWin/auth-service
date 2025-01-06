package storage

import (
	"database/sql"

	_ "github.com/mattn/go-sqlite3"
)

func CreateDB() (*sql.DB, error) {
	db, err := sql.Open("sqlite3", "storage/user.db")
	if err != nil {
		return nil, err
	}
	stmt, err := db.Prepare(`
	CREATE TABLE IF NOT EXISTS user(
		login TEXT PRIMARY KEY,
		password TEXT NOT NULL
		);
	`)
	if err != nil {
		return nil, err
	}
	_, err = stmt.Exec()
	if err != nil {
		return nil, err
	}
	return db, err
}
