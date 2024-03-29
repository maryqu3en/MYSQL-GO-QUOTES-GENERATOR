package main

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"math/rand"
	"net/http"
	"strconv" //convert string to int

	_ "github.com/go-sql-driver/mysql"
)

type Quote struct {
	ID     int    `json:"id"`
	Book   string `json:"book"`
	Author string `json:"author"`
	Quote  string `json:"quote"`
}

var quotes []Quote

func loadQuotes() {
	db, err := sql.Open("mysql", "root:rootpassword@tcp(127.0.0.1:3306)/Quotes")
	if err != nil {
        panic(err)
    }
    defer db.Close()

	rows, err := db.Query("SELECT * FROM Quote")
	if err != nil {
        panic(err)
    }
    defer rows.Close()

	for rows.Next() {
		var quote Quote
		err = rows.Scan(&quote.ID, &quote.Book, &quote.Author, &quote.Quote)
		if err != nil {
			panic(err.Error())
		}
		quotes = append(quotes, quote)
	}
}

func enableCors(w *http.ResponseWriter) { //enable CORS for API access from different domains
	(*w).Header().Set("Access-Control-Allow-Origin", "*")
	(*w).Header().Set("Access-Control-Allow-Methods", "POST, GET, OPTIONS, PUT, DELETE")
	(*w).Header().Set("Access-Control-Allow-Headers", "Accept, Content-Type, Content-Length, Accept-Encoding, X-CSRF-Token, Authorization")
}

func getQuotes(w http.ResponseWriter, r *http.Request) {
	enableCors(&w)
	json.NewEncoder(w).Encode(quotes)
}

func getRandomQuote(w http.ResponseWriter, r *http.Request) {
	enableCors(&w)
	randomQuote := quotes[rand.Intn(len(quotes))]
	json.NewEncoder(w).Encode(randomQuote)
}
func getQuoteByID(w http.ResponseWriter, r *http.Request) {
	enableCors(&w)
	id, _ := strconv.Atoi(r.URL.Path[len("/api/quotes/"):])
	for _, quote := range quotes {
		if quote.ID == id {
			json.NewEncoder(w).Encode(quote)
			return
		}
	}
	http.Error(w, "Quote not found", http.StatusNotFound)
}

func main() {
	loadQuotes()

	http.HandleFunc("/api/quotes", getQuotes)
	http.HandleFunc("/api/quotes/random", getRandomQuote)
	http.HandleFunc("/api/quotes/", getQuoteByID)

	fmt.Println("Server running on http://localhost:3030")
	http.ListenAndServe(":3030", nil)  // nil = null value for pointers
}
