package telemetry

import (
	"context"
	"errors"
	"fmt"
	expmetergrpc "go.opentelemetry.io/otel/exporters/otlp/otlpmetric/otlpmetricgrpc"
	exptracegrpc "go.opentelemetry.io/otel/exporters/otlp/otlptrace/otlptracegrpc"
	metricsdk "go.opentelemetry.io/otel/sdk/metric"
	sdkrsc "go.opentelemetry.io/otel/sdk/resource"
	sdktrace "go.opentelemetry.io/otel/sdk/trace"
	semconv "go.opentelemetry.io/otel/semconv/v1.27.0"
	"time"
)

type Builder struct {
	tel *Telemetry
	res *sdkrsc.Resource
	err error
}

func NewBuilder() *Builder {
	return &Builder{tel: &Telemetry{}}
}

func (b *Builder) WithDefaultResource(ctx context.Context, service, namespace, version, instance string) *Builder {
	if b.err != nil {
		return b
	}
	var err error
	b.res, err = sdkrsc.New(
		ctx,
		sdkrsc.WithFromEnv(),
		sdkrsc.WithTelemetrySDK(),
		sdkrsc.WithProcess(),
		sdkrsc.WithOS(),
		sdkrsc.WithContainer(),
		sdkrsc.WithHost(),
		sdkrsc.WithAttributes(
			semconv.ServiceName(service),
			semconv.ServiceNamespace(namespace),
			semconv.ServiceVersion(version),
			semconv.ServiceInstanceID(instance),
		),
	)
	if err != nil {
		b.err = fmt.Errorf("create otel resource: %w", err)
	}
	return b
}

type TracerConfig struct {
	Endpoint string
	Insecure bool
}

func (b *Builder) WithTracer(ctx context.Context, config TracerConfig) *Builder {
	if b.err != nil {
		return b
	}
	if b.res == nil {
		b.err = errors.New("tracer requires a resource to be configured")
		return b
	}
	opt := []exptracegrpc.Option{
		exptracegrpc.WithEndpoint(config.Endpoint),
	}
	if config.Insecure {
		opt = append(opt, exptracegrpc.WithInsecure())
	}
	exp, err := exptracegrpc.New(ctx, opt...)
	if err != nil {
		b.err = fmt.Errorf("create otel grpc trace exporter: %w", err)
		return b
	}
	bsp := sdktrace.NewBatchSpanProcessor(exp)
	dropCheckProcessor := NewDropCheckSpanProcessor(bsp)
	b.tel.Tracer = sdktrace.NewTracerProvider(
		sdktrace.WithIDGenerator(&randomIDGenerator{}),
		sdktrace.WithResource(b.res),
		sdktrace.WithSampler(sdktrace.ParentBased(AttributeDropSampler(DropSpanAttributeName))),
		sdktrace.WithSpanProcessor(dropCheckProcessor),
	)
	return b
}

type MeterConfig struct {
	Endpoint string
	Insecure bool
	Interval time.Duration
}

func (b *Builder) WithMeter(ctx context.Context, config MeterConfig) *Builder {
	if b.err != nil {
		return b
	}
	if b.res == nil {
		b.err = errors.New("meter requires a resource to be configured")
		return b
	}
	opt := []expmetergrpc.Option{
		expmetergrpc.WithEndpoint(config.Endpoint),
	}
	if config.Insecure {
		opt = append(opt, expmetergrpc.WithInsecure())
	}
	exp, err := expmetergrpc.New(ctx, opt...)
	if err != nil {
		b.err = fmt.Errorf("create otel grpc metric exporter: %w", err)
		return b
	}
	periodicReader := metricsdk.NewPeriodicReader(exp, metricsdk.WithInterval(config.Interval))
	b.tel.Meter = metricsdk.NewMeterProvider(
		metricsdk.WithResource(b.res),
		metricsdk.WithReader(periodicReader),
	)
	return b
}

func (b *Builder) Build() (*Telemetry, error) {
	return b.tel, b.err
}
