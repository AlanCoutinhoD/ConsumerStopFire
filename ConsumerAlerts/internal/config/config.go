package config

import (
    "fmt"
    "os"
    "github.com/joho/godotenv"
)

type RabbitMQConfig struct {
    Host     string
    Port     string
    User     string
    Password string
    Queues   QueueConfig
}

type QueueConfig struct {
    KY026 string
    MQ2   string
    MQ135 string
    DHT22 string
}

type Config struct {
    RabbitMQ RabbitMQConfig
    APIURL   string
}

func LoadConfig(envPath string) (*Config, error) {
    if err := godotenv.Load(envPath); err != nil {
        return nil, fmt.Errorf("error loading .env file: %w", err)
    }

    return &Config{
        RabbitMQ: RabbitMQConfig{
            Host:     os.Getenv("RABBITMQ_HOST"),
            Port:     os.Getenv("RABBITMQ_PORT"),
            User:     os.Getenv("RABBITMQ_USER"),
            Password: os.Getenv("RABBITMQ_PASSWORD"),
            Queues: QueueConfig{
                KY026: os.Getenv("RABBITMQ_QUEUE_KY026"),
                MQ2:   os.Getenv("RABBITMQ_QUEUE_MQ2"),
                MQ135: os.Getenv("RABBITMQ_QUEUE_MQ135"),
                DHT22: os.Getenv("RABBITMQ_QUEUE_DHT22"),
            },
        },
        APIURL: "http://98.85.68.200:8080/api/alerts",
    }, nil
}