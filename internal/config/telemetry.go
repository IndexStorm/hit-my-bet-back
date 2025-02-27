package config

type Telemetry struct {
	OpenTelemetryEndpoint string `env:"OTEL_ENDPOINT,notEmpty,unset"`
}
