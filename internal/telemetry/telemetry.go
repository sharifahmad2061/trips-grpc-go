package telemetry

import (
	"context"
	"log"

	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/exporters/otlp/otlpmetric/otlpmetricgrpc"
	"go.opentelemetry.io/otel/exporters/otlp/otlptrace/otlptracegrpc"
	"go.opentelemetry.io/otel/sdk/metric"
	"go.opentelemetry.io/otel/sdk/resource"
	"go.opentelemetry.io/otel/sdk/trace"
	semconv "go.opentelemetry.io/otel/semconv/v1.17.0"
)

func Init(ctx context.Context) (func(), error) {
	rsc, err := resource.New(ctx,
		resource.WithAttributes(
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

	shutdown := func() {
		if err := tp.Shutdown(ctx); err != nil {
			log.Printf("Error shutting down tracer: %v", err)
		}
		if err := mp.Shutdown(ctx); err != nil {
			log.Printf("Error shutting down meter: %v", err)
		}
	}

	return shutdown, nil
}
