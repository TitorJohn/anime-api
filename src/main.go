package main

import (
	"fmt"
	"log"
	"net/http"
)

func byAnime(w http.ResponseWriter, r *http.Request) {
	queryParams := r.URL.Query()
	search := queryParams.Get("s")

	fmt.Printf("You are searching for this value: %s", search)
	fmt.Fprint(w, quoteByTitle(search))
}

func random(w http.ResponseWriter, r *http.Request) {
	fmt.Fprint(w, randomQuote())
}

func main() {
	http.HandleFunc("/anime", byAnime)
	http.HandleFunc("/", random)

	fmt.Println("Server listening on port 3030")
	log.Fatal(http.ListenAndServe(":3030", nil))
}
