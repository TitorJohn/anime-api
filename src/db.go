package main

import (
	"database/sql"
	"fmt"
	"os"
	"sync"

	_ "github.com/mattn/go-sqlite3"
)

type AnimeQuote struct {
	ID        int    `json:"id"`
	Name      string `json:"anime"`
	Character string `json:"character"`
	Quote     string `json:"quote"`
}

var (
	db   *sql.DB
	once sync.Once
)

func GetDbInstance() (*sql.DB, error) {
	filePath := os.Getenv("DB_PATH")

	var err error
	once.Do(func() {
		db, err = sql.Open("sqlite3", filePath)
		if err != nil {
			fmt.Println("Error connecting to the database:", err)
		}
	})

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
	}

	return db, err

}

func getQuotesByTitle(title string) ([]AnimeQuote, error) {
	db, err := GetDbInstance()
	if err != nil {
		return []AnimeQuote{}, err
	}

	query := "SELECT id, name, character, quote FROM anime_quote where name=?"
	rows, err := db.Query(query, title)
	if err != nil {
		return nil, err
	}

	defer rows.Close()

	var quotes []AnimeQuote

	for rows.Next() {
		var quote AnimeQuote
		err := rows.Scan(&quote.ID, &quote.Name, &quote.Character, &quote.Quote)
		if err != nil {
			return nil, err
		}
		quotes = append(quotes, quote)
	}

	if err := rows.Err(); err != nil {
		return nil, err
	}

	return quotes, nil

}

func insertQuote(quote AnimeQuote) (int64, error) {
	db, err := GetDbInstance()
	if err != nil {
		return 0, err
	}

	stmt, err := db.Prepare("INSERT INTO anime_quote (name, character, quote) VALUES (?, ?, ?)")
	if err != nil {
		return 0, err
	}
	defer stmt.Close()

	result, err := stmt.Exec(quote.Name, quote.Character, quote.Quote)
	if err != nil {
		return 0, err
	}

	lastInsertID, err := result.LastInsertId()
	if err != nil {
		return 0, err
	}

	return lastInsertID, nil
}
