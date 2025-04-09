package app

import (
   
    "encoding/json"
    "sync"
    "time"
    "log"
    "github.com/yourusername/ConsumerAlerts/internal/config"
    "github.com/yourusername/ConsumerAlerts/internal/domain"
    "github.com/yourusername/ConsumerAlerts/internal/infrastructure"
)

type ConsumerService struct {
    rabbitMQ    domain.MessageQueueConnection
    processor   *MessageProcessor
    apiClient   *infrastructure.APIClient
    config      *config.Config
    logger      Logger
}

func NewConsumerService(
    rabbitMQ domain.MessageQueueConnection,
    processor *MessageProcessor,
    apiClient *infrastructure.APIClient,
    cfg *config.Config,
    logger Logger,
) *ConsumerService {
    return &ConsumerService{
        rabbitMQ:    rabbitMQ,
        processor:   processor,
        apiClient:   apiClient,
        config:      cfg,
        logger:      logger,
    }
}

func (s *ConsumerService) StartConsuming() error {
    queues := []string{
        s.config.RabbitMQ.Queues.KY026,
        s.config.RabbitMQ.Queues.MQ2,
        s.config.RabbitMQ.Queues.MQ135,
        s.config.RabbitMQ.Queues.DHT22,
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


func (s *ConsumerService) consumeQueue(queueName string) {
	msgs, err := s.rabbitMQ.ConsumeQueue(queueName)
	if err != nil {
		log.Printf("Error consuming from queue %s: %v", queueName, err)
		return
	}

	log.Printf("Successfully connected to queue: %s - Waiting for messages...", queueName)
	
	
	heartbeat := time.NewTicker(30 * time.Second)
	defer heartbeat.Stop()
	
	// Heartbeat goroutine
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

		
		log.Printf("📊 PROCESSED MESSAGE DETAILS:")
		log.Printf("  📌 Queue: %s", queueName)
		log.Printf("  🔢 Número Serie: %s", sensorMsg.NumeroSerie)
		log.Printf("  🔍 Sensor: %s", sensorMsg.Sensor)
		log.Printf("  🕒 Fecha Activación: %s", sensorMsg.FechaActivacion)
		log.Printf("  🕓 Fecha Desactivación: %s", sensorMsg.FechaDesactivacion)
		log.Printf("  🚦 Estado: %s", sensorMsg.GetEstadoAsString())
		
		
		if sensorMsg.IsDHT22() {
			log.Printf("  🌡️ Temperature reading: %s", sensorMsg.GetEstadoAsString())
		} else if sensorMsg.IsAlert() {
			log.Printf("  ⚠️ ALERT: Sensor %s detected an event! ⚠️", sensorMsg.Sensor)
		}
		
		// Forward the message to the API
		s.sendToAPI(sensorMsg)
		
		
	}
}

// Update the sendToAPI method to use apiClient instead of client
func (s *ConsumerService) sendToAPI(msg domain.SensorMessage) {
    if err := s.apiClient.SendAlert(msg); err != nil {
        s.logger.Printf("Error sending message to API: %v", err)
    } else {
        s.logger.Printf("✅ Successfully sent message to API")
    }
}