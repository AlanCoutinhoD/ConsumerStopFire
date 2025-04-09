package domain

import "fmt"

// ErrQueueOperation representa errores relacionados con operaciones de cola
type ErrQueueOperation struct {
    QueueName string
    Operation string
    Err       error
}

func (e *ErrQueueOperation) Error() string {
    return fmt.Sprintf("queue operation failed: %s on queue %s: %v", e.Operation, e.QueueName, e.Err)
}

// ErrMessageProcessing representa errores en el procesamiento de mensajes
type ErrMessageProcessing struct {
    MessageType string
    Err         error
}

func (e *ErrMessageProcessing) Error() string {
    return fmt.Sprintf("message processing error for %s: %v", e.MessageType, e.Err)
}

type ErrAPIRequest struct {
    StatusCode int
    URL       string
}

func (e *ErrAPIRequest) Error() string {
    return fmt.Sprintf("API request failed with status code %d for URL %s", e.StatusCode, e.URL)
}