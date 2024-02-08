package server

import (
	"net/http"
	"time"

	"github.com/google/uuid"
)

type statusRecorder struct {
	http.ResponseWriter
	status int // this allows us to log the status code
}

func (rec *statusRecorder) WriteHeader(code int) {
	rec.status = code
	rec.ResponseWriter.WriteHeader(code)
}

func (s *server) middlewareLogging(h http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		startTime := time.Now()
		requestID := uuid.New().String()
		recorder := statusRecorder{w, 200}

		s.logger.Infow("request",
			"request_id", requestID,
			"method", r.Method,
			"url", r.URL.String(),
			"remote_addr", r.RemoteAddr,
			"user_agent", r.UserAgent(),
		)

		h.ServeHTTP(&recorder, r)

		duration := time.Since(startTime)

		s.logger.Infow("response",
			"request_id", requestID,
			"method", r.Method,
			"url", r.URL.String(),
			"response_time", duration,
			"status_code", recorder.status,
		)
	})
}
