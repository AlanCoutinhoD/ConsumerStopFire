package infrastructure

import (
	"fmt"
	"log"
	"os"

	amqp "github.com/rabbitmq/amqp091-go"
)

// RabbitMQConnection handles the connection to RabbitMQ
type RabbitMQConnection struct {
	Connection *amqp.Connection
	Channel    *amqp.Channel
}

// NewRabbitMQConnection creates a new RabbitMQ connection
func NewRabbitMQConnection() (*RabbitMQConnection, error) {
	host := os.Getenv("RABBITMQ_HOST")
	port := os.Getenv("RABBITMQ_PORT")
	user := os.Getenv("RABBITMQ_USER")
	password := os.Getenv("RABBITMQ_PASSWORD")

	connectionString := fmt.Sprintf("amqp://%s:%s@%s:%s/", user, password, host, port)
	
	conn, err := amqp.Dial(connectionString)
	if err != nil {
		return nil, fmt.Errorf("failed to connect to RabbitMQ: %w", err)
	}

	ch, err := conn.Channel()
	if err != nil {
		conn.Close()
		return nil, fmt.Errorf("failed to open a channel: %w", err)
	}

	return &RabbitMQConnection{
		Connection: conn,
		Channel:    ch,
	}, nil
}

// Close closes the RabbitMQ connection and channel
func (r *RabbitMQConnection) Close() {
	if r.Channel != nil {
		r.Channel.Close()
	}
	if r.Connection != nil {
		r.Connection.Close()
	}
}

// ConsumeQueue starts consuming messages from a queue
func (r *RabbitMQConnection) ConsumeQueue(queueName string) (<-chan amqp.Delivery, error) {
	// Declare the queue to ensure it exists
	_, err := r.Channel.QueueDeclare(
		queueName, // name
		true,      // durable
		false,     // delete when unused
		false,     // exclusive
		false,     // no-wait
		nil,       // arguments
	)
	if err != nil {
		return nil, fmt.Errorf("failed to declare queue %s: %w", queueName, err)
	}

	// Start consuming messages
	msgs, err := r.Channel.Consume(
		queueName, // queue
		"",        // consumer
		true,      // auto-ack
		false,     // exclusive
		false,     // no-local
		false,     // no-wait
		nil,       // args
	)
	if err != nil {
		return nil, fmt.Errorf("failed to register a consumer for queue %s: %w", queueName, err)
	}

	log.Printf("Started consuming from queue: %s", queueName)
	return msgs, nil
}