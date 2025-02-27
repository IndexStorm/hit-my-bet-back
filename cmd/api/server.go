package main

import (
	"github.com/IndexStorm/hit-my-bet-back/internal/repository/prediction"
	"github.com/goccy/go-json"
	"github.com/gofiber/fiber/v2"
	"github.com/rs/zerolog"
	"go.opentelemetry.io/otel/trace"
	"time"
)

type server struct {
	app            *fiber.App
	logger         zerolog.Logger
	tracer         trace.Tracer
	predictionRepo prediction.Repository
}

func newServer(logger zerolog.Logger, tr trace.Tracer, predictionRepo prediction.Repository) *server {
	app := fiber.New(fiber.Config{
		DisableStartupMessage: true,
		ReadTimeout:           time.Second * 15,
		WriteTimeout:          time.Second * 15,
		IdleTimeout:           time.Second * 30,
		JSONEncoder:           json.Marshal,
		JSONDecoder:           json.Unmarshal,
	})
	return &server{
		app:            app,
		logger:         logger,
		tracer:         tr,
		predictionRepo: predictionRepo,
	}
}

func (s *server) start(address string) error {
	s.configureEndpoints()
	return s.app.Listen(address)
}

func (s *server) Close() error {
	return s.app.Shutdown()
}
