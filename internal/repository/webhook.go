package repository

import (
	"database/sql"

	"github.com/davidtaing/go-webhook-server/internal/models"
)

type WebhookRepository struct {
	DB *sql.DB
}

// Adds a new webhook to the database.
// Returns nil if the webhook was added successfully, or the error upon failure.
func (w *WebhookRepository) Create(webhook models.Webhook) error {
	statement := `INSERT INTO webhooks (id, source, event, payload) VALUES ($1, $2, $3, $4)`

	_, err := w.DB.Exec(statement, webhook.ID, webhook.Source, webhook.Event, webhook.Payload)
	return err
}

// FindByID retrieves a webhook from the database by its ID.
// Returns nil if the webhook does not exist.
func (w *WebhookRepository) FindByID(id string) (*models.Webhook, error) {
	statement := `SELECT * FROM webhooks WHERE id = $1`

	var webhook models.Webhook

	row := w.DB.QueryRow(statement, id)
	switch err := row.Scan(&webhook.ID, &webhook.Source, &webhook.Event, &webhook.Created, &webhook.Updated, &webhook.Payload); err {
	case sql.ErrNoRows:
		return nil, nil
	case nil:
		return &webhook, nil
	default:
		return nil, err
	}
}

func (w *WebhookRepository) Get() ([]models.Webhook, error) {
	statement := `SELECT * FROM webhooks`

	rows, err := w.DB.Query(statement)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var webhooks []models.Webhook
	for rows.Next() {
		var webhook models.Webhook
		err := rows.Scan(&webhook.ID, &webhook.Source, &webhook.Event, &webhook.Created, &webhook.Updated, &webhook.Payload)
		if err != nil {
			return nil, err
		}
		webhooks = append(webhooks, webhook)
	}

	if err := rows.Err(); err != nil {
		return nil, err
	}

	return webhooks, nil
}
