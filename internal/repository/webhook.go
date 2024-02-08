package repository

import (
	"database/sql"

	"github.com/davidtaing/go-webhook-server/internal/models"
)

type WebhookRepository struct {
	db *sql.DB
}

// Adds a new webhook to the database.
// Returns nil if the webhook was added successfully, or the error upon failure.
func (w *WebhookRepository) Create(webhook models.Webhook) error {
	statement := `INSERT INTO webhooks (id, event, payload) VALUES ($1, $2, $3)`

	_, err := w.db.Exec(statement, webhook.ID, webhook.Event, webhook.Payload)
	return err
}

// FindByID retrieves a webhook from the database by its ID.
// Returns nil if the webhook does not exist.
func (w *WebhookRepository) FindByID(id string) (*models.Webhook, error) {
	statement := `SELECT * FROM webhooks WHERE id = $1`

	var webhook models.Webhook

	row := w.db.QueryRow(statement, id)
	switch err := row.Scan(&webhook.ID, &webhook.Event, &webhook.Payload); err {
	case sql.ErrNoRows:
		return nil, nil
	case nil:
		return &webhook, nil
	default:
		return nil, err
	}
}
