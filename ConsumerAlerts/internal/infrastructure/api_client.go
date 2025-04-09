package infrastructure

import (
    "bytes"
    "encoding/json"
    "net/http"
    "time"
    "github.com/yourusername/ConsumerAlerts/internal/domain"
)

type APIClient struct {
    client  *http.Client
    baseURL string
}

func NewAPIClient(baseURL string) *APIClient {
    return &APIClient{
        client:  &http.Client{Timeout: 10 * time.Second},
        baseURL: baseURL,
    }
}

func (c *APIClient) SendAlert(msg domain.SensorMessage) error {
    payload, err := json.Marshal(msg)
    if err != nil {
        return err
    }

    resp, err := c.client.Post(c.baseURL, "application/json", bytes.NewBuffer(payload))
    if err != nil {
        return err
    }
    defer resp.Body.Close()

    if resp.StatusCode >= 400 {
        return &domain.ErrAPIRequest{
            StatusCode: resp.StatusCode,
            URL:       c.baseURL,
        }
    }

    return nil
}