package main

import (
	"log"
	"pasetoAuth/server"
)

func main() {
	if _, err := server.NewServer("0.0.0.0:8080"); err != nil {
		log.Fatal("Failed to start server")
	}
}
