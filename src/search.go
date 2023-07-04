package main

import (
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
	encondedTitle := url.QueryEscape(title)
	response, err := http.Get("https://animechan.xyz/api/quotes/anime?title=" + encondedTitle)
	if err != nil {
		log.Fatal("Request failed:", err)
	}
	defer response.Body.Close()

	body, err := ioutil.ReadAll(response.Body)
	if err != nil {
		log.Fatal("Failed to read response body:", err)
	}

	return string(body)
}
