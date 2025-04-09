package main

import (
	"log"
	"os"
	"os/signal"
	"path/filepath"
	"syscall"

	"github.com/yourusername/ConsumerAlerts/internal/app"
	"github.com/yourusername/ConsumerAlerts/internal/config"
	"github.com/yourusername/ConsumerAlerts/internal/infrastructure"
)

// Create a standard logger that implements the app.Logger interface
type stdLogger struct{}

func (l *stdLogger) Printf(format string, v ...interface{})  { log.Printf(format, v...) }
func (l *stdLogger) Println(v ...interface{})               { log.Println(v...) }
func (l *stdLogger) Fatalf(format string, v ...interface{}) { log.Fatalf(format, v...) }

func main() {
	logger := &stdLogger{}

	// Cargar configuración
	cfg, err := config.LoadConfig(filepath.Join("..", ".env"))
	if err != nil {
		logger.Fatalf("Failed to load config: %v", err)
	}

	// Crear conexión RabbitMQ usando la configuración
	rabbitMQ, err := infrastructure.NewRabbitMQConnection(cfg.RabbitMQ)
	if err != nil {
		logger.Fatalf("Failed to connect to RabbitMQ: %v", err)
	}
	defer rabbitMQ.Close()

	// Create required dependencies
	messageProcessor := app.NewMessageProcessor(logger)
	apiClient := infrastructure.NewAPIClient(cfg.APIURL)

	// Crear servicio consumidor
	consumerService := app.NewConsumerService(
		rabbitMQ,
		messageProcessor,
		apiClient,
		cfg,
		logger,
	)

	// Manejar cierre graceful
	sigChan := make(chan os.Signal, 1)
	signal.Notify(sigChan, syscall.SIGINT, syscall.SIGTERM)

	// Iniciar consumo en una goroutine separada
	go func() {
		logger.Println("Starting to consume messages from queues...")
		if err := consumerService.StartConsuming(); err != nil {
			logger.Fatalf("Error consuming messages: %v", err)
		}
	}()

	// Esperar señal de terminación
	sig := <-sigChan
	logger.Printf("Received signal %v, shutting down...", sig)
}