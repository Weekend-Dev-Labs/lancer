package services

import (
	"bytes"
	"encoding/json"
	"net/http"
	"time"
)

type WebhookEvent string

type Webhook struct {
	endpoint string
}

const (
	EventSessionCreate     = WebhookEvent("SESSION_CREATED")
	EventSessionCancelled  = WebhookEvent("SESSION_CANCELLED")
	EventSessionCompleteed = WebhookEvent("SESSION_COMPLETED")

	EventFileUpload = WebhookEvent("FILE_UPLOAD_")
)

func NewWebhookNotifier(endpoint string) *Webhook {
	return &Webhook{
		endpoint: endpoint,
	}
}

func (wh *Webhook) SendEvent(event WebhookEvent, payload interface{}) error {
	if wh.endpoint != "" {
		data := map[string]interface{}{
			"event": event,
			"data":  payload,
		}

		jsonData, err := json.Marshal(data)

		if err != nil {
			return err
		}

		req, err := http.NewRequest(http.MethodPost, wh.endpoint, bytes.NewBuffer(jsonData))

		if err != nil {
			return err
		}

		req.Header.Set("Content-Type", "application/json")

		client := http.Client{
			Timeout: 10 * time.Second,
		}

		_, err = client.Do(req)

		if err != nil {
			return err
		}

		// TODO : metrics for webhook success & failure rates

		return nil
	}

	return nil
}
