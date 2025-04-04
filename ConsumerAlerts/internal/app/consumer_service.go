package app

import (
	"encoding/json"
	"log"
	"os"
	"sync"

	"github.com/yourusername/ConsumerAlerts/internal/domain"
	"github.com/yourusername/ConsumerAlerts/internal/infrastructure"
)

// ConsumerService handles the consumption of messages from RabbitMQ
type ConsumerService struct {
	rabbitMQ *infrastructure.RabbitMQConnection
}

// NewConsumerService creates a new consumer service
func NewConsumerService(rabbitMQ *infrastructure.RabbitMQConnection) *ConsumerService {
	return &ConsumerService{
		rabbitMQ: rabbitMQ,
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
		log.Printf("  🚦 Estado: %d", sensorMsg.Estado)
		
		// Add more visible alert for state 1
		if sensorMsg.Estado == 1 {
			log.Printf("  ⚠️ ALERT: Sensor %s detected an event! ⚠️", sensorMsg.Sensor)
		}
		log.Printf("=======================================\n")
	}
}