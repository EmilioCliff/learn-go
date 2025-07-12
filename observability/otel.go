package main

import (
	"context"
	"fmt"

	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/exporters/otlp/otlplog/otlploggrpc"
	"go.opentelemetry.io/otel/exporters/otlp/otlpmetric/otlpmetricgrpc"
	"go.opentelemetry.io/otel/exporters/otlp/otlptrace/otlptracegrpc"
	"go.opentelemetry.io/otel/exporters/stdout/stdoutlog"
	"go.opentelemetry.io/otel/log/global"
	"go.opentelemetry.io/otel/propagation"
	"go.opentelemetry.io/otel/sdk/log"
	sdkmetric "go.opentelemetry.io/otel/sdk/metric"
	"go.opentelemetry.io/otel/sdk/resource"
	sdktrace "go.opentelemetry.io/otel/sdk/trace"
	semconv "go.opentelemetry.io/otel/semconv/v1.17.0"
	"google.golang.org/grpc"
)

type openTelemetry struct {
	res *resource.Resource
}

func NewOpenTelemetry() (*openTelemetry, error) {
	res, err := resource.New(context.Background(),
		resource.WithAttributes(
			semconv.ServiceName("Backend"),
			semconv.ContainerName("backend"),
			semconv.ServiceVersion("v0.0.1"),
		),
	)
	if err != nil {
		return nil, fmt.Errorf("failed to create resource: %w", err)
	}

	return &openTelemetry{
		res: res,
	}, nil
}

func (o *openTelemetry) InitializeMeterProvider(ctx context.Context, conn *grpc.ClientConn) (func(context.Context) error, error) {
	// Create a metric exporter using OTLP over gRPC
	metricExporter, err := otlpmetricgrpc.New(ctx, otlpmetricgrpc.WithGRPCConn(conn))
	if err != nil {
		return nil, err
	}

	// Create a MeterProvider with periodic reading and attach the resource
	meterProvider := sdkmetric.NewMeterProvider(
		sdkmetric.WithReader(sdkmetric.NewPeriodicReader(metricExporter)),
		sdkmetric.WithResource(o.res),
	)
	otel.SetMeterProvider(meterProvider)

	// Return a shutdown function to gracefully stop the provider
	return meterProvider.Shutdown, nil
}

func (o *openTelemetry) InitializeLoggerProvider(ctx context.Context, conn *grpc.ClientConn) (func(context.Context) error, error) {
	// Create a logger exporter using OTLP over gRPC
	loggerExplorter, err := otlploggrpc.New(ctx, otlploggrpc.WithGRPCConn(conn))
	if err != nil {
		return nil, err
	}

	// Create an optional stdout exporter for local debugging
	loggerExporterStdOut, err := stdoutlog.New()
	if err != nil {
		return nil, err
	}

	// Build a logger provider with both exporters: OTLP(batch exported) and stdout(simple/normal logger)
	loggerProvider := log.NewLoggerProvider(
		log.WithProcessor(log.NewBatchProcessor(loggerExplorter)),
		log.WithProcessor(log.NewSimpleProcessor(loggerExporterStdOut)),
		log.WithResource(o.res),
	)
	global.SetLoggerProvider(loggerProvider)

	// Return a shutdown function to gracefully stop the provider
	return loggerExplorter.Shutdown, nil
}

func (o *openTelemetry) InitializeTracerProvider(ctx context.Context, conn *grpc.ClientConn) (func(context.Context) error, error) {
	// Create a tracer exporter using OTLP over gRPC
	traceExporter, err := otlptracegrpc.New(ctx, otlptracegrpc.WithGRPCConn(conn))
	if err != nil {
		return nil, err
	}

	// Use a batch processor for performance in production settings
	bsp := sdktrace.NewBatchSpanProcessor(traceExporter)
	tracerProvider := sdktrace.NewTracerProvider(
		sdktrace.WithSampler(sdktrace.AlwaysSample()), // Always sample spans (good for testing/dev)
		sdktrace.WithResource(o.res),
		sdktrace.WithSpanProcessor(bsp),
	)
	otel.SetTracerProvider(tracerProvider)

	otel.SetTextMapPropagator(propagation.TraceContext{})

	// Return a shutdown function to gracefully stop the provider
	return tracerProvider.Shutdown, nil
}
