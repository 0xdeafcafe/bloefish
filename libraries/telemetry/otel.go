package telemetry

import (
	"context"
	"errors"
	"time"

	"github.com/0xdeafcafe/bloefish/libraries/clog"
	"github.com/0xdeafcafe/bloefish/libraries/contexts"
	"github.com/0xdeafcafe/bloefish/libraries/version"
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/exporters/otlp/otlplog/otlploghttp"
	"go.opentelemetry.io/otel/exporters/otlp/otlpmetric/otlpmetrichttp"
	"go.opentelemetry.io/otel/exporters/otlp/otlptrace/otlptracehttp"
	"go.opentelemetry.io/otel/log/global"
	"go.opentelemetry.io/otel/propagation"
	"go.opentelemetry.io/otel/sdk/log"
	"go.opentelemetry.io/otel/sdk/metric"
	"go.opentelemetry.io/otel/sdk/resource"
	"go.opentelemetry.io/otel/sdk/trace"
	semconv "go.opentelemetry.io/otel/semconv/v1.26.0"
)

type Config struct {
	Enable      bool   `env:"ENABLED"`
	EndpointURL string `env:"ENDPOINT_URL"`
}

func (c Config) MustSetup(ctx context.Context) (shutdown func(context.Context) error) {
	shutdown, err := c.setup(ctx)
	if err != nil {
		panic(err)
	}

	return shutdown
}

func (c Config) setup(ctx context.Context) (shutdown func(context.Context) error, err error) {
	if !c.Enable {
		return func(context.Context) error { return nil }, nil
	}
	if c.EndpointURL == "" {
		c.EndpointURL = "http://util_telemetry:4318"
	}

	var shutdownFuncs []func(context.Context) error

	// shutdown calls cleanup functions registered via shutdownFuncs.
	// The errors from the calls are joined.
	// Each registered cleanup will be invoked once.
	shutdown = func(ctx context.Context) error {
		var err error

		for _, fn := range shutdownFuncs {
			err = errors.Join(err, fn(ctx))
		}

		shutdownFuncs = nil

		return err
	}

	// handleErr calls shutdown for cleanup and makes sure that all errors are returned.
	handleErr := func(inErr error) {
		err = errors.Join(inErr, shutdown(ctx))
	}

	// Set up propagator.
	prop := c.newPropagator()
	otel.SetTextMapPropagator(prop)

	resource, err := c.newResource(ctx)
	if err != nil {
		handleErr(err)
		return
	}

	// Set up trace provider.
	tracerProvider, err := c.newTracerProvider(ctx, resource)
	if err != nil {
		handleErr(err)
		return
	}
	shutdownFuncs = append(shutdownFuncs, tracerProvider.Shutdown)
	otel.SetTracerProvider(tracerProvider)

	// Set up meter provider.
	meterProvider, err := c.newMeterProvider(ctx, resource)
	if err != nil {
		handleErr(err)
		return
	}
	shutdownFuncs = append(shutdownFuncs, meterProvider.Shutdown)
	otel.SetMeterProvider(meterProvider)

	// Set up logger provider.
	loggerProvider, err := c.newLoggerProvider(ctx, resource)
	if err != nil {
		handleErr(err)
		return
	}
	shutdownFuncs = append(shutdownFuncs, loggerProvider.Shutdown)
	global.SetLoggerProvider(loggerProvider)

	return
}

func (c Config) newPropagator() propagation.TextMapPropagator {
	return propagation.NewCompositeTextMapPropagator(
		propagation.TraceContext{},
		propagation.Baggage{},
	)
}

func (c Config) newTracerProvider(ctx context.Context, res *resource.Resource) (*trace.TracerProvider, error) {
	traceExporter, err := otlptracehttp.New(ctx, otlptracehttp.WithEndpointURL(c.EndpointURL))
	if err != nil {
		return nil, err
	}

	return trace.NewTracerProvider(
		trace.WithBatcher(traceExporter, trace.WithBatchTimeout(time.Second)),
		trace.WithResource(res),
	), nil
}

func (c Config) newMeterProvider(ctx context.Context, res *resource.Resource) (*metric.MeterProvider, error) {
	metricExporter, err := otlpmetrichttp.New(ctx, otlpmetrichttp.WithEndpointURL(c.EndpointURL))
	if err != nil {
		return nil, err
	}

	return metric.NewMeterProvider(
		metric.WithReader(metric.NewPeriodicReader(metricExporter, metric.WithInterval(time.Minute))),
		metric.WithResource(res),
	), nil
}

func (c Config) newLoggerProvider(ctx context.Context, res *resource.Resource) (*log.LoggerProvider, error) {
	logExporter, err := otlploghttp.New(ctx, otlploghttp.WithEndpointURL(c.EndpointURL))
	if err != nil {
		return nil, err
	}

	return log.NewLoggerProvider(
		log.WithProcessor(log.NewBatchProcessor(logExporter)),
		log.WithResource(res),
	), nil
}

func (c Config) newResource(ctx context.Context) (*resource.Resource, error) {
	svcInfo := contexts.GetServiceInfo(ctx)

	res, err := resource.New(
		context.Background(),
		resource.WithTelemetrySDK(),
		resource.WithProcess(),
		resource.WithOS(),
		resource.WithContainer(),
		resource.WithHost(),
		resource.WithAttributes(
			semconv.ServiceNameKey.String(svcInfo.Service),
			semconv.ServiceVersionKey.String(version.Revision),
		),
	)

	if errors.Is(err, resource.ErrPartialResource) || errors.Is(err, resource.ErrSchemaURLConflict) {
		clog.Get(ctx).WithError(err).Warn("failed to create telemetry resource")
	} else if err != nil {
		return nil, err
	}

	return res, nil
}
