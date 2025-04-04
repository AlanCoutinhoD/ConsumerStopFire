package main

import (
	"log"
	"os"
	"os/signal"
	"path/filepath"
	"syscall"

	"github.com/joho/godotenv"
	"github.com/yourusername/ConsumerAlerts/internal/app"
	"github.com/yourusername/ConsumerAlerts/internal/infrastructure"
)

func main() {
	
	envPath := filepath.Join("..", ".env")
	if err := godotenv.Load(envPath); err != nil {
		log.Fatalf("Error loading .env file: %v", err)
	}

	// Create RabbitMQ connection
	rabbitMQ, err := infrastructure.NewRabbitMQConnection()
	if err != nil {
		log.Fatalf("Failed to connect to RabbitMQ: %v", err)
	}
	defer rabbitMQ.Close()

	// Create consumer service
	consumerService := app.NewConsumerService(rabbitMQ)

	// Handle graceful shutdown
	sigChan := make(chan os.Signal, 1)
	signal.Notify(sigChan, syscall.SIGINT, syscall.SIGTERM)

	// Start consuming in a separate goroutine
	go func() {
		log.Println("Starting to consume messages from queues...")
		if err := consumerService.StartConsuming(); err != nil {
			log.Fatalf("Error consuming messages: %v", err)
		}
	}()

	// Wait for termination signal
	sig := <-sigChan
	log.Printf("Received signal %v, shutting down...", sig)
}