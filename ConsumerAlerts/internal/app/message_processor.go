package app

import (
    "encoding/json"
    "github.com/yourusername/ConsumerAlerts/internal/domain"
)

type MessageProcessor struct {
    logger Logger
}

func NewMessageProcessor(logger Logger) *MessageProcessor {
    return &MessageProcessor{logger: logger}
}

func (p *MessageProcessor) ProcessMessage(rawMessage []byte) (*domain.SensorMessage, error) {
    var sensorMsg domain.SensorMessage
    if err := json.Unmarshal(rawMessage, &sensorMsg); err != nil {
        return nil, &domain.ErrMessageProcessing{
            MessageType: "sensor",
            Err:        err,
        }
    }
    return &sensorMsg, nil
}