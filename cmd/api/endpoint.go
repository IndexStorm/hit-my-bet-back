package main

func (s *server) configureEndpoints() {
	api := s.app.Group("/v1")

	api.Post("/markets/create", s.createMarket)
	api.Post("/markets/init", s.initMarket)
}
