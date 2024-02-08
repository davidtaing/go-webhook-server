package server

import (
	"bytes"
	"encoding/json"
	"io"
	"net/http"

	"github.com/davidtaing/go-webhook-server/internal/models"
	"github.com/davidtaing/go-webhook-server/internal/repository"
)

func (s *server) handleWebhook() http.HandlerFunc {
	type request struct {
		ID    string `json:"id"`
		Event string `json:"event"`
	}

	return func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPost {
			http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		}

		var (
			body       request      // Unmarshalled request body / payload.
			jsonBuffer bytes.Buffer // Raw JSON payload to be saved in db.
		)

		tee := io.TeeReader(r.Body, &jsonBuffer)
		err := json.NewDecoder(tee).Decode(&body)
		if err != nil {
			s.logger.Error(err)
			http.Error(w, "Bad request", http.StatusBadRequest)
			return
		}

		// check if webhook exists
		webhookRepo := &repository.WebhookRepository{DB: s.db}
		exists, err := webhookRepo.FindByID(body.ID)
		if err != nil {
			s.logger.Error(err)
			http.Error(w, "Internal server error", http.StatusInternalServerError)
			return
		}

		if exists != nil {
			// non-200 statuses are treated as errors,
			// leading to webhooks being retried / resent by the sender
			w.WriteHeader(http.StatusOK)
			return
		}

		// create if webhook doesn't exist
		webhook := models.Webhook{
			ID:      body.ID,
			Event:   body.Event,
			Payload: jsonBuffer.String(),
			Source:  "development",
		}

		err = webhookRepo.Create(webhook)
		if err != nil {
			s.logger.Errorw("Failed to create webhook",
				"error", err,
			)
			http.Error(w, "Internal server error", http.StatusInternalServerError)
			return
		}

		w.WriteHeader(http.StatusOK)
	}
}
