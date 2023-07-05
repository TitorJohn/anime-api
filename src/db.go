package main

import (
	"database/sql"
	"fmt"
	"os"

	_ "github.com/mattn/go-sqlite3"
)

type AnimeQuote struct {
	id        int
	name      string
	character string
	quote     string
}

func db() {
	filePath := os.Getenv("DB_PATH")
	db, err := sql.Open("sqlite3", filePath)

	if err != nil {
		fmt.Println("Error connecting to the database:", err)
		return
	}
	defer db.Close()

	createTableSQL := `
CREATE TABLE IF NOT EXISTS anime_quote (
	id INTEGER PRIMARY KEY AUTOINCREMENT,
	name TEXT,
	character TEXT,
	quote TEXT
)`
	_, err = db.Exec(createTableSQL)
	if err != nil {
		fmt.Println("Error creating table:", err)
		return
	}
}
