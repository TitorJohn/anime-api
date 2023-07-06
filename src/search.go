package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"net/url"
)

func randomQuote() string {
	response, err := http.Get("https://animechan.xyz/api/random")

	if err != nil {
		log.Fatal("Request failed: ", err)
	}
	defer response.Body.Close()

	body, err := ioutil.ReadAll(response.Body)
	if err != nil {
		log.Fatal("Failed to read response body:", err)
	}

	return string(body)
}

func quoteByTitle(title string) string {
	quotes, err := getQuotesByTitle(title)
	if err != nil {
		log.Fatal("Failed to get quote", err)
	}

	if quotes == nil {
		quotes := []AnimeQuote{}
		encodedTitle := url.QueryEscape(title)

		response, err := http.Get("https://animechan.xyz/api/quotes/anime?title=" + encodedTitle)
		if err != nil {
			log.Fatal("Request failed:", err)
		}
		defer response.Body.Close()

		body, err := ioutil.ReadAll(response.Body)
		fmt.Println(string(body))
		if err != nil {
			log.Fatal("Failed to read response body:", err)
		}

		err = json.Unmarshal(body, &quotes)
		if err != nil {
			log.Fatal("Failed to parse API response:", err)
		}

		for _, quote := range quotes {
			apiQuote := AnimeQuote{
				Name:      quote.Name,
				Character: quote.Character,
				Quote:     quote.Quote,
			}
			_, err = insertQuote(apiQuote)
			if err != nil {
				log.Fatal("Failed to insert quote into the database:", err)
			}
		}

		quotesJSON, err := json.Marshal(quotes)
		if err != nil {
			log.Fatal("Failed to serialize quotes to JSON:", err)
		}

		return string(quotesJSON)
	}

	quotesJSON, err := json.Marshal(quotes)
	if err != nil {
		log.Fatal("Failed to serialize quotes to JSON:", err)
	}

	return string(quotesJSON)
}
