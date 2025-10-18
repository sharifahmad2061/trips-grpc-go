package telemetry

import (
	"context"
	"log"

	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/exporters/otlp/otlplog/otlploggrpc"
	"go.opentelemetry.io/otel/exporters/otlp/otlpmetric/otlpmetricgrpc"
	"go.opentelemetry.io/otel/exporters/otlp/otlptrace/otlptracegrpc"
	"go.opentelemetry.io/otel/log/global"
	sdklog "go.opentelemetry.io/otel/sdk/log"
	"go.opentelemetry.io/otel/sdk/metric"
	"go.opentelemetry.io/otel/sdk/resource"
	"go.opentelemetry.io/otel/sdk/trace"
	semconv "go.opentelemetry.io/otel/semconv/v1.37.0"
)

func Init(ctx context.Context) (func(), error) {
	rsc, err := resource.Merge(
		resource.Default(),
		resource.NewWithAttributes(
			semconv.SchemaURL,
			semconv.ServiceNameKey.String("trip-grpc-go"),
		),
	)
	if err != nil {
		return nil, err
	}

	te, err := otlptracegrpc.New(ctx,
		otlptracegrpc.WithEndpoint("localhost:4317"),
		otlptracegrpc.WithInsecure(),
	)
	if err != nil {
		return nil, err
	}

	me, err := otlpmetricgrpc.New(ctx,
		otlpmetricgrpc.WithEndpoint("localhost:4317"),
		otlpmetricgrpc.WithInsecure(),
	)
	if err != nil {
		return nil, err
	}

	le, err := otlploggrpc.New(ctx,
		otlploggrpc.WithEndpoint("localhost:4317"),
		otlploggrpc.WithInsecure(),
	)

	if err != nil {
		return nil, err
	}

	tp := trace.NewTracerProvider(
		trace.WithResource(rsc),
		trace.WithBatcher(te),
	)
	otel.SetTracerProvider(tp)

	mp := metric.NewMeterProvider(
		metric.WithResource(rsc),
		metric.WithReader(metric.NewPeriodicReader(me)),
	)
	otel.SetMeterProvider(mp)

	lp := sdklog.NewLoggerProvider(
		sdklog.WithResource(rsc),
		sdklog.WithProcessor(
			sdklog.NewBatchProcessor(le),
		),
	)
	global.SetLoggerProvider(lp)

	shutdown := func() {
		if err := tp.Shutdown(ctx); err != nil {
			log.Printf("Error shutting down tracer: %v", err)
		}
		if err := mp.Shutdown(ctx); err != nil {
			log.Printf("Error shutting down meter: %v", err)
		}
		if err := lp.Shutdown(ctx); err != nil {
			log.Printf("Error shutting down logger: %v", err)
		}
	}

	return shutdown, nil
}
