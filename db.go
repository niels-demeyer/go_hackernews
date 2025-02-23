package main

import (
	"database/sql"
	"fmt"

	_ "github.com/mattn/go-sqlite3"
)

func InitDB(dbPath string) (*sql.DB, error) {
    db, err := sql.Open("sqlite3", dbPath)
    if err != nil {
        return nil, err
    }

    // Create table if it doesn't exist with UNIQUE constraint on URL
    createTableSQL := `
    CREATE TABLE IF NOT EXISTS hackernews_items (
        id INTEGER PRIMARY KEY AUTOINCREMENT,
        title TEXT NOT NULL,
        site TEXT,
        url TEXT NOT NULL UNIQUE,
        created_at DATETIME DEFAULT CURRENT_TIMESTAMP
    );`

    _, err = db.Exec(createTableSQL)
    if err != nil {
        return nil, err
    }

    return db, nil
}

func InsertItem(db *sql.DB, title, site, url string) error {
    // Use INSERT OR IGNORE to handle duplicates
    insertSQL := `
    INSERT OR IGNORE INTO hackernews_items (title, site, url)
    VALUES (?, ?, ?)`

    result, err := db.Exec(insertSQL, title, site, url)
    if err != nil {
        return err
    }

    // Check if the row was actually inserted
    rowsAffected, err := result.RowsAffected()
    if err != nil {
        return err
    }

    if rowsAffected == 0 {
        return fmt.Errorf("duplicate URL: %s", url)
    }

    return nil
}