package server

import (
	"fmt"
	"log"
	"net/http"

	"github.com/davidtaing/go-webhook-server/internal/server/db"
)

func WebhookHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
	}

	// Handle the POST request here
	fmt.Println("Webhook received!")
}

func Start() {
	db, err := db.Open()
	if err != nil {
		log.Fatal(err)
	}

	defer db.Close()

	http.HandleFunc("/webhook", WebhookHandler)

	fmt.Printf("Starting server at port 8080\n")
	if err := http.ListenAndServe(":8080", nil); err != nil {
		log.Fatal(err)
	}
}
