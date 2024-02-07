package server

import (
	"database/sql"
	"fmt"
	"log"
	"net/http"

	"github.com/davidtaing/go-webhook-server/internal/database"
)

// Abitrary typedef to represent an event
type WebhookEvent struct {
	ID    string `json:"id"`
	Event string `json:"event"`
}

type Env struct {
	db *sql.DB
}

func (env *Env) WebhookHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
	}

	// Handle the POST request here
	fmt.Println("Webhook received!")
}

func Start() {
	db, err := database.Open("./db/database.db")
	if err != nil {
		log.Fatal(err)
	}

	defer db.Close()

	env := Env{db: db}

	http.HandleFunc("/webhook", env.WebhookHandler)

	fmt.Printf("Starting server at port 8080\n")
	if err := http.ListenAndServe(":8080", nil); err != nil {
		log.Fatal(err)
	}
}
