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
	metricExporter, err := otlpmetricgrpc.New(ctx, otlpmetricgrpc.WithGRPCConn(conn))
	if err != nil {
		return nil, err
	}

	meterProvider := sdkmetric.NewMeterProvider(
		sdkmetric.WithReader(sdkmetric.NewPeriodicReader(metricExporter)),
		sdkmetric.WithResource(o.res),
	)
	otel.SetMeterProvider(meterProvider)

	return meterProvider.Shutdown, nil
}

func (o *openTelemetry) InitializeLoggerProvider(ctx context.Context, conn *grpc.ClientConn) (func(context.Context) error, error) {
	loggerExplorter, err := otlploggrpc.New(ctx, otlploggrpc.WithGRPCConn(conn))
	if err != nil {
		return nil, err
	}

	loggerExporterStdOut, err := stdoutlog.New()
	if err != nil {
		return nil, err
	}

	loggerProvider := log.NewLoggerProvider(
		log.WithProcessor(log.NewBatchProcessor(loggerExplorter)),
		log.WithProcessor(log.NewSimpleProcessor(loggerExporterStdOut)),
		log.WithResource(o.res),
	)
	global.SetLoggerProvider(loggerProvider)

	return loggerExplorter.Shutdown, nil
}

func (o *openTelemetry) InitializeTracerProvider(ctx context.Context, conn *grpc.ClientConn) (func(context.Context) error, error) {
	traceExporter, err := otlptracegrpc.New(ctx, otlptracegrpc.WithGRPCConn(conn))
	if err != nil {
		return nil, err
	}

	bsp := sdktrace.NewBatchSpanProcessor(traceExporter)
	tracerProvider := sdktrace.NewTracerProvider(
		sdktrace.WithSampler(sdktrace.AlwaysSample()),
		sdktrace.WithResource(o.res),
		sdktrace.WithSpanProcessor(bsp),
	)
	otel.SetTracerProvider(tracerProvider)

	otel.SetTextMapPropagator(propagation.TraceContext{})

	return tracerProvider.Shutdown, nil
}
