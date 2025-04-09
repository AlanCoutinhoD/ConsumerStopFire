package domain

import amqp "github.com/rabbitmq/amqp091-go"

// MessageQueueConnection define la interfaz para conexiones de cola de mensajes
type MessageQueueConnection interface {
    ConsumeQueue(queueName string) (<-chan amqp.Delivery, error)
    Close()
}

// MessageConsumer define la interfaz para el consumidor de mensajes
type MessageConsumer interface {
    StartConsuming() error
}