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

type server struct {
	db     *sql.DB
	router *mux.Router
	logger *logger.Logger
}

func setup(s *server) *server {
	s.routes()
	return s
}

func Run() {
	l := logger.New()
	s := setup(&server{
		logger: l,
		db:     database.Open("./db/database.db", l),
		router: mux.NewRouter(),
	})

	defer s.db.Close()

	fmt.Printf("Starting server at port 8080\n")
	if err := http.ListenAndServe(":8080", s.router); err != nil {
		log.Fatal(err)
	}
}
