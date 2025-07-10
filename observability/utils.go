package main

import (
	"context"
	"log"
	"runtime"

	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/metric"
)

func GetCPUMetrics() (metric.Registration, error) {
	meter := otel.Meter("")

	totalAllocGauge, _ := meter.Int64ObservableGauge(
		"memory_total_alloc",
		metric.WithDescription("Total allocated bytes"),
	)

	heapAllocGauge, _ := meter.Int64ObservableGauge(
		"memory_heap_alloc",
		metric.WithDescription("Heap allocated bytes"),
	)

	heapInuseGauge, _ := meter.Int64ObservableGauge(
		"memory_heap_inuse",
		metric.WithDescription("Heap in-use bytes"),
	)

	stackInuseGauge, _ := meter.Int64ObservableGauge(
		"memory_stack_inuse",
		metric.WithDescription("Stack in-use bytes"),
	)

	sysMemGauge, _ := meter.Int64ObservableGauge(
		"memory_sys",
		metric.WithDescription("Total system memory used"),
	)

	unregister, err := meter.RegisterCallback(func(ctx context.Context, observer metric.Observer) error {
		log.Println("getting metrics")
		var memStats runtime.MemStats
		runtime.ReadMemStats(&memStats)

		observer.ObserveInt64(totalAllocGauge, int64(memStats.TotalAlloc))
		observer.ObserveInt64(heapAllocGauge, int64(memStats.HeapAlloc))
		observer.ObserveInt64(heapInuseGauge, int64(memStats.HeapInuse))
		observer.ObserveInt64(stackInuseGauge, int64(memStats.StackInuse))
		observer.ObserveInt64(sysMemGauge, int64(memStats.Sys))

		return nil
	}, totalAllocGauge, heapAllocGauge, heapInuseGauge, stackInuseGauge, sysMemGauge)
	if err != nil {
		return nil, err
	}

	return unregister, nil
}
