package main

import (
	"flag"
	"os"
	"os/signal"
	"strconv"
	"syscall"
)

func main() {
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, os.Interrupt, syscall.SIGINT, syscall.SIGTERM)

	port := flag.Uint("port", 5000, "TCP Port Number for Blockchain Server")
	flag.Parse()

	portStr := os.Getenv("PORT")
	if os.Getenv("ENVIRONMENT") == "production" && portStr != "" {
		p, err := strconv.ParseUint(portStr, 10, 16)
		if err != nil {
			panic(err)
		}
		portVal := uint(p)
		port = &portVal
	}

	bcs := NewBlockchainServer(uint16(*port))

	if err := bcs.Start(); err != nil {
		panic(err)
	}

	<-quit

	if err := bcs.Stop(); err != nil {
		panic(err)
	}
}
