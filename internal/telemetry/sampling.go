package telemetry

import (
	"go.opentelemetry.io/otel/attribute"
	sdktrace "go.opentelemetry.io/otel/sdk/trace"
	oteltrace "go.opentelemetry.io/otel/trace"
)

const DropSpanAttributeName = "span.drop"

func DropSpan(v bool) attribute.KeyValue {
	return attribute.Bool(DropSpanAttributeName, v)
}

type attributeDropSampler struct {
	attrName string
}

func AttributeDropSampler(attrName string) sdktrace.Sampler {
	return &attributeDropSampler{attrName: attrName}
}

func (s *attributeDropSampler) ShouldSample(p sdktrace.SamplingParameters) sdktrace.SamplingResult {
	decision := sdktrace.RecordAndSample
	for _, attr := range p.Attributes {
		if string(attr.Key) == s.attrName {
			if attr.Value.AsBool() {
				decision = sdktrace.Drop
				break
			}
		}
	}
	return sdktrace.SamplingResult{
		Decision:   decision,
		Attributes: nil,
		Tracestate: oteltrace.SpanContextFromContext(p.ParentContext).TraceState(),
	}
}

func (s *attributeDropSampler) Description() string {
	return "AttributeDropSampler{" + s.attrName + "}"
}
