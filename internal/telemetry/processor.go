package telemetry

import (
	"context"
	sdktrace "go.opentelemetry.io/otel/sdk/trace"
	"go.opentelemetry.io/otel/trace"
	"sync/atomic"
)

type dropCheckSpanProcessor struct {
	queue    chan sdktrace.ReadOnlySpan
	next     sdktrace.SpanProcessor
	shutdown atomic.Bool
	drained  chan struct{}
	m        map[[24]byte][]sdktrace.ReadOnlySpan
}

// NewDropCheckSpanProcessor creates a new processor that checks parent spans for drop sample attribute
func NewDropCheckSpanProcessor(next sdktrace.SpanProcessor) sdktrace.SpanProcessor {
	sp := &dropCheckSpanProcessor{
		queue:   make(chan sdktrace.ReadOnlySpan, sdktrace.DefaultMaxQueueSize),
		next:    next,
		drained: make(chan struct{}, 1),
		m:       make(map[[24]byte][]sdktrace.ReadOnlySpan, sdktrace.DefaultMaxQueueSize),
	}
	go sp.processQueue()
	return sp
}

func (p *dropCheckSpanProcessor) OnStart(parent context.Context, s sdktrace.ReadWriteSpan) {
}

func (p *dropCheckSpanProcessor) OnEnd(s sdktrace.ReadOnlySpan) {
	if p.shutdown.Load() {
		return
	}
	select {
	case p.queue <- s:
	default:
	}
	//p.next.OnEnd(s)
}

func (p *dropCheckSpanProcessor) Shutdown(ctx context.Context) error {
	p.shutdown.Store(true)
	close(p.queue)
	<-p.drained
	return p.next.Shutdown(ctx)
}

func (p *dropCheckSpanProcessor) ForceFlush(ctx context.Context) error {
	return p.next.ForceFlush(ctx)
}

func (p *dropCheckSpanProcessor) processQueue() {
	for span := range p.queue {
		p.process(span)
	}
	p.drained <- struct{}{}
}

func getSpanID(trid trace.TraceID, spid trace.SpanID) [24]byte {
	var data [24]byte
	copy(data[:16], trid[:])
	copy(data[16:], spid[:])
	return data
}

func (p *dropCheckSpanProcessor) process(s sdktrace.ReadOnlySpan) {
	for _, attr := range s.Attributes() {
		if string(attr.Key) == DropSpanAttributeName && attr.Value.AsBool() {
			ctx := s.SpanContext()
			spanID := getSpanID(ctx.TraceID(), ctx.SpanID())
			if children, ok := p.m[spanID]; ok {
				for _, child := range children {
					childCtx := child.SpanContext()
					delete(p.m, getSpanID(childCtx.TraceID(), childCtx.SpanID()))
				}
			}
			delete(p.m, spanID)
			return
		}
	}
	parent := s.Parent()
	if !parent.IsValid() {
		ctx := s.SpanContext()
		spanID := getSpanID(ctx.TraceID(), ctx.SpanID())
		if children, ok := p.m[spanID]; ok {
			for _, child := range children {
				childCtx := child.SpanContext()
				delete(p.m, getSpanID(childCtx.TraceID(), childCtx.SpanID()))
				p.next.OnEnd(child)
			}
		}
		p.next.OnEnd(s)
		delete(p.m, spanID)
		return
	}
	ctx := s.Parent()
	parentID := getSpanID(ctx.TraceID(), ctx.SpanID())
	children, ok := p.m[parentID]
	if !ok {
		children = make([]sdktrace.ReadOnlySpan, 0, 2)
	}
	children = append(children, s)
	p.m[parentID] = children
}
