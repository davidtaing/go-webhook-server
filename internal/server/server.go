package server

import (
	"database/sql"
	"fmt"
	"log"
	"net/http"

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
	router *mux.Router
	logger *logger.Logger
}

func newServer() *server {
	s := &server{}
	s.routes()
	return s
}

func Run() {
	s := newServer()
	s.db = database.Open("./db/database.db")
	s.logger = logger.New()
	s.router = mux.NewRouter()

	defer s.db.Close()

	fmt.Printf("Starting server at port 8080\n")
	if err := http.ListenAndServe(":8080", s.router); err != nil {
		log.Fatal(err)
	}
}
