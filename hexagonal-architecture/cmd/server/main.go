package main

import (
	"database/sql"
	"log"

	"myapp/internal/handlers"
	"myapp/internal/postgres"
	"myapp/internal/rabbitmq"
	"myapp/internal/repository"
	"myapp/internal/services"

	"github.com/gin-gonic/gin"
	"github.com/streadway/amqp"
)

func main() {
	var paymentRepo repository.PaymentRepository

	// Initialize database connections

	db, err := sql.Open("postgres", "postgres://user:password@localhost/dbname?sslmode=disable") // for PostgreSQL
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()

	paymentRepo = postgres.NewPostgresPaymentRepository(db)

	// Uncomment this section to use MySQL
	// db, err := sql.Open("mysql", "user:password@tcp(127.0.0.1:3306)/dbname") // for MySQL
	// if err != nil {
	// 	log.Fatal(err)
	// }
	// defer db.Close()

	// paymentRepo = mysql.NewMySQLPaymentRepository(db)

	// Initialize RabbitMQ connections
	conn, err := amqp.Dial("amqp://user:password@localhost:5672/")
	if err != nil {
		log.Fatal(err)
	}
	defer conn.Close()

	// Set up repository and message publisher
	messagePublisher := rabbitmq.NewMessagePublisher(conn)

	// Create service with dependencies
	paymentService := services.NewPaymentService(paymentRepo, messagePublisher)

	// Set up handlers and start server
	router := gin.Default()
	handlers.RegisterRoutes(router, paymentService)
	router.Run(":8080")
}
