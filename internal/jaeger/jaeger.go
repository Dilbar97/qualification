package jaeger

import (
	"errors"

	"go.opentelemetry.io/otel"
	exJ "go.opentelemetry.io/otel/exporters/jaeger"
	"go.opentelemetry.io/otel/sdk/resource"
	tracesdk "go.opentelemetry.io/otel/sdk/trace"
	semconv "go.opentelemetry.io/otel/semconv/v1.18.0"
)

func InitTracer(url, serviceName string) (*tracesdk.TracerProvider, error) {
	tp, err := tracerProvider(url, serviceName)
	if err != nil {
		return nil, err
	}

	otel.SetTracerProvider(tp)

	return tp, nil
}

func tracerProvider(dsn, serviceName string) (*tracesdk.TracerProvider, error) {
	if len(serviceName) == 0 {
		return nil, errors.New("empty opentelemetry service name")
	}

	// Create the Jaeger exporter
	exp, err := exJ.New(exJ.WithCollectorEndpoint(exJ.WithEndpoint(dsn)))
	if err != nil {
		return nil, err
	}

	tp := tracesdk.NewTracerProvider(
		// Always be sure to batch in production.
		tracesdk.WithBatcher(exp),
		// Record information about this application in a Resource.
		tracesdk.WithResource(resource.NewWithAttributes(
			semconv.SchemaURL,
			semconv.ServiceNameKey.String(serviceName),
		)),
	)
	return tp, nil
}
