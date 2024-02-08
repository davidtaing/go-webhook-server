package server

import (
	"net/http"
	"time"
)

func (s *server) middlewareLogging(h http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		startTime := time.Now()

		h.ServeHTTP(w, r)

		duration := time.Since(startTime)

		s.logger.Infow("handled request",
			"method", r.Method,
			"url", r.URL.String(),
			"response_time", duration,
		)
	})
}
