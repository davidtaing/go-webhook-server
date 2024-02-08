package server

func (s *server) routes() {
	s.router.HandleFunc("/webhook", s.handleWebhook())
	s.router.Use(s.middlewareRequestID)
	s.router.Use(s.middlewareLogging)
}
