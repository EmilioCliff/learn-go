package main

import (
	"context"
	"errors"
	"log"
	"log/slog"
	"net"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"go.opentelemetry.io/contrib/bridges/otelslog"
	"go.opentelemetry.io/contrib/instrumentation/github.com/gin-gonic/gin/otelgin"
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/metric"
	"go.opentelemetry.io/otel/trace"
)

type Server struct {
	router *gin.Engine  // HTTP router using Gin
	ln     net.Listener // TCP listener
	srv    *http.Server // Underlying HTTP server

	meter  metric.Meter // OpenTelemetry Meter for metrics
	logger *slog.Logger // Structured logger (bridged with OpenTelemetry)
	tracer trace.Tracer // Tracer for distributed tracing

	startTime time.Time // Used to calculate uptime in health check
}

func NewServer() *Server {
	// Initialize observability tools
	meter := otel.Meter("")          // Get global Meter
	logger := otelslog.NewLogger("") // Structured logger via Otel-Slog bridge
	tracer := otel.Tracer("Backend") // Tracer named "Backend"

	// Create server instance
	s := &Server{
		router: gin.Default(),
		meter:  meter,
		logger: logger,
		tracer: tracer,
	}

	// Apply OpenTelemetry middleware for Gin (tracing)
	s.router.Use(otelgin.Middleware("Backend"))

	// Register routes
	s.router.GET("/health", s.healthCheckHandler)  // Simple health check
	s.router.GET("/compute", s.computeOneHandler)  // Simulated CPU-heavy GET route
	s.router.POST("/compute", s.computeTwoHandler) // POST route with input parsing and delay

	// Setup the HTTP server config
	s.srv = &http.Server{
		Addr:         ":3030", // Server will listen on port 3030
		Handler:      s.router.Handler(),
		ReadTimeout:  10 * time.Second,
		WriteTimeout: 10 * time.Second,
	}

	return s
}

func (s *Server) Start() error {
	var err error
	if s.ln, err = net.Listen("tcp", s.srv.Addr); err != nil {
		return err
	}

	s.startTime = time.Now()

	go func(s *Server) {
		err := s.srv.Serve(s.ln)
		if err != nil && err != http.ErrServerClosed {
			panic(err)
		}
	}(s)

	return nil
}

func (s *Server) Stop() error {
	log.Println("Shutting down http server...")

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	return s.srv.Shutdown(ctx)
}

func (s *Server) healthCheckHandler(c *gin.Context) {
	// Start a new span for the health check
	ctx, span := s.tracer.Start(c.Request.Context(), "HealthCheck")
	defer span.End()

	// Log that the health check was triggered
	s.logger.InfoContext(ctx, "Health check requested")

	c.JSON(http.StatusOK, gin.H{
		"status":  "ok",
		"uptime":  time.Since(s.startTime).String(),
		"message": "Service is healthy",
	})
}

func (s *Server) computeOneHandler(c *gin.Context) {
	// Start a span for tracing the compute operation
	ctx, span := s.tracer.Start(c.Request.Context(), "ComputeOne")
	defer span.End()

	start := time.Now()
	s.logger.InfoContext(ctx, "ComputeOneHandler started")

	// Simulate a CPU-heavy calculation
	result := 0
	for i := 0; i < 10_000_000; i++ {
		result += i % 3
	}

	duration := time.Since(start)

	s.logger.InfoContext(ctx, "ComputeOneHandler completed", "duration", duration, "result", result)

	// Add attributes to the span
	span.SetAttributes(
		attribute.Int("compute.iterations", 10_000_000),
		attribute.String("result.type", "int"),
	)

	c.JSON(http.StatusOK, gin.H{
		"result":   result,
		"duration": duration.String(),
	})
}

type ComputeInput struct {
	Numbers []int `json:"numbers"`
}

func (s *Server) computeTwoHandler(c *gin.Context) {
	// Start a span for tracing the compute operation
	ctx, span := s.tracer.Start(c.Request.Context(), "ComputeTwo")
	defer span.End()

	var input ComputeInput
	if err := c.ShouldBindJSON(&input); err != nil {
		s.logger.ErrorContext(ctx, "Invalid JSON input", "error", err)
		span.RecordError(err)
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid input"})
		return
	}

	if len(input.Numbers) == 0 {
		err := errors.New("no numbers provided")
		s.logger.WarnContext(ctx, "Empty input received", "error", err)
		span.RecordError(err)
		c.JSON(http.StatusUnprocessableEntity, gin.H{"error": err.Error()})
		return
	}

	// Log and tag trace with size of input
	s.logger.InfoContext(ctx, "Processing computeTwo request", "input_length", len(input.Numbers))
	span.SetAttributes(attribute.Int("input.count", len(input.Numbers)))

	// Perform the sum
	sum := 0
	for _, num := range input.Numbers {
		sum += num
	}

	s.logger.InfoContext(ctx, "computeTwo result", "sum", sum)

	c.JSON(http.StatusOK, gin.H{
		"sum":     sum,
		"message": "Computation completed successfully",
	})
}
