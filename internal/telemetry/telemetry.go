package telemetry

import (
	"context"
	"errors"
	"fmt"
	metricsdk "go.opentelemetry.io/otel/sdk/metric"
	sdktrace "go.opentelemetry.io/otel/sdk/trace"
)

type Telemetry struct {
	Tracer *sdktrace.TracerProvider
	Meter  *metricsdk.MeterProvider
}

func (t *Telemetry) Close() error {
	var errs []error
	if tr := t.Tracer; tr != nil {
		if err := tr.Shutdown(context.Background()); err != nil {
			errs = append(errs, fmt.Errorf("shutdown tracer: %w", err))
		}
	}
	if mt := t.Meter; mt != nil {
		if err := mt.Shutdown(context.Background()); err != nil {
			errs = append(errs, fmt.Errorf("shutdown meter: %w", err))
		}
	}
	return errors.Join(errs...)
}
