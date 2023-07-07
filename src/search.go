package main

import (
	"encoding/json"
	"io/ioutil"
	"log"
	"net/http"
	"net/url"
)

func randomQuote() string {
	response, err := makeRequest("https://animechan.xyz/api/random")
	if err != nil {
		log.Fatal("Request failed: ", err)
	}
	defer response.Body.Close()

	body, err := getResponseBody(response)
	if err != nil {
		log.Fatal("Failed to read response body:", err)
	}

	return string(body)
}

func makeRequest(url string) (*http.Response, error) {
	response, err := http.Get(url)
	if err != nil {
		return nil, err
	}
	return response, nil
}

func getResponseBody(response *http.Response) ([]byte, error) {
	body, err := ioutil.ReadAll(response.Body)
	if err != nil {
		return nil, err
	}
	return body, nil
}

func quoteByTitle(title string) string {
	db, _ := NewDatabase()
	quotes, err := db.GetQuotesByTitle(title)
	if err != nil {
		log.Fatal("Failed to get quote", err)
	}

	if quotes == nil {
		log.Printf("MISS: not found in DB %s", title)
		encodedTitle := url.QueryEscape(title)
		// TODO: if the quotes aren't in DB, they are returned without IDs
		quotes = fetchQuotesByTitle(encodedTitle)
		saveQuotesToDatabase(quotes)
	} else {
		log.Printf("HIT: found in DB %s", title)
	}

	quotesJSON, err := json.Marshal(quotes)
	if err != nil {
		log.Fatal("Failed to serialize quotes to JSON:", err)
	}

	return string(quotesJSON)
}

func fetchQuotesByTitle(encodedTitle string) []AnimeQuote {
	response, err := makeRequest("https://animechan.xyz/api/quotes/anime?title=" + encodedTitle)
	if err != nil {
		log.Fatal("Request failed:", err)
	}
	defer response.Body.Close()

	body, err := getResponseBody(response)
	if err != nil {
		log.Fatal("Failed to read response body:", err)
	}

	var quotes []AnimeQuote
	err = json.Unmarshal(body, &quotes)
	if err != nil {
		log.Fatal("Failed to parse API response:", err)
	}

	return quotes
}

func saveQuotesToDatabase(quotes []AnimeQuote) {
	// TODO: improve database connection
	db, _ := NewDatabase()

	for _, quote := range quotes {
		apiQuote := AnimeQuote{
			Name:      quote.Name,
			Character: quote.Character,
			Quote:     quote.Quote,
		}
		_, err := db.InsertQuote(apiQuote)
		if err != nil {
			log.Fatal("Failed to insert quote into the database:", err)
		}
	}
}
