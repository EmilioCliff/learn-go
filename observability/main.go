package main

import (
	"context"
	"log"
	"os"
	"os/signal"
	"syscall"

	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

func main() {
	// setup signal handlers
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, os.Interrupt, syscall.SIGINT, syscall.SIGTERM)

	// establish a grpc connection for sending telementry data to otel collector - grafana alloy
	conn, err := grpc.NewClient("alloy:4317", grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		panic(err)
	}
	log.Println("Connection state: ", conn.GetState())

	// create an instance of openTelemetry and initialize the providers
	openTelemetry, err := NewOpenTelemetry()
	if err != nil {
		panic(err)
	}

	ctx := context.Background()
	shutdownTracer, err := openTelemetry.InitializeTracerProvider(ctx, conn)
	if err != nil {
		panic(err)
	}

	shutdownMeter, err := openTelemetry.InitializeMeterProvider(ctx, conn)
	if err != nil {
		panic(err)
	}

	shutdownLogger, err := openTelemetry.InitializeLoggerProvider(ctx, conn)
	if err != nil {
		panic(err)
	}

	// Register memory metrics (custom gauges) for periodic export
	unregisterCPUMetrics, err := GetCPUMetrics()
	if err != nil {
		panic(err)
	}

	server := NewServer()
	log.Println("Starting server on :3030")
	if err := server.Start(); err != nil {
		log.Fatalf("Failed to start server: %v", err)
	}

	<-quit

	log.Println("Received shutdown signal, stopping server...")
	if err := server.Stop(); err != nil {
		log.Fatalf("Failed to stop server: %v", err)
	}

	// Clean up and shutdown all the telemetry providers
	if err := shutdownTracer(context.Background()); err != nil {
		log.Fatal(err)
	}
	if err := shutdownMeter(context.Background()); err != nil {
		log.Fatal(err)
	}
	if err := shutdownLogger(context.Background()); err != nil {
		log.Fatal(err)
	}

	// Unregister custom metrics
	if err := unregisterCPUMetrics.Unregister(); err != nil {
		log.Fatal(err)
	}

	// Finally, close the grpc connection
	conn.Close()
}
