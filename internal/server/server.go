package server

import (
	"database/sql"
	"fmt"
	"log"
	"net/http"
	"time"

	"github.com/davidtaing/go-webhook-server/internal/database"
	"github.com/davidtaing/go-webhook-server/internal/logger"
	"github.com/gorilla/mux"
)

// Abitrary typedef to represent an event
type WebhookEvent struct {
	ID    string `json:"id"`
	Event string `json:"event"`
}

type server struct {
	db     *sql.DB
	logger *logger.Logger
}

func (s *server) WebhookHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
	}

	fmt.Println("Webhook received!")
}

func LogMiddleware(logger *logger.Logger, next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		startTime := time.Now()

		next.ServeHTTP(w, r)

		duration := time.Since(startTime)

		logger.Infow("handled request",
			"method", r.Method,
			"url", r.URL.String(),
			"response_time", duration,
		)
	})
}

func Start() {
	logger := logger.New()

	db, err := database.Open("./db/database.db")
	if err != nil {
		log.Fatal(err)
	}

	defer db.Close()

	env := server{db: db, logger: logger}

	r := mux.NewRouter()

	r.HandleFunc("/webhook", env.WebhookHandler)
	r.Use(func(next http.Handler) http.Handler {
		return LogMiddleware(logger, next)
	})

	fmt.Printf("Starting server at port 8080\n")
	if err := http.ListenAndServe(":8080", r); err != nil {
		log.Fatal(err)
	}
}
