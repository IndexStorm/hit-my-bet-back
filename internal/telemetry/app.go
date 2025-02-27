package telemetry

import (
	"context"
	"time"
)

type AppConfig struct {
	Service, Namespace, Version, Instance string

	TraceEndpoint string
	TraceInsecure bool

	MeterEndpoint string
	MeterInsecure bool
	MeterInterval time.Duration
}

func NewApp(ctx context.Context, config AppConfig) (*Telemetry, error) {
	return NewBuilder().
		WithDefaultResource(ctx, config.Service, config.Namespace, config.Version, config.Instance).
		WithTracer(ctx, TracerConfig{
			Endpoint: config.TraceEndpoint,
			Insecure: config.TraceInsecure,
		}).
		WithMeter(ctx, MeterConfig{
			Endpoint: config.MeterEndpoint,
			Insecure: config.MeterInsecure,
			Interval: config.MeterInterval,
		}).
		Build()
}
