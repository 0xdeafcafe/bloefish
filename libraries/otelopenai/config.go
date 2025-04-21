package otelopenai

import (
	"go.opentelemetry.io/otel/propagation"
	oteltrace "go.opentelemetry.io/otel/trace"
)

// config is used to configure the middleware.
type config struct {
	tracerProvider                oteltrace.TracerProvider
	propagators                   propagation.TextMapPropagator
	traceIDResponseHeaderKey      string
	traceSampledResponseHeaderKey string
}

// Option specifies instrumentation configuration options.
type Option interface {
	apply(*config)
}

type optionFunc func(*config)

func (o optionFunc) apply(c *config) {
	o(c)
}

// WithTracerProvider specifies a tracer provider to use for creating a tracer.
// If none is specified, the global provider is used.
func WithTracerProvider(provider oteltrace.TracerProvider) Option {
	return optionFunc(func(c *config) {
		c.tracerProvider = provider
	})
}

// WithPropagators specifies propagators to use for extracting
// information from the HTTP requests. If none are specified, global
// ones will be used.
func WithPropagators(propagators propagation.TextMapPropagator) Option {
	return optionFunc(func(c *config) {
		c.propagators = propagators
	})
}

// WithTraceIDResponseHeader enables adding a header with the trace ID to the response.
// If not set, no trace ID header will be added to the response.
func WithTraceIDResponseHeader(key string) Option {
	return optionFunc(func(c *config) {
		c.traceIDResponseHeaderKey = key
	})
}

// WithTraceSampledResponseHeader enables adding a header with the trace sampling decision to the response.
// If not set, no trace sampling header will be added to the response.
func WithTraceSampledResponseHeader(key string) Option {
	return optionFunc(func(c *config) {
		c.traceSampledResponseHeaderKey = key
	})
}
