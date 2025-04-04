package app

import (
	"bytes"
	"encoding/json"
	"log"
	"net/http"
	"os"
	"sync"
	"time"

	"github.com/yourusername/ConsumerAlerts/internal/domain"
	"github.com/yourusername/ConsumerAlerts/internal/infrastructure"
)

// ConsumerService handles the consumption of messages from RabbitMQ
type ConsumerService struct {
	rabbitMQ *infrastructure.RabbitMQConnection
	client   *http.Client
}

// NewConsumerService creates a new consumer service
func NewConsumerService(rabbitMQ *infrastructure.RabbitMQConnection) *ConsumerService {
	return &ConsumerService{
		rabbitMQ: rabbitMQ,
		client:   &http.Client{Timeout: 10 * time.Second},
	}
}

// StartConsuming starts consuming messages from all queues
func (s *ConsumerService) StartConsuming() error {
	queues := []string{
		os.Getenv("RABBITMQ_QUEUE_KY026"),
		os.Getenv("RABBITMQ_QUEUE_MQ2"),
		os.Getenv("RABBITMQ_QUEUE_MQ135"),
		os.Getenv("RABBITMQ_QUEUE_DHT22"),
	}

	var wg sync.WaitGroup

	for _, queue := range queues {
		wg.Add(1)
		go func(queueName string) {
			defer wg.Done()
			s.consumeQueue(queueName)
		}(queue)
	}

	wg.Wait()
	return nil
}

// consumeQueue consumes messages from a specific queue
func (s *ConsumerService) consumeQueue(queueName string) {
	msgs, err := s.rabbitMQ.ConsumeQueue(queueName)
	if err != nil {
		log.Printf("Error consuming from queue %s: %v", queueName, err)
		return
	}

	log.Printf("Successfully connected to queue: %s - Waiting for messages...", queueName)
	
	// Create a ticker for heartbeat messages
	heartbeat := time.NewTicker(30 * time.Second)
	defer heartbeat.Stop()
	
	// Use a separate goroutine for the heartbeat
	go func() {
		for {
			select {
			case <-heartbeat.C:
				log.Printf("Still listening on queue: %s - No messages received yet", queueName)
			}
		}
	}()

	for msg := range msgs {
		log.Printf("\n==== NEW MESSAGE RECEIVED FROM %s ====", queueName)
		log.Printf("Raw message: %s", string(msg.Body))
		
		var sensorMsg domain.SensorMessage
		if err := json.Unmarshal(msg.Body, &sensorMsg); err != nil {
			log.Printf("Error unmarshaling message from queue %s: %v", queueName, err)
			continue
		}

		// Process the message with more visible formatting
		log.Printf("📊 PROCESSED MESSAGE DETAILS:")
		log.Printf("  📌 Queue: %s", queueName)
		log.Printf("  🔢 Número Serie: %s", sensorMsg.NumeroSerie)
		log.Printf("  🔍 Sensor: %s", sensorMsg.Sensor)
		log.Printf("  🕒 Fecha Activación: %s", sensorMsg.FechaActivacion)
		log.Printf("  🕓 Fecha Desactivación: %s", sensorMsg.FechaDesactivacion)
		log.Printf("  🚦 Estado: %s", sensorMsg.GetEstadoAsString())
		
		// Add special handling for DHT22 sensor
		if sensorMsg.IsDHT22() {
			log.Printf("  🌡️ Temperature reading: %s", sensorMsg.GetEstadoAsString())
		} else if sensorMsg.IsAlert() {
			log.Printf("  ⚠️ ALERT: Sensor %s detected an event! ⚠️", sensorMsg.Sensor)
		}
		
		// Forward the message to the API
		s.sendToAPI(sensorMsg)
		
		log.Printf("=======================================\n")
	}
}

// sendToAPI sends the sensor message to the API endpoint
func (s *ConsumerService) sendToAPI(msg domain.SensorMessage) {
	apiURL := "http://localhost:8080/api/alerts"
	
	// Create the payload
	payload, err := json.Marshal(msg)
	if err != nil {
		log.Printf("Error marshaling message for API: %v", err)
		return
	}
	
	// Send the request
	resp, err := s.client.Post(apiURL, "application/json", bytes.NewBuffer(payload))
	if err != nil {
		log.Printf("Error sending message to API: %v", err)
		return
	}
	defer resp.Body.Close()
	
	if resp.StatusCode >= 200 && resp.StatusCode < 300 {
		log.Printf("✅ Successfully sent message to API: %s", apiURL)
	} else {
		log.Printf("❌ Failed to send message to API. Status code: %d", resp.StatusCode)
	}
}