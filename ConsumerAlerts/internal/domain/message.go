package domain

import (
	"fmt"
)

// SensorMessage represents a message from a sensor
type SensorMessage struct {
	NumeroSerie        string      `json:"numeroSerie"`
	Sensor             string      `json:"sensor"`
	FechaActivacion    string      `json:"fecha_activacion"`
	FechaDesactivacion string      `json:"fecha_desactivacion"`
	Estado             interface{} `json:"estado"` // Using interface{} to handle both int and string
}

// GetEstadoAsString returns the estado field as a string regardless of its original type
func (s *SensorMessage) GetEstadoAsString() string {
	switch v := s.Estado.(type) {
	case int:
		return fmt.Sprintf("%d", v)
	case float64:
		return fmt.Sprintf("%.2f", v)
	case string:
		return v
	default:
		return fmt.Sprintf("%v", v)
	}
}

// IsDHT22 checks if the sensor is a DHT22
func (s *SensorMessage) IsDHT22() bool {
	return s.Sensor == "DHT_22"
}

// IsAlert checks if the message represents an alert condition
func (s *SensorMessage) IsAlert() bool {
	// For DHT22, we don't have a simple alert condition since it's temperature data
	if s.IsDHT22() {
		return false
	}
	
	// For other sensors, estado=1 means alert
	if estado, ok := s.Estado.(float64); ok {
		return estado == 1
	}
	
	return false
}