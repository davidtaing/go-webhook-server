package main

import (
	"fmt"
	"log"
	"net/http"
)

func WebhookHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
	}

	// Handle the POST request here
	fmt.Println("Webhook received!")
}

func main() {
	http.HandleFunc("/webhook", WebhookHandler)

	fmt.Printf("Starting server at port 8080\n")
	if err := http.ListenAndServe(":8080", nil); err != nil {
		log.Fatal(err)
	}
}
