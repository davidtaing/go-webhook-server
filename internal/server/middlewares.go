package server

import (
	"context"
	"net/http"
	"time"

	"github.com/google/uuid"
)

// Used in the middlewareLogging function to capture and log the status code.
type statusRecorder struct {
	http.ResponseWriter
	status int
}

func (rec *statusRecorder) WriteHeader(code int) {
	rec.status = code
	rec.ResponseWriter.WriteHeader(code)
}

// Generates a UUID for each request and adds it to the context.
func (s *server) middlewareRequestID(h http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		requestID := uuid.New().String()
		ctx := context.WithValue(r.Context(), "request_id", requestID)
		r = r.WithContext(ctx)

		h.ServeHTTP(w, r)
	})
}

// Logs the request and response details.
func (s *server) middlewareLogging(h http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		startTime := time.Now()
		recorder := statusRecorder{w, 200}

		s.logger.Infow("request",
			"request_id", r.Context().Value("request_id"),
			"method", r.Method,
			"url", r.URL.String(),
			"remote_addr", r.RemoteAddr,
			"user_agent", r.UserAgent(),
		)

		h.ServeHTTP(&recorder, r)

		duration := time.Since(startTime)

		s.logger.Infow("response",
			"request_id", r.Context().Value("request_id"),
			"method", r.Method,
			"url", r.URL.String(),
			"response_time", duration,
			"status_code", recorder.status,
		)
	})
}
